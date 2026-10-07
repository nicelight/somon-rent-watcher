//go:build cgo
package app
import (
 "context"
 "sync"
 "testing"
 "time"
 "github.com/nicelight/somon-rent-watcher/internal/model"
)
func TestVerifierStaleEvaluationOutcome(t *testing.T){
 for _,kind:=range []string{"delete","edit-reject","edit-send"}{t.Run(kind,func(t *testing.T){
 h:=newKeywordPollingHarness(t,3);max:=200;s:=h.search("verifier",&max);entered,release:=make(chan struct{}),make(chan struct{});var once,open sync.Once
 h.mu.Lock();h.detailHook=func(){once.Do(func(){close(entered);<-release})};if kind=="edit-reject"{h.detailPrice=270};h.mu.Unlock()
 ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second);defer cancel();done:=make(chan error,1);go func(){done<-h.a.pollOnce(ctx)}();t.Cleanup(func(){open.Do(func(){close(release)})})
 select{case <-entered:case <-ctx.Done():open.Do(func(){close(release)});<-done;t.Fatal("detail boundary not reached")}
 switch kind {case "delete":if ok,e:=h.a.DeleteKeywordSearch(s.ID);e!=nil||!ok{open.Do(func(){close(release)});<-done;t.Fatal(ok,e)}
 default:current,ok,e:=h.a.KeywordSearch(s.ID);if e!=nil||!ok{t.Fatal(ok,e)};limit:=300;if kind=="edit-send"{limit=100};current.PriceMax=&limit;updated,ok,e:=h.a.UpdateKeywordSearch(current);if e!=nil||!ok||updated.Revision<=s.Revision{t.Fatal("revision edit",ok,e)};t.Logf("old revision=%d new=%d",s.Revision,updated.Revision)}
 open.Do(func(){close(release)});if e:=<-done;e!=nil{t.Fatal(e)};h.mu.Lock();sends:=len(h.captions);h.detailHook=nil;if kind=="edit-send"{h.price=80};h.mu.Unlock()
 if sends!=0||len(h.history(s))!=0{t.Fatal("stale evaluation sent/committed rejection",sends,h.history(s))}
 if kind=="delete"{for _,flag:=range []bool{false,true}{if e:=h.db.RecordKeywordAdState(s.ID,1001,s.Revision,flag);e!=nil{t.Fatal(e)}};h.reopen();h.poll();if _,found,e:=h.a.KeywordSearch(s.ID);e!=nil||found||len(h.history(s))!=0{t.Fatal("deleted ID/history returned")};h.mu.Lock();sends=len(h.captions);h.mu.Unlock();if sends!=0{t.Fatal("future deleted poll sent")}
 }else{h.poll();hist:=h.history(s);h.mu.Lock();sends=len(h.captions);h.mu.Unlock();if sends!=1||!hist[1001].Delivered||hist[1001].EvaluatedRevision<=s.Revision{t.Fatal("new revision not reevaluated",sends,hist)}}
 t.Logf("AC002 %s detail-start→mutation-complete→detail-release→0 stale sends/history; subsequent current poll checked",kind)
 })}
}
func TestVerifierStartedRequestDeletionOutcome(t *testing.T){
 h:=newKeywordPollingHarness(t,3);s:=h.search("already-started",nil);entered,release:=make(chan struct{}),make(chan struct{});var once,open sync.Once
 h.mu.Lock();h.sendHook=func(){once.Do(func(){close(entered);<-release})};h.mu.Unlock();ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second);defer cancel();done:=make(chan error,1);go func(){done<-h.a.pollOnce(ctx)}();t.Cleanup(func(){open.Do(func(){close(release)})})
 select{case <-entered:case <-ctx.Done():open.Do(func(){close(release)});<-done;t.Fatal("HTTP request not received")}
 type result struct{found bool;err error};deleted:=make(chan result,1);invoked:=make(chan struct{});go func(){close(invoked);ok,e:=h.a.DeleteKeywordSearch(s.ID);deleted<-result{ok,e}}();<-invoked
 open.Do(func(){close(release)});if e:=<-done;e!=nil{t.Fatal(e)};var res result;select{case res=<-deleted:case <-ctx.Done():t.Fatal("delete did not finish")};if !res.found||res.err!=nil{t.Fatal(res)}
 h.mu.Lock();sends:=len(h.captions);h.sendHook=nil;h.mu.Unlock();if sends!=1||len(h.history(s))!=0{t.Fatal("started request/deletion outcome",sends,h.history(s))}
 if sent,e:=h.a.deliverKeywordAd(ctx,s,model.Ad{Card:model.Card{ID:1002}});e!=nil||sent{t.Fatal("captured stale candidate started send",sent,e)}
 for _,flag:=range []bool{true,false}{if e:=h.db.RecordKeywordAdState(s.ID,1002,s.Revision,flag);e!=nil{t.Fatal(e)}};h.reopen();h.poll();h.mu.Lock();sends=len(h.captions);h.mu.Unlock();if _,found,e:=h.a.KeywordSearch(s.ID);e!=nil||found||sends!=1||len(h.history(s))!=0{t.Fatal("post-delete resurrection/send",found,sends,e)}
 t.Log("AC002 Telegram HTTP received→concurrent delete invoked→success allowed→delete completed→reopen absent; exactly1 total started request, no later captured-candidate/poll send or conditional history resurrection")
}
