import pathlib,sqlite3,hashlib,subprocess,json,os
expected='d5fdc120400acba96d088c465aa71db9ab1402ab'
backup=pathlib.Path('/var/backups/somonwatch/price-hotfix-20260909T085033Z')
def run(*a):return subprocess.check_output(a,text=True).strip()
def sha(p):return hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
pid=run('systemctl','show','somonwatch.service','-p','MainPID','--value')
assert sha('/proc/'+pid+'/exe')==sha('/opt/somonwatch/somonwatch')==sha('/root/somon-rent-watcher/dist/somonwatch')
assert run('git','-C','/root/somon-rent-watcher','rev-parse','HEAD')==expected
assert not run('git','-C','/root/somon-rent-watcher','status','--porcelain')
assert expected[:12] in run('/opt/somonwatch/somonwatch','version')
old=sqlite3.connect('file:'+str(backup/'somonwatch.db')+'?mode=ro',uri=True);now=sqlite3.connect('file:/var/lib/somonwatch/somonwatch.db?mode=ro',uri=True)
assert now.execute('pragma integrity_check').fetchone()[0]=='ok'
a=set(old.execute('select ad_id,first_seen_at from seen_ads'));b=set(now.execute('select ad_id,first_seen_at from seen_ads'));assert a<=b
assert old.execute('select json from settings where id=1').fetchone()==now.execute('select json from settings where id=1').fetchone()
assert sha(backup/'somonwatch.env')==sha('/etc/somonwatch/somonwatch.env')
assert '7f9c5f50d659' in run(str(backup/'somonwatch'),'version')
assert run('systemctl','show','somonwatch.service','-p','NRestarts','--value')=='0'
print(json.dumps({'functional_postflight':'PASS','running_installed_built_sha256':sha('/opt/somonwatch/somonwatch'),'exact_source_commit':expected,'backup_rows_preserved':len(a),'current_rows':len(b),'settings_and_environment_preserved':True,'rollback_binary_verified':True,'sqlite_integrity':'ok','restarts':0},indent=2))
