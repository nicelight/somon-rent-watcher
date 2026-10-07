#!/usr/bin/env python3
"""Disposable snapshot verification; no repository source/runtime writes."""
import datetime
import hashlib
import json
import pathlib
import shutil
import subprocess
import tempfile

root=pathlib.Path(__file__).resolve().parents[2]
artifact=pathlib.Path(__file__).resolve().parent
directories=['cmd','internal','scripts','testdata']
files=['go.mod','VERSION']
source_files=sorted([p for d in directories for p in (root/d).rglob('*') if p.is_file()]+[root/f for f in files])
hashes={str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in source_files}
image=subprocess.check_output(['docker','image','inspect','somon-price-hotfix-builder:latest','--format','{{.Id}}'],text=True).strip()
state={'captured_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'image':image,'files':hashes}
(artifact/'verifier-source-state.json').write_text(json.dumps(state,indent=2)+'\n')
results=[]
with tempfile.TemporaryDirectory(prefix='disposable-verifier-',dir=artifact) as temporary:
    copied=pathlib.Path(temporary)
    for d in directories: shutil.copytree(root/d,copied/d)
    for f in files: shutil.copy2(root/f,copied/f)
    def run(name, command):
        args=['docker','run','--rm','--network','none','-v',f'{copied}:/src','-w','/src',image]+command
        result=subprocess.run(args,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
        (artifact/f'verifier-{name}.log').write_text(result.stdout)
        results.append({'name':name,'command':args,'exit_code':result.returncode,'completed_at':datetime.datetime.now(datetime.timezone.utc).isoformat(),'evidence':f'verifier-{name}.log'})
        print(name,'exit_code=',result.returncode,flush=True)
        print(result.stdout[-1200:],flush=True)
        return result.returncode
    run('focused-gate',['sh','-c','CGO_ENABLED=1 go test -count=1 ./internal/somon ./internal/app ./internal/store'])
    run('native-gate',['bash','scripts/build.sh'])
    shutil.copy2(artifact/'verifier_history_probe_test.go',copied/'internal/store/verifier_history_probe_test.go')
    shutil.copy2(artifact/'verifier_polling_history_probe_test.go',copied/'internal/app/verifier_polling_history_probe_test.go')
    run('outcome',['sh','-c','gofmt -w internal/store/verifier_history_probe_test.go internal/app/verifier_polling_history_probe_test.go && CGO_ENABLED=1 go test -count=1 -v -run "^TestVerifierCurrentFeed" ./internal/store ./internal/app'])
state['source_unchanged_after']=all(p.exists() and hashlib.sha256(p.read_bytes()).hexdigest()==hashes[str(p.relative_to(root))] for p in source_files)
(artifact/'verifier-source-state.json').write_text(json.dumps(state,indent=2)+'\n')
(artifact/'verifier-commands-results.json').write_text(json.dumps(results,indent=2)+'\n')
print('source_unchanged_after=',state['source_unchanged_after'],flush=True)
raise SystemExit(0 if all(r['exit_code']==0 for r in results) and state['source_unchanged_after'] else 1)
