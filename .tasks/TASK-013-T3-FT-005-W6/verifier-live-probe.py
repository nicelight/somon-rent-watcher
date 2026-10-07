import subprocess, pathlib, sqlite3, json, hashlib, os, datetime, re

EXPECTED='8b48c4cef11237716e1dbc471cf363018e48c613'
BACKUP=pathlib.Path('/var/backups/somonwatch/keyword-release-20261007T163309Z')
REPO=pathlib.Path('/root/somon-rent-watcher')
DB=pathlib.Path('/var/lib/somonwatch/somonwatch.db')
BINARY=pathlib.Path('/opt/somonwatch/somonwatch')
def run(*args): return subprocess.check_output(args,text=True).strip()
def sha(p): return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
def digest(s): return hashlib.sha256(s.encode()).hexdigest()
def read_db(p):
    with sqlite3.connect('file:'+str(p)+'?mode=ro',uri=True) as c:
        c.execute('BEGIN')
        assert c.execute('pragma integrity_check').fetchone()[0]=='ok'
        tables={n:c.execute('select * from "'+n.replace('"','""')+'"').fetchall() for n, in c.execute("select name from sqlite_master where type='table' and name not like 'sqlite_%'")}
        return tables
unit=dict(line.split('=',1) for line in run('systemctl','show','somonwatch.service','-p','ActiveState','-p','SubState','-p','MainPID','-p','User','-p','NRestarts').splitlines())
assert unit['ActiveState']=='active' and unit['SubState']=='running' and unit['User']=='somonwatch' and unit['NRestarts']=='0'
assert run('pgrep','-cx','somonwatch')=='1'
pid=unit['MainPID']; proc=pathlib.Path('/proc')/pid
assert os.readlink(proc/'exe')==str(BINARY)
assert '/system.slice/somonwatch.service' in (proc/'cgroup').read_text()
fds=[]
for fd in (proc/'fd').iterdir():
    try: fds.append(os.readlink(fd))
    except FileNotFoundError: pass
assert str(DB) in fds
assert sha(proc/'exe')==sha(BINARY)==sha(REPO/'dist/somonwatch')
assert run('git','-C',str(REPO),'rev-parse','HEAD')==EXPECTED
assert not run('git','-C',str(REPO),'status','--porcelain')
assert run('git','-C',str(REPO),'rev-parse','refs/remotes/origin/main')==EXPECTED
assert not run('git','-C',str(REPO),'diff',EXPECTED,'--','cmd','internal','scripts','deploy','go.mod','go.sum','VERSION')
version=run(str(BINARY),'version'); assert EXPECTED[:12] in version
assert 'not found' not in run('ldd',str(BINARY))
assert sha('/etc/systemd/system/somonwatch.service')==sha(REPO/'deploy/somonwatch.service')==sha(BACKUP/'somonwatch.service')
assert BACKUP.stat().st_uid==0 and BACKUP.stat().st_mode & 0o077==0
for name in ('somonwatch.db','somonwatch.env','SHA256SUMS'):
    st=(BACKUP/name).stat(); assert st.st_uid==0 and st.st_mode & 0o077==0
assert sha('/etc/somonwatch/somonwatch.env')==sha(BACKUP/'somonwatch.env')
old=read_db(BACKUP/'somonwatch.db'); live=read_db(DB)
assert set(old)<=set(live)
preserved={}
for n,rows in old.items():
    if n!='state':
        assert set(rows)<=set(live[n]), 'old rows absent/changed: '+n
        preserved[n]={'backup_rows':len(rows),'live_rows':len(live[n]),'all_backup_rows_preserved':True}
assert old['settings']==live['settings']
before=dict(old['state']); after=dict(live['state']); assert set(before)<=set(after)
for k,v in before.items():
    if k=='telegram_offset': assert int(after[k])>=int(v)
    elif k not in ('last_successful_poll_at','previous_ordinary_ids'): assert v==after[k], 'state changed: '+k
assert 'search_monitors' in live and 'search_ad_state' in live
assert len(live['search_monitors'])==0 and len(live['search_ad_state'])==0
containers=[]
for cid in run('docker','ps','-aq').splitlines():
    o=json.loads(run('docker','inspect',cid))[0]
    identity=' '.join([o.get('Name',''),o['Config'].get('Image',''),str(o['Config'].get('Cmd')),str(o['Config'].get('Entrypoint'))]).lower()
    assert 'somonwatch' not in identity and 'somon-rent-watcher' not in identity
    containers.append([o['Id'],o['Image'],o['State']['Status'],o['State']['StartedAt'],o['RestartCount']])
services='\n'.join(sorted(line.split()[0] for line in run('systemctl','list-units','--type=service','--state=running','--no-legend','--plain').splitlines() if line.split()[0]!='somonwatch.service'))
host={'env_sha256':sha('/etc/somonwatch/somonwatch.env'),'containers_sha256':digest(json.dumps(sorted(containers))),'container_count':len(containers),'running_services_sha256':digest(services),'listeners_sha256':digest('\n'.join(sorted(run('ss','-H','-lntu').splitlines()))),'routes_sha256':digest(run('ip','route','show','table','all')),'firewall_sha256':digest(re.sub(r'\[\d+:\d+\]', '[COUNTERS]', '\n'.join(line for line in run('iptables-save').splitlines() if not line.startswith('#')))),'selinux':run('getenforce'),'failed_units_sha256':digest(run('systemctl','--failed','--no-legend','--plain'))}
logs=run('journalctl','-u','somonwatch.service','--since','2026-10-07 16:33:09 UTC','--no-pager','-o','cat')
errors=len(re.findall(r'level=ERROR|"level":"ERROR"',logs))
assert errors==0
assert not pathlib.Path('/opt/somonwatch/somonwatch.release-check').exists()
assert not pathlib.Path('/var/lib/somonwatch/release-doctor-'+EXPECTED[:12]).exists()
print(json.dumps({'observed_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'expected_commit':EXPECTED,'checkout_commit':EXPECTED,'version':version,'unit':unit,'process_installed_build_sha256':sha(BINARY),'process_unit_db_bound':True,'sqlite_integrity':'ok','backup_root_only':True,'old_tables':preserved,'prior_state_keys_preserved':True,'stable_state_values_preserved':True,'telegram_offset_monotonic':True,'settings_exactly_preserved':True,'db_identity':{'inode':DB.stat().st_ino,'device':DB.stat().st_dev},'search_monitors':len(live['search_monitors']),'search_ad_state':len(live['search_ad_state']),'unrelated':host,'scoped_log_errors':errors,'new_baseline_log_count':logs.count('baseline created'),'stage_and_scratch_absent':True},indent=2))
