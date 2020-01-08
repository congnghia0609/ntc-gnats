# ntc-gnats
ntc-gnats is module [NAST](https://nats.io/) golang client.  

## Install
```bash
go get -u github.com/congnghia0609/ntc-gnats
```

## Build from source
```bash
// Install dependencies
make deps

// Build
make build

// Clean file build
make clean
```

## Usage
### Publisher
```go
//// InitPub
npub.InitPubConf("notify")
// Case 1: PubSub.
for i:=0; i<10; i++ {
    subj, msg := "msg.test", "hello " + strconv.Itoa(i)
    npub.Publish(subj, msg)
    log.Printf("Published PubSub [%s] : '%s'\n", subj, msg)
}
```

### Subscriber
struct **NSubscriber** and **PoolNSubscriber** reference file [nsub_template.go](https://github.com/congnghia0609/ntc-gnats/blob/master/nsub/nsub_template.go).  
```go
//// InitSub
nsub.InitSubConf("chat")
// Init PoolNSubscriber
var poolnsub nsub.PoolNSubscriber
for i:=0; i<2; i++ {
    ns := nsub.NSubscriber{strconv.Itoa(i), "msg.test", nil, nil}
    poolnsub.AddNSub(ns)
}
poolnsub.RunPoolNSub()
```

### Queue Worker
struct **NWorker** and **PoolNWorker** reference file [nworker_template.go](https://github.com/congnghia0609/ntc-gnats/blob/master/nworker/nworker_template.go).  
```go
//// InitWorker
nworker.InitWorkerConf("email")
// Init PoolNWorker
var poolnworker nworker.PoolNWorker
for i:=0; i<2; i++ {
    nw := nworker.NWorker{strconv.Itoa(i), "worker.email", "worker.email", nil, nil}
    poolnworker.AddNWorker(nw)
}
poolnworker.RunPoolNWorker()
```

### Request
```go
//// InitNReq
nreq.InitReqConf("dbreq")
for i:=0; i<10; i++ {
    subj, payload := "reqres", "this is request " + strconv.Itoa(i)
    msg, err := nreq.Request(subj, payload)
    if err != nil {
        log.Fatalf("%v for request", err)
    }
    log.Printf("Published [%s] : '%s'", subj, payload)
    log.Printf("Received  [%v] : '%s'", msg.Subject, string(msg.Data))
}
```

### Respond
struct **NRes** and **PoolNRes** reference file [nres_template.go](https://github.com/congnghia0609/ntc-gnats/blob/master/nres/nres_template.go).  
```go
// InitNRes
nres.InitResConf("dbres")
// Init PoolNRes
var poolnres nres.PoolNRes
for i:=0; i<2; i++ {
    nrs := nres.NRes{strconv.Itoa(i), "reqres", "dbquery", nil, nil}
    poolnres.AddNRes(nrs)
}
poolnres.RunPoolNRes()
```

## License
This code is under the [Apache Licence v2](https://www.apache.org/licenses/LICENSE-2.0).  
