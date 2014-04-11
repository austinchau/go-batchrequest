package request

import (
  "testing"
  "log"
)

func GetRequestURI(uri string) *Request{
  return NewRequest(map[string]interface{}{
    "uri": uri,
    "method": "get",
    // "headers": map[string]string{
    //   "a": "b",
    // },
    // "body": "austin=austu",
  })
}

func GetRequest() *Request{
  return NewRequest(map[string]interface{}{
    "uri": "http://www.yahoo.com",
    "method": "get",
    "headers": map[string]string{
      "a": "b",
    },
    // "body": "austin=austu",
  })
}

func TestFetch(t *testing.T) {
 r := GetRequest() 
 status, headers, body, err := r.Fetch()
 log.Println(status)
 log.Println(headers)
 if err != nil {
   t.Error(err)
 }
 if len(body) == 0 {
   t.Error("empty body")
 }
}

func TestBatch(t *testing.T) {
  batch := NewBatch()
  batch.Add(GetRequestURI("http://www.ebay.com"))
  batch.Add(GetRequestURI("http://www.google.com"))
  batch.Add(GetRequestURI("http://www.reputation.com"))
  batch.Add(GetRequestURI("http://www.bing.com"))
  batch.Add(GetRequestURI("http://www.apple.com"))
  batch.Add(GetRequestURI("http://www.facebook.com"))
  for uri,data := range batch.Process(5) {
    m := data.(map[string]interface{})
    log.Println(uri)
    log.Printf("%v", m["status"]) 
    // log.Printf("%v", m["headers"]) 
    log.Printf("%v", m["body"]) 
    log.Printf("%v", m["err"]) 
  }
}