//go:build cgo

package app

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestVerifierCurrentFeedPollingHistory(t *testing.T) {
	h:=newKeywordPollingHarness(t,10)
	max:=100
	s:=h.search("first",&max)
	other:=h.search("second",nil)
	seed:=func(monitor,id,revision int64, delivered bool) {
		t.Helper(); if err:=h.db.RecordKeywordAdState(monitor,id,revision,delivered);err!=nil {t.Fatal(err)}
	}
	seed(s.ID,1001,s.Revision,true)
	seed(s.ID,1002,s.Revision,false)
	seed(s.ID,2001,s.Revision,true)
	seed(s.ID,2002,s.Revision,false)
	seed(other.ID,1001,other.Revision,false)
	seed(other.ID,2002,other.Revision,true)
	otherBefore:=h.history(other)
	h.mu.Lock();h.ids=[]int64{1001,1002,1003};h.mu.Unlock()
	remaining:=10
	stats,count,err:=h.a.pollKeywordSearch(context.Background(),s,&remaining)
	if err!=nil || count!=3 || stats.NewIDs!=1 || stats.DetailRequests!=0 || stats.Sent!=0 {t.Fatal(stats,count,err)}
	initial:=h.history(s)
	if len(initial)!=5 || !initial[1001].Delivered || !initial[2001].Delivered || initial[1003].EvaluatedRevision!=s.Revision {t.Fatal(initial)}
	max=200
	s.PriceMax=&max
	s,found,err:=h.a.UpdateKeywordSearch(s)
	if err!=nil || !found {t.Fatal(found,err)}
	h.reopen()
	remaining=10
	stats,count,err=h.a.pollKeywordSearch(context.Background(),s,&remaining)
	if err!=nil || count!=3 || stats.Sent!=2 || stats.DetailRequests!=2 {t.Fatal(stats,count,err)}
	final:=h.history(s)
	if len(final)!=5 || !final[1001].Delivered || !final[1002].Delivered || !final[1003].Delivered || final[2001]!=initial[2001] || final[2002]!=initial[2002] {t.Fatal(final)}
	if got:=h.history(other); !reflect.DeepEqual(got,otherBefore) {t.Fatal("other monitor changed",got,otherBefore)}
	h.mu.Lock()
	sends:=len(h.captions)
	for _,r:=range h.trace {if strings.HasPrefix(r.path,"/adv/1001_") {t.Error("delivered ID was fetched again",r)}}
	h.mu.Unlock()
	if sends!=2 {t.Fatal("expected only revision reevaluation deliveries",sends)}
	h.reopen();remaining=10
	stats,_,err=h.a.pollKeywordSearch(context.Background(),s,&remaining)
	if err!=nil || stats.NewIDs!=0 || stats.Sent!=0 || stats.DetailRequests!=0 || !reflect.DeepEqual(h.history(s),final) {t.Fatal("restart did not preserve dedup",stats,err,h.history(s))}
	t.Log("FT-005-AC-009 production App polling preserved 2 off-feed historical states + second monitor; edit/reopen reevaluated only current rejections; delivered current ID skipped; restart skipped every delivered current ID")
}
