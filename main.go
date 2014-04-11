package request

import (
  "log"
  "net/http"
  "io/ioutil"
  "strings"
  // "sync"
)

type Request struct {
  Uri string
  Method string
  Body string
  Headers map[string]string
}

type Batch struct {
  Requests []*Request
}

func NewRequest(options map[string]interface{}) *Request{
  r := new(Request)
  for k,v := range options {
    switch k {
    case "uri":
      r.Uri = v.(string)
    case "method":
      r.Method = strings.ToUpper(v.(string))
    case "headers":
      r.Headers = v.(map[string]string)
    case "body":
      r.Body = v.(string)
    }
  }
  return r
}

func NewBatch() *Batch {
  batch := new(Batch)
  batch.Requests = make([]*Request, 0, 10)
  return batch
}

func (self *Request) Fetch() (string, http.Header, string, error) {
  client := new(http.Client)
  req, err := http.NewRequest(self.Method, self.Uri, nil)
  if err != nil {
    return "", nil, "", err
  }
  if self.Body != "" {
    req.Body = ioutil.NopCloser(strings.NewReader(self.Body))
  }
  
  for k,v := range self.Headers {
    req.Header.Add(k,v)
  }
  resp, err := client.Do(req)
  if err != nil {
    return "", nil, "", err
  }
  
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return resp.Status, resp.Header, "", err
  }
  // log.Println(resp.Header)
  log.Printf("url=%s size=%d", self.Uri, len(body))
  return resp.Status, resp.Header, string(body), err
}

func (self *Batch) Add(request *Request) {
  self.Requests = append(self.Requests, request)
}

func (self *Batch) Process(maxConcurrent int) map[string]interface{} {
  m := make(map[string]interface{})
  done := make(chan map[string]interface{})

  go func(){
    sema := make(chan int, maxConcurrent)

    for _, request := range self.Requests {
      sema<-1
      go func(r *Request) {
        log.Printf("processing %s", r.Uri)
        status, headers, body, err := r.Fetch()
        if err != nil {
          log.Println(err)
        }
        log.Printf("finished url=%s status=%s size=%d", r.Uri, status, len(body))

        <-sema

        ret := make(map[string]interface{})
        ret["uri"] = r.Uri
        ret["status"] = status
        ret["body"] = body
        ret["headers"] = headers
        ret["err"] = err
        done<-ret
      }(request)
    }
  }()
  
  count := 0
  for {
    ret := <-done
    count++
    m[ret["uri"].(string)] = ret
    if count == len(self.Requests) {
      break
    }
  }   

  log.Println("DONE!") 
  return m
}