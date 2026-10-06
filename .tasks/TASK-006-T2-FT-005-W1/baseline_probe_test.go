package somon
import("context";"errors";"net/http";"net/http/httptest";"os";"reflect";"testing";"time")
func TestKeywordSourceBaseline(t *testing.T){
 for _,tc:=range []struct{name string;want []int64}{{"small",[]int64{21000001}},{"empty",[]int64{}}}{t.Run(tc.name,func(t *testing.T){
 body,err:=os.ReadFile("../../testdata/keyword-search-primary-"+tc.name+".html");if err!=nil{t.Fatal(err)}
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/search/vose/"||r.URL.Query().Get("q")!="стол & стул"||r.URL.Query().Get("ordering")!="relevance"{t.Errorf("request=%s",r.URL)};w.Write(body)}));defer server.Close()
 cards,_,err:=NewClient("test",0,time.Second,1<<20).FetchCategory(context.Background(),server.URL+"/search/vose/?q=%D1%81%D1%82%D0%BE%D0%BB+%26+%D1%81%D1%82%D1%83%D0%BB&ordering=relevance")
 if err!=nil{t.Fatalf("expected primary success: %v",err)};ids:=[]int64{};for _,c:=range cards{ids=append(ids,c.ID)};if !reflect.DeepEqual(ids,tc.want){t.Fatalf("primary IDs=%v want %v; foreign-region contamination",ids,tc.want)}
 })}
 t.Run("blocked",func(t *testing.T){_,err:=ParseCategory(DefaultCategoryURL,[]byte(`<h1>Access denied</h1>`));var blocked *BlockedPageError;if !errors.As(err,&blocked){t.Fatal(err)}})
 t.Run("malformed",func(t *testing.T){_,err:=ParseCategory(DefaultCategoryURL,[]byte(`<h1>Broken page</h1>`));if err==nil{t.Fatal("missing parse error")}})
 t.Run("transport",func(t *testing.T){s:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.WriteHeader(403)}));defer s.Close();_,_,err:=NewClient("test",0,time.Second,1<<20).FetchCategory(context.Background(),s.URL);var httpErr *HTTPError;if !errors.As(err,&httpErr)||httpErr.StatusCode!=403{t.Fatal(err)}})
}
