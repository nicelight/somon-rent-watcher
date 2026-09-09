import subprocess,json,hashlib,sqlite3,pathlib,re

def run(*args):return subprocess.check_output(args,text=True).strip()
def digest(s):return hashlib.sha256(s.encode()).hexdigest()
unit=dict(line.split('=',1) for line in run('systemctl','show','somonwatch.service','-p','ActiveState','-p','SubState','-p','NRestarts','-p','MainPID','-p','User').splitlines())
assert unit['ActiveState']=='active' and unit['SubState']=='running' and unit['User']=='somonwatch',unit
assert run('pgrep','-cx','somonwatch')=='1','ambiguous watcher process count'
repo='/root/somon-rent-watcher';assert not run('git','-C',repo,'status','--porcelain'),'dirty production checkout'
c=sqlite3.connect('file:/var/lib/somonwatch/somonwatch.db?mode=ro',uri=True)
integrity=c.execute('pragma integrity_check').fetchone()[0];assert integrity=='ok',integrity
settings=c.execute('select json from settings where id=1').fetchone()[0]
seen=c.execute('select ad_id,first_seen_at from seen_ads order by ad_id').fetchall()
ids=run('docker','ps','-aq').splitlines()
containers=[]
for cid in ids:
 o=json.loads(run('docker','inspect',cid))[0];containers.append([o['Id'],o['Image'],o['State']['Status'],o['State']['StartedAt'],o['RestartCount']])
services='\n'.join(sorted(line.split()[0] for line in run('systemctl','list-units','--type=service','--state=running','--no-legend','--plain').splitlines() if line.split()[0]!='somonwatch.service'))
result={'unit':unit,'version':run('/opt/somonwatch/somonwatch','version'),'checkout_commit':run('git','-C',repo,'rev-parse','HEAD'),'sqlite_integrity':integrity,'settings_sha256':digest(settings),'seen_count':len(seen),'seen_sha256':digest(json.dumps(seen)),'env_sha256':hashlib.sha256(pathlib.Path('/etc/somonwatch/somonwatch.env').read_bytes()).hexdigest(),'containers_sha256':digest(json.dumps(sorted(containers))),'container_count':len(containers),'running_services_sha256':digest(services),'listeners_sha256':digest('\n'.join(sorted(run('ss','-H','-lntu').splitlines()))),'routes_sha256':digest(run('ip','route','show','table','all')),'firewall_sha256':digest(re.sub(r'\[\d+:\d+\]', '[COUNTERS]', '\n'.join(line for line in run('iptables-save').splitlines() if not line.startswith('#')))),'selinux':run('getenforce'),'failed_units_sha256':digest(run('systemctl','--failed','--no-legend','--plain'))}
print(json.dumps(result,indent=2))
