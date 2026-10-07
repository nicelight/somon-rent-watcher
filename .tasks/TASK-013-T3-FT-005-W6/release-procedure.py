import os,sys,pathlib,subprocess,sqlite3,json,hashlib,shutil,datetime,time
EXPECTED=sys.argv[1]
BASE='93cbe9ebcfe6421040461f4ee4e049430bf5028d'
REPO=pathlib.Path('/root/somon-rent-watcher')
BINARY=pathlib.Path('/opt/somonwatch/somonwatch')
STAGE=pathlib.Path('/opt/somonwatch/somonwatch.release-check')
ENV=pathlib.Path('/etc/somonwatch/somonwatch.env')
DB=pathlib.Path('/var/lib/somonwatch/somonwatch.db')
UNIT=pathlib.Path('/etc/systemd/system/somonwatch.service')
def out(s):print(s,flush=True)
def run(args,**kw):return subprocess.check_output(args,text=True,**kw).strip()
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
def state():
 with sqlite3.connect('file:'+str(DB)+'?mode=ro',uri=True) as c:
  assert c.execute('pragma integrity_check').fetchone()[0]=='ok'
  return {'settings':c.execute('select json from settings where id=1').fetchone()[0], 'seen':c.execute('select ad_id,first_seen_at from seen_ads order by ad_id').fetchall(),'state':dict(c.execute('select key,value from state').fetchall()), 'tables':{n:c.execute('select * from "'+n+'"').fetchall() for (n,) in c.execute("select name from sqlite_master where type='table' and name not like 'sqlite_%'").fetchall()}, 'inode':DB.stat().st_ino,'device':DB.stat().st_dev}
def stable_equal(a,b):
 assert a['settings']==b['settings'],'settings changed'
 assert a['inode']==b['inode'] and a['device']==b['device'],'database replaced'
 assert set(a['seen']).issubset(set(b['seen'])),'seen rows lost or changed'
 assert set(a['state']).issubset(b['state']),'state keys lost'
 for k,v in a['state'].items():
  if k=='telegram_offset':assert int(b['state'][k])>=int(v),'Telegram offset reset'
  elif k not in ('last_successful_poll_at','previous_ordinary_ids'):assert b['state'][k]==v,'persisted state changed: '+k
 for n,rows in a['tables'].items():
  if n!='state':assert set(rows).issubset(set(b['tables'].get(n,[]))),'table rows lost: '+n
def health():
 v=dict(line.split('=',1) for line in run(['systemctl','show','somonwatch.service','-p','ActiveState','-p','SubState','-p','User','-p','NRestarts','-p','MainPID']).splitlines())
 assert v['ActiveState']=='active' and v['SubState']=='running' and v['User']=='somonwatch',v
 assert run(['pgrep','-cx','somonwatch'])=='1','ambiguous runtime'
 pid=v['MainPID'];proc=pathlib.Path('/proc')/pid
 assert os.readlink(proc/'exe')==str(BINARY),'unexpected running executable'
 assert sha(proc/'exe')==sha(BINARY),'running executable differs from installed file'
 assert '/system.slice/somonwatch.service' in (proc/'cgroup').read_text(),'process is outside watcher unit'
 fds=[]
 for fd in (proc/'fd').iterdir():
  try:fds.append(os.readlink(fd))
  except FileNotFoundError:pass
 assert str(DB) in fds,'expected SQLite file is not open by watcher'
 for cid in run(['docker','ps','-q']).splitlines():
  o=json.loads(run(['docker','inspect',cid]))[0]
  identity=' '.join([o.get('Name',''),o['Config'].get('Image',''),str(o['Config'].get('Cmd')),str(o['Config'].get('Entrypoint'))]).lower()
  assert 'somonwatch' not in identity and 'somon-rent-watcher' not in identity,'competing watcher container'
 return v
def wait_health():
 deadline=time.monotonic()+15
 while True:
  try:return health()
  except (AssertionError,FileNotFoundError,ProcessLookupError,subprocess.CalledProcessError):
   if time.monotonic()>=deadline:raise
   time.sleep(0.25)
assert len(EXPECTED)==40 and all(c in '0123456789abcdef' for c in EXPECTED)
health();assert not run(['git','-C',str(REPO),'status','--porcelain'])
current_head=run(['git','-C',str(REPO),'rev-parse','HEAD'])
assert current_head in (BASE,EXPECTED),'production baseline changed'
installed_version=run([str(BINARY),'version']);out('AC010 installed baseline: '+installed_version+'; expected='+EXPECTED)
assert BASE[:12] in installed_version or EXPECTED[:12] in installed_version,'installed version differs from accepted baseline/hotfix'
assert sha(UNIT)==sha(REPO/'deploy/somonwatch.service'),'installed unit differs from accepted source'
env_hash=sha(ENV);initial=state();old_hash=sha(BINARY)
out('Preflight PASS: healthy single watcher, process/binary/database bound, clean baseline, matching unit, SQLite integrity OK')
if EXPECTED[:12] in run([str(BINARY),'version']) and current_head==EXPECTED:
 out('Hotfix is already installed and active; no install replay needed');sys.exit(0)
subprocess.run(['git','-C',str(REPO),'fetch','origin','main'],check=True)
assert run(['git','-C',str(REPO),'rev-parse','FETCH_HEAD'])==EXPECTED
subprocess.run(['git','-C',str(REPO),'merge','--ff-only',EXPECTED],check=True)
assert not run(['git','-C',str(REPO),'status','--porcelain'])
out('Source synchronized to '+EXPECTED)
build_env=os.environ.copy();build_env.update({'COMMIT':EXPECTED[:12],'GOMAXPROCS':'2','GOFLAGS':'-p=2'})
subprocess.run(['nice','-n','10','./scripts/build.sh'],cwd=REPO,env=build_env,check=True)
assert run(['git','-C',str(REPO),'rev-parse','HEAD'])==EXPECTED
assert not run(['git','-C',str(REPO),'status','--porcelain'])
BUILT=REPO/'dist/somonwatch';new_hash=sha(BUILT)
assert EXPECTED[:12] in run([str(BUILT),'version'])
assert 'not found' not in run(['ldd',str(BUILT)])
assert not STAGE.exists(),'unexpected staging file already exists'
def doctor_check(binary,clone):
 scratch=pathlib.Path('/var/lib/somonwatch/release-doctor-'+EXPECTED[:12])
 assert not scratch.exists(),'unexpected scratch directory'
 try:
  scratch.mkdir(mode=0o700)
  clone_db=scratch/'doctor.db'
  if clone:
   with sqlite3.connect('file:'+str(DB)+'?mode=ro',uri=True) as src,sqlite3.connect(str(clone_db)) as dst:src.backup(dst)
  import pwd
  owner=pwd.getpwnam('somonwatch')
  for f in [scratch]+list(scratch.iterdir()):os.chown(f,owner.pw_uid,owner.pw_gid)
  subprocess.run(['restorecon','-RF',str(scratch)],check=True)
  command='set -a; . /etc/somonwatch/somonwatch.env; '
  if clone:command+='export DB_PATH='+str(clone_db)+'; '
  command+='export DEBUG_DIR='+str(scratch)+'; set +a; exec runuser -u somonwatch --preserve-environment -- '+str(binary)+' doctor'
  doctor=subprocess.run(['bash','-c',command],capture_output=True,text=True)
  for line in doctor.stdout.splitlines():
   if line.startswith('[OK] Telegram bot:'):out('[OK] Telegram identity checked (redacted)')
   elif line.startswith('[OK] Target chat:'):out('[OK] Target chat checked (redacted)')
   elif line.startswith('[OK]') or line.startswith('[WARN]'):out(line)
  assert doctor.returncode==0,'doctor failed, exit='+str(doctor.returncode)+' (sensitive output withheld)'
  out('Doctor PASS; clone='+str(clone))
 finally:
  if scratch.exists():shutil.rmtree(scratch)
try:
 shutil.copy2(BUILT,STAGE);STAGE.chmod(0o755)
 subprocess.run(['restorecon',str(STAGE)],check=True)
 doctor_check(STAGE,True)
finally:
 STAGE.unlink(missing_ok=True)
assert sha(ENV)==env_hash;stable_equal(initial,state());health()
assert sha(BINARY)==old_hash
BACKUP=pathlib.Path('/var/backups/somonwatch/keyword-release-'+datetime.datetime.now(datetime.timezone.utc).strftime('%Y%m%dT%H%M%SZ'))
subprocess.run(['./scripts/backup-installed.sh',str(BACKUP)],cwd=REPO,check=True)
assert sha(BACKUP/'somonwatch')==old_hash
out('Backup verified; stopping only somonwatch.service')
stopped=False
try:
 stopped=True
 subprocess.run(['systemctl','stop','somonwatch.service'],check=True)
 before=state()
 subprocess.run(['./scripts/install-almalinux.sh'],cwd=REPO,check=True)
 assert sha(BINARY)==new_hash and sha(ENV)==env_hash
 after_install=state();stable_equal(before,after_install);assert before['tables']==after_install['tables'],'existing tables changed during stopped install'
 out('Stopped-state comparison PASS: same DB inode, exact settings/seen rows and environment preserved')
 subprocess.run(['systemctl','start','somonwatch.service'],check=True)
 current=wait_health();assert current['NRestarts']=='0',current
 doctor_check(BINARY,False)
 after=state();stable_equal(before,after);assert sha(ENV)==env_hash
 assert 'search_monitors' in after['tables'] and 'search_ad_state' in after['tables'],'missing additive tables'
 assert not after['tables']['search_monitors'],'unexpected auto-created monitor'
 out(json.dumps({'installation':'PASS; final host comparison pending','commit':EXPECTED,'binary_sha256':new_hash,'backup':str(BACKUP),'seen_before':len(before['seen']),'seen_after':len(after['seen']),'settings_preserved':True,'environment_preserved':True,'database_identity_preserved':True,'prior_state_keys_preserved':True,'telegram_offset_not_reset':True,'old_tables_preserved':True,'keyword_monitors':len(after['tables']['search_monitors']),'additive_tables_present':True,'doctor_scratch_cleaned':True,'unit':current}))
except BaseException:
 if stopped:
  out('Release failed after stop; restoring only the watcher executable and service')
  subprocess.run(['systemctl','stop','somonwatch.service'],check=True)
  rollback=BINARY.with_name('somonwatch.rollback-tmp');shutil.copy2(BACKUP/'somonwatch',rollback);os.replace(rollback,BINARY)
  subprocess.run(['systemctl','start','somonwatch.service'],check=True)
  out('Rollback health: '+json.dumps(wait_health()))
 raise
