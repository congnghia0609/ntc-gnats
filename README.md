# ntc-gnats
ntc-gnats is module NAST golang client.  

## Install dependencies
```bash
make deps
```

## Build
```bash
make build
```

## Publisher
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

## Subscriber
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

## Queue Worker
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

## Request
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

## Respond
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
