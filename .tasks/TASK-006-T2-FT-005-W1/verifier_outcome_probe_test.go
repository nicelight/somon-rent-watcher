package somon

import (
 "context"
 "errors"
 "net/http"
 "net/http/httptest"
 "net/url"
 "os"
 "reflect"
 "strings"
 "testing"
 "time"

 "github.com/nicelight/somon-rent-watcher/internal/model"
)

type reviewer006Transport struct { endpoint *url.URL; observe func(*http.Request) }
func (r reviewer006Transport) RoundTrip(req *http.Request) (*http.Response,error) {
 r.observe(req)
 clone := req.Clone(req.Context())
 clone.URL.Scheme,clone.URL.Host=r.endpoint.Scheme,r.endpoint.Host
 return http.DefaultTransport.RoundTrip(clone)
}
func reviewer006IDs(cards []model.Card) []int64 {
 result:=[]int64{}
 for _,card:=range cards {result=append(result,card.ID)}
 return result
}

func TestReviewer006Outcome(t *testing.T) {
 // Expectations derive from the accepted boundary map, not the catalog function.
 paths:=[]struct{key,path string}{{"all","/search/"},{"furniture","/vse-dlya-doma/mebel/"},{"tables_chairs","/vse-dlya-doma/mebel/mebel-dlya-kuhni/stolyi-stulya/"},{"services","/biznes-i-uslugi/"}}
 cities:=[]struct{key,slug string}{{"country",""},{"dushanbe","dushanbe/"},{"vose","vose/"},{"dangara","dangara/"}}
 for _,category:=range paths {
  for _,city:=range cities {
   t.Run("native/"+category.key+"/"+city.key,func(t *testing.T){
    // Two source primary entries; phrase deliberately differs from title. Nested
    // region boundary, duplicate link, absent fields and foreign RSC isolate AC003.
    body:=[]byte(`<html><body><h1>Поиск 2</h1><div><article><a href="/adv/26000001_service/">Ремонт стула</a><a href="/adv/26000001_service/">Ремонт стула</a><span class="price">0 сомони</span><p>Восе</p><img src="/actual.jpg"></article><article><a href="/adv/26000002_unknown/">Без данных</a></article><section><h2>Объявления <em>из других регионов</em></h2><article><a href="/adv/26000003_foreign/">Чужой стол</a><p>Душанбе</p></article></section></div><script>self.__next_f.push([1,"0:{\"adverts\":[{\"id\":26000003,\"url\":\"/adv/26000003_foreign/\",\"title\":\"Чужой стол\"}]}"])</script></body></html>`)
    var requests int
    server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Write(body)}));defer server.Close()
    endpoint,_:=url.Parse(server.URL)
    client:=NewClient("reviewer006",0,time.Second,1<<20)
    phrase:="  стол + & стул / ? # %  "
    client.httpClient.Transport=reviewer006Transport{endpoint,func(r *http.Request){
     requests++
     if r.URL.Scheme!="https" || r.URL.Host!="somon.tj" || r.URL.Path!=category.path+city.slug || r.URL.Query().Get("q")!=strings.TrimSpace(phrase) || r.URL.Query().Get("ordering")!="relevance" || len(r.URL.Query())!=2 || r.URL.Fragment!="" || r.Method!="GET" || r.Header.Get("User-Agent")!="reviewer006" {t.Errorf("native request changed: %s",r.URL)}
    }}
    cards,raw,err:=client.FetchKeywordSearch(context.Background(),phrase,category.key,city.key)
    if err!=nil || requests!=1 || string(raw)!=string(body) || !reflect.DeepEqual(reviewer006IDs(cards),[]int64{26000001,26000002}) {t.Fatalf("primary outcome ids=%v requests=%d err=%v",reviewer006IDs(cards),requests,err)}
    first,missing:=cards[0],cards[1]
    if first.Price==nil || *first.Price!=0 || first.Currency!="TJS" || first.City!="Восе" || first.ImageURL!="https://somon.tj/actual.jpg" || first.Position!=0 || first.Rooms!=nil || missing.Price!=nil || missing.Currency!="" || missing.City!="" || missing.ImageURL!="" || missing.Position!=1 {t.Fatalf("source fields changed: %+v",cards)}
    t.Logf("native path=%s q=%q ordering=relevance IDs=%v; no local matching/apartment predicate or foreign resurrection",category.path+city.slug,strings.TrimSpace(phrase),reviewer006IDs(cards))
   })
  }
 }
 for _,fixture:=range []struct{name string;ids []int64}{{"small",[]int64{21000001}},{"empty",[]int64{}}} {
  t.Run("unchanged-fixture/"+fixture.name,func(t *testing.T){
   body,err:=os.ReadFile("../../testdata/keyword-search-primary-"+fixture.name+".html");if err!=nil{t.Fatal(err)}
   cards,err:=ParseKeywordSearch("https://somon.tj/search/vose/?q=стол",body)
   if err!=nil || !reflect.DeepEqual(reviewer006IDs(cards),fixture.ids){t.Fatalf("ids=%v want=%v err=%v",reviewer006IDs(cards),fixture.ids,err)}
   t.Logf("same prospective fixture IDs=%v; error=%v",reviewer006IDs(cards),err)
  })
 }
 cases:=[]struct{name,body string; status int;wantError,blocked bool;ids []int64}{
  {"confirmed-zero",`<h1>Потери и находки Дангара 0</h1><div><h2>Объявления из других регионов</h2><a href="/adv/26000003_foreign/">Чужой стол</a></div>`,200,false,false,[]int64{}},
  {"unconfirmed",`<h1>Поиск стол</h1><h2>Объявления из других регионов</h2><a href="/adv/26000003_foreign/">Чужой стол</a>`,200,true,false,[]int64{}},
  {"malformed","not a result page",200,true,false,[]int64{}},
  {"contradictory-zero",`<h1>Поиск 0</h1><article><a href="/adv/26000001_x/">Стол</a></article>`,200,true,false,[]int64{}},
  {"unscoped-rsc",`<script>self.__next_f.push([1,"0:{\"adverts\":[{\"id\":26000003,\"url\":\"/adv/26000003_foreign/\",\"title\":\"Чужой стол\"}]}"])</script>`,200,true,false,[]int64{}},
  {"blocked-html",`<h1>Access denied</h1><h1>Поиск 0</h1>`,200,true,true,[]int64{}},
  {"http403","",403,true,true,[]int64{}},
  {"http429","",429,true,true,[]int64{}},
 }
 for _,tc:=range cases {t.Run("classification/"+tc.name,func(t *testing.T){
  server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Header().Set("Retry-After","90");w.WriteHeader(tc.status);w.Write([]byte(tc.body))}));defer server.Close()
  endpoint,_:=url.Parse(server.URL)
  client:=NewClient("reviewer006",0,time.Second,1<<20)
  client.httpClient.Transport=reviewer006Transport{endpoint,func(*http.Request){}}
  cards,raw,err:=client.FetchKeywordSearch(context.Background(),"стол","all","dangara")
  _,retry,blocked:=IsBlocked(err)
  if (err!=nil)!=tc.wantError || blocked!=tc.blocked || !reflect.DeepEqual(reviewer006IDs(cards),tc.ids) || tc.status==429 && retry!=90*time.Second || tc.status==200 && string(raw)!=tc.body {t.Fatalf("classification cards=%v error=%v blocked=%v retry=%v raw=%q",cards,err,blocked,retry,raw)}
  if tc.status==403 || tc.status==429 {var httpErr *HTTPError;if !errors.As(err,&httpErr)||httpErr.StatusCode!=tc.status{t.Fatal("typed HTTP error lost")}}
  t.Logf("status=%d IDs=%v error=%v blocked=%v retry=%v",tc.status,reviewer006IDs(cards),err!=nil,blocked,retry)
 })}
 t.Run("foreign-currency-not-invented",func(t *testing.T){
  cards,err:=ParseKeywordSearch("https://somon.tj/search/",[]byte(`<article><a href="/adv/26000001_x/">Стол</a><meta itemprop="price" content="200"><meta itemprop="priceCurrency" content="USD"></article>`))
  if err!=nil || len(cards)!=1 || cards[0].Price!=nil || cards[0].Currency!="USD" || cards[0].City!="" || cards[0].ImageURL!="" {t.Fatalf("foreign or missing fields invented: %+v err=%v",cards,err)}
 })
 t.Run("invalid-before-http",func(t *testing.T){
  client:=NewClient("reviewer006",0,time.Second,1<<20)
  client.httpClient.Transport=reviewer006Transport{nil,func(*http.Request){t.Fatal("invalid scope reached HTTP")}}
  for _,values:=range [][3]string{{" ","all","country"},{"стол","../","vose"},{"стол","all","https://evil.test/"}} {if _,_,err:=client.FetchKeywordSearch(context.Background(),values[0],values[1],values[2]);err==nil{t.Fatalf("invalid=%v accepted",values)}}
 })
}
