//go:build cgo

package app

import (
 "context"
 "fmt"
 "strings"
 "sync"
 "testing"
 "time"
 "github.com/nicelight/somon-rent-watcher/internal/model"
)

// Verifier-owned assertions use the existing local HTTP/DB fixture wiring only.
// Production behavior and executor TestKeywordPolling assertions are unmodified.
func verifierEdit(t *testing.T, h *keywordPollingHarness, s model.KeywordSearch, max int) model.KeywordSearch {
 t.Helper()
 current, found, err := h.a.KeywordSearch(s.ID)
 if err != nil || !found { t.Fatal(found,err) }
 current.PriceMax=&max
 updated, found, err := h.a.UpdateKeywordSearch(current)
 if err != nil || !found { t.Fatal(found,err) }
 return updated
}
func verifierWait(t *testing.T, condition func() bool) {
 t.Helper(); deadline:=time.Now().Add(2*time.Second)
 for !condition() { if time.Now().After(deadline) { t.Fatal("local observation timeout") }; time.Sleep(time.Millisecond) }
}
func TestVerifierSearchDurableOutcome(t *testing.T) {
 h:=newKeywordPollingHarness(t,8)
 h.mu.Lock();h.ids=[]int64{7311};h.price=190;h.mu.Unlock()
 budget:=200
 a:=h.search("рабочий стол",&budget)
 if err:=h.db.MarkSeen([]int64{7311},time.Now());err!=nil {t.Fatal(err)}
 h.poll()
 first:=h.history(a)[7311]
 if !first.Delivered || first.EvaluatedRevision!=a.Revision || len(h.captions)!=1 {t.Fatal(first,h.captions)}
 for _,value:=range []string{"рабочий стол","Стол деревянный","190 c.","Душанбе"} {if !strings.Contains(h.captions[0],value){t.Fatal("missing payload",value,h.captions)}}
 if !strings.Contains(h.photos[0],"7311.jpg") || !strings.Contains(h.buttons[0],"7311_table") {t.Fatal(h.photos,h.buttons)}
 // Enable another search after the first ID is already delivered in A.
 b:=h.search("стол для дома",&budget);h.poll()
 if !h.history(b)[7311].Delivered || len(h.captions)!=2 {t.Fatal("cross-search exclusion",h.history(b))}
 h.mu.Lock();h.ids=[]int64{7311,7312};h.price=220;h.mu.Unlock();h.poll()
 for _,s:=range []model.KeywordSearch{a,b} {v:=h.history(s);if len(v)!=2 || v[7312].Delivered || v[7312].EvaluatedRevision!=s.Revision {t.Fatal(v)}}
 // A raised budget reevaluates only A; B retains its revision rejection.
 a=verifierEdit(t,h,a,240);h.reopen();h.poll()
 if !h.history(a)[7312].Delivered || h.history(b)[7312].Delivered || len(h.captions)!=3 {t.Fatal(h.history(a),h.history(b),h.captions)}
 if h.history(a)[7311]!=first {t.Fatal("delivered evidence rewritten",h.history(a))}
 h.mu.Lock();h.price=100;h.ids=[]int64{7311,7312,7313};h.mu.Unlock();h.poll()
 if len(h.captions)!=5 || !h.history(a)[7313].Delivered || !h.history(b)[7313].Delivered || h.history(b)[7312].Delivered {t.Fatal("new ID/revision lifecycle",h.history(a),h.history(b),h.captions)}
 b=verifierEdit(t,h,b,230);h.poll();h.reopen();h.poll()
 if len(h.captions)!=6 || len(h.history(a))!=3 || len(h.history(b))!=3 || !h.history(b)[7312].Delivered {t.Fatal("restart/edit/drop repeat",h.captions,h.history(a),h.history(b))}
 t.Logf("AC005 six confirmed sends; independent histories A=%+v B=%+v; original delivery row preserved=%+v",h.history(a),h.history(b),first)
}

func TestVerifierSearchFailuresPreserveExistingHistory(t *testing.T) {
 for _,failure:=range []string{"source-http","source-parse","detail-http","telegram-error","telegram-ambiguity"} {
 t.Run(failure,func(t *testing.T){
 h:=newKeywordPollingHarness(t,3);s:=h.search("письменный стол",nil);h.poll();before:=h.history(s)[1001]
 h.mu.Lock();h.ids=[]int64{1001,7422}
 switch failure {case "source-http":h.sourceStatus=503;case "source-parse":h.sourceBroken=true;case "detail-http":h.detailStatus=502;case "telegram-error":h.telegramStatus=500;case "telegram-ambiguity":h.telegramAmbiguous=true};h.mu.Unlock()
 err:=h.a.pollOnce(context.Background());if strings.HasPrefix(failure,"source")&&err==nil {t.Fatal("source failure not surfaced")}
 if v:=h.history(s);len(v)!=1||v[1001]!=before {t.Fatal("failed candidate advanced history",v)}
 h.mu.Lock();failedSends:=len(h.captions);trace:=append([]keywordTrace(nil),h.trace...);h.sourceStatus=0;h.sourceBroken=false;h.detailStatus=0;h.telegramStatus=0;h.telegramAmbiguous=false;h.mu.Unlock()
 h.reopen();h.poll()
 if v:=h.history(s);len(v)!=2||!v[7422].Delivered||v[1001]!=before {t.Fatal("retry lost",v)}
 h.mu.Lock();sends:=len(h.captions);h.mu.Unlock()
 if sends!=failedSends+1 {t.Fatal("wrong retry count",failedSends,sends)}
 t.Logf("AC006 %s existing row preserved=%+v; absent7422 -> delivered after reopen; trace=%v",failure,before,trace)
 })}
}

func TestVerifierSearchCapFairness(t *testing.T) {
 h:=newKeywordPollingHarness(t,1);a:=h.search("стол один",nil);b:=h.search("стол два",nil)
 h.poll();settings,err:=h.a.LoadSettings();if err!=nil {t.Fatal(err)};settings.Enabled=true;if err=h.a.SaveSettings(settings);err!=nil {t.Fatal(err)}
 h.mu.Lock();h.ids=[]int64{1001,7555};h.rentalNew=true;h.mu.Unlock()
 starts:=[]string{a.Phrase,b.Phrase,""};totalDetails:=0
 for cycle:=0;cycle<3;cycle++ {
 h.mu.Lock();n:=len(h.trace);h.mu.Unlock();h.poll();h.mu.Lock();trace:=append([]keywordTrace(nil),h.trace[n:]...);h.mu.Unlock()
 if len(trace)==0||trace[0].phrase!=starts[cycle] {t.Fatal("rotation",cycle,trace)}
 details:=0;for _,r:=range trace {if strings.HasPrefix(r.path,"/adv/"){details++}}
 if details!=1 {t.Fatal("shared rental + keyword cap",cycle,details,trace)};totalDetails+=details
 t.Logf("AC006 cycle%d one shared detail, rotating start=%q trace=%v",cycle,starts[cycle],trace)
 }
 if !h.history(a)[7555].Delivered||!h.history(b)[1001].Delivered||h.history(b)[7555].Delivered {t.Fatal("cap-deferred state",h.history(a),h.history(b))}
 seen,err:=h.db.SeenIDs([]int64{9002});if err!=nil||!seen[9002] {t.Fatal("rental budget starvation",seen,err)}
 h.poll();if !h.history(b)[7555].Delivered {t.Fatal("deferred new ID not retried",h.history(b))}
 h.mu.Lock();trace:=append([]keywordTrace(nil),h.trace...);max:=h.maxInFlight;h.mu.Unlock()
 if max!=1||totalDetails!=3 {t.Fatal(max,totalDetails)}
 spacing:=time.Hour;for i:=1;i<len(trace);i++ {d:=trace[i].at.Sub(trace[i-1].at);if d<spacing {spacing=d};if d<4*time.Millisecond {t.Fatal("shared configured5ms delay absent",d)}}
 t.Logf("AC006 sequential HTTP max=%d; configured5ms minimum server spacing=%s (1ms arrival tolerance)",max,spacing)
}

func TestVerifierSearchRevisionBarrier(t *testing.T) {
 for _,staleReject:=range []bool{true,false} {t.Run(fmt.Sprint(staleReject),func(t *testing.T){
 h:=newKeywordPollingHarness(t,2);budget:=200;s:=h.search("стол",&budget)
 entered,release:=make(chan struct{}),make(chan struct{});var once sync.Once;var releaseOnce sync.Once
 h.mu.Lock();h.detailHook=func(){once.Do(func(){close(entered);<-release})};if staleReject {h.detailPrice=250};h.mu.Unlock()
 done:=make(chan error,1);go func(){done<-h.a.pollOnce(context.Background())}()
 defer releaseOnce.Do(func(){close(release)})
 select {case <-entered:case <-time.After(time.Second):t.Fatal("detail barrier not entered")}
 newBudget:=100;if staleReject {newBudget=300};updated:=verifierEdit(t,h,s,newBudget)
 releaseOnce.Do(func(){close(release)});if err:=<-done;err!=nil {t.Fatal(err)}
 h.mu.Lock();sends:=len(h.captions);h.detailHook=nil;if !staleReject {h.price=90};h.mu.Unlock()
 if len(h.history(s))!=0||sends!=0 {t.Fatal("stale result survived revision",h.history(s),sends)}
 h.poll();if v:=h.history(s)[1001];!v.Delivered||v.EvaluatedRevision!=updated.Revision {t.Fatal("current revision not evaluated",v)}
 t.Logf("AC006 stale rejection=%v history/sends empty before current revision%d succeeds",staleReject,updated.Revision)
 })}
}

func TestVerifierSearchBackoffAndRunningManual(t *testing.T) {
 for _,detailBlock:=range []bool{false,true} {for _,code:=range []int{403,429} {t.Run(fmt.Sprintf("detail%v-%d",detailBlock,code),func(t *testing.T){
 h:=newKeywordPollingHarness(t,2);s:=h.search("стол один",nil);h.search("стол два",nil)
 h.mu.Lock();if detailBlock {h.detailStatus=code}else{h.sourceStatus=code};h.mu.Unlock()
 ctx,cancel:=context.WithCancel(context.Background());done:=make(chan error,1);go func(){done<-h.a.pollLoop(ctx)}();defer func(){cancel();<-done}()
 verifierWait(t,func()bool{return h.a.RuntimeStatus().BackoffUntil.After(time.Now())})
 status:=h.a.RuntimeStatus();if !strings.Contains(status.Mode,fmt.Sprint(code)) {t.Fatal(status)}
 if err:=h.a.RequestPollNow(-200,1);err==nil {t.Fatal("manual bypassed shared block")}
 if len(h.history(s))!=0 {t.Fatal(h.history(s))}
 h.mu.Lock();trace:=append([]keywordTrace(nil),h.trace...);h.mu.Unlock();expected:=2;if detailBlock {expected=3};if len(trace)!=expected {t.Fatal("requests continued after block",trace)}
 t.Logf("AC006 source/detail block%d backoff=%s manual refused trace=%v",code,status.BackoffUntil,trace)
 })}}
 t.Run("running-manual",func(t *testing.T){
 h:=newKeywordPollingHarness(t,1);s:=h.search("стол",nil);entered,release:=make(chan struct{}),make(chan struct{});var once sync.Once;var releaseOnce sync.Once
 h.mu.Lock();h.detailHook=func(){once.Do(func(){close(entered);<-release})};h.mu.Unlock()
 if err:=h.a.RequestPollNow(-200,1);err!=nil {t.Fatal(err)}
 if err:=h.a.RequestPollNow(-200,2);err==nil {t.Fatal("queued duplicate")}
 ctx,cancel:=context.WithCancel(context.Background());done:=make(chan error,1);go func(){done<-h.a.pollLoop(ctx)}();defer func(){releaseOnce.Do(func(){close(release)});cancel();<-done}()
 select {case <-entered:case <-time.After(time.Second):t.Fatal("running barrier timeout")}
 for i:=0;i<3;i++ {if err:=h.a.RequestPollNow(-200,int64(i+3));err==nil {t.Fatal("running duplicate",i)}}
 releaseOnce.Do(func(){close(release)})
 verifierWait(t,func()bool{return !h.a.RuntimeStatus().NextPoll.IsZero()})
 if !h.history(s)[1001].Delivered {t.Fatal("running scan lost delivery")}
 h.mu.Lock();count:=len(h.captions);max:=h.maxInFlight;h.mu.Unlock();if count!=1||max!=1 {t.Fatal(count,max)}
 t.Logf("AC006 queued/running duplicate manual triggers refused; sends%d inFlight%d",count,max)
 })
}

func TestVerifierSearchSendWriteAmbiguity(t *testing.T) {
 h:=newKeywordPollingHarness(t,2);s:=h.search("стол",nil)
 h.mu.Lock();h.sendHook=func(){if err:=h.db.Close();err!=nil {t.Error(err)}};h.mu.Unlock()
 if err:=h.a.pollOnce(context.Background());err==nil {t.Fatal("post-send store error absent")}
 h.mu.Lock();h.sendHook=nil;h.mu.Unlock();h.reopen()
 if len(h.history(s))!=0||len(h.captions)!=1 {t.Fatal("uncommitted success became durable",h.history(s),h.captions)}
 h.poll();h.reopen();h.poll()
 if len(h.captions)!=2||!h.history(s)[1001].Delivered {t.Fatal("accepted repeat/restart lifecycle",h.captions,h.history(s))}
 t.Log("AC006 first successful send lost only history write; one possible repeat, then durable dedup after reopen")
}
