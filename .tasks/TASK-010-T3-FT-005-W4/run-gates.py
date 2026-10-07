from pathlib import Path
import datetime,hashlib,json,subprocess
root=Path.cwd();e=root/'.tasks/TASK-010-T3-FT-005-W4';results=[]
for name,body in [('focused-package-gate','CGO_ENABLED=1 go test -count=1 ./internal/app ./internal/store ./internal/telegram'),('native-build-gate','./scripts/build.sh')]:
 files=sorted(p for prefix in ['internal','cmd','scripts','testdata'] for p in (root/prefix).rglob('*') if p.is_file())+[root/'go.mod',root/'VERSION']
 basis={'head':subprocess.check_output(['git','rev-parse','HEAD'],text=True).strip(),'status':subprocess.check_output(['git','status','--short'],text=True),'hashes':{str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in files},'at':datetime.datetime.now().astimezone().isoformat(),'toolchain':'existing somon-price-hotfix-builder:latest Go1.21 CGO/gcc/sqlite header; network none; only temporary test DB/HTTP'}
 (e/(name+'-input-state.json')).write_text(json.dumps(basis,indent=2)+'\n')
 command=['docker','run','--rm','--network','none','-v',str(root)+':/src','-w','/src','somon-price-hotfix-builder:latest','sh','-c',body]
 with (e/(name+'.log')).open('w') as out:r=subprocess.run(command,stdout=out,stderr=subprocess.STDOUT)
 results.append({'name':name,'command':command,'cwd':str(root),'exit_code':r.returncode,'completed_at':datetime.datetime.now().astimezone().isoformat(),'input_snapshot':str((e/(name+'-input-state.json')).relative_to(root)),'evidence':str((e/(name+'.log')).relative_to(root)),'log_sha256':hashlib.sha256((e/(name+'.log')).read_bytes()).hexdigest()})
 (e/'commands-results.json').write_text(json.dumps(results,indent=2)+'\n')
 print(name+' exit '+str(r.returncode),flush=True)
 if r.returncode:break
