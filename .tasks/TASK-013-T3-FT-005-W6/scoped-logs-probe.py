import subprocess,json,re,pathlib
logs=subprocess.check_output(['journalctl','-u','somonwatch.service','--since','2026-10-07 16:33:09 UTC','--no-pager','-o','cat'],text=True)
messages=['Somon poll completed','monitoring paused; fresh baseline advanced','baseline created','Somon poll failed','keyword detail failed','Telegram delivery failed','detail fetch/parse failed; ad remains unseen','keyword Telegram delivery failed; will retry','new exact ad sent','new fallback ad sent']
result={'window_start':'2026-10-07T16:33:09Z','messages':{m:logs.count(m) for m in messages},'error_count':len(re.findall(r'level=ERROR|"level":"ERROR"',logs)),'poll_counts':[],'staged_binary_absent':not pathlib.Path('/opt/somonwatch/somonwatch.release-check').exists(),'doctor_scratch_absent':not pathlib.Path('/var/lib/somonwatch/release-doctor-8b48c4cef112').exists()}
for line in logs.splitlines():
 if 'Somon poll completed' in line:result['poll_counts'].append({k:int(v) for k,v in re.findall(r'\b(cards|ordinary|new_ids|details|sent)=(\d+)',line)})
print(json.dumps(result,indent=2))
