# ntc-gnats
ntc-gnats is a module [NATS](https://nats.io/) golang client.  

## Install
```bash
go get -u github.com/congnghia0609/ntc-gnats
```

## 1. Publish-Subscribe
### Publisher
```go
////// Publish
//// Case 1: PubSub.
//// Cach 1.1.
name := "notify"
subj := "msg.test"
for i:=0; i<10; i++ {
    msg := "hello " + strconv.Itoa(i)
    npub.Publish(name, subj, msg)
    log.Printf("Published PubSub[%s] : '%s'\n", subj, msg)
}

//// Cach 1.2.
name := "notify"
subj := "msg.test"
np := npub.GetInstance(name)
for i := 0; i < 10; i++ {
    msg := "hello "+strconv.Itoa(i)
    np.Publish(subj, msg)
    log.Printf("Published PubSub[%s] : '%s'\n", subj, msg)
}
```

### Subscriber  
```go
func StartSimpleSubscriber() {
	name := "chat"
	ns := nsub.NewNSubscriber(name)
	processChan := make(chan *nats.Msg)
	ns.Start(processChan)
	// go-routine process message.
	go func() {
		for {
			select {
			case msg := <-processChan:
				// Process message in here.
				log.Printf("SimpleSubscriber[#%s] Received on PubSub [%s]: '%s'", ns.Name, ns.Subject, string(msg.Data))
			}
		}
	}()
	fmt.Printf("SimpleSubscriber[#%s] start...\n", ns.Name)
}
```

## 2. Queue Groups
### Queue Worker  
```go
func StartSimpleWorker() {
	name := "email"
	nw := nworker.NewNWorker(name)
	processChan := make(chan *nats.Msg)
	nw.Start(processChan)
	// go-routine process message.
	go func() {
		for {
			select {
			case msg := <-processChan:
				// Process message in here.
				log.Printf("SimpleNWorker[%s][#%s] Received on QueueWorker[%s]: '%s'", nw.Group, nw.Name, nw.Subject, string(msg.Data))
			}
		}
	}()
	fmt.Printf("SimpleNWorker[#%s] start...\n", nw.Name)
}
```

## 3. Request-Reply
### Request
```go
////// Request
//// Cach 1.
name := "dbreq"
subj := "reqres"
for i:=0; i<10; i++ {
    payload := "this is request " + strconv.Itoa(i)
    msg, err := nreq.Request(name, subj, payload)
    if err != nil {
        log.Fatalf("%v for request", err)
    }
    log.Printf("NReq Published [%s] : '%s'", subj, payload)
    log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
}

//// Cach 2.
name := "dbreq"
subj := "reqres"
nr := nreq.GetInstance(name)
for i := 0; i < 10; i++ {
    payload := "this is request "+strconv.Itoa(i)
    msg, err := nr.Request(subj, payload)
    if err != nil {
        log.Fatalf("%v for request", err)
    }
    log.Printf("NReq[%s] Published [%s] : '%s'", nr.Name, subj, payload)
    log.Printf("NReq[%s] Received  [%v] : '%s'", nr.Name, msg.Subject, string(msg.Data))
}
```

### Reply
```go
func StartSimpleResponse() {
	name := "dbres"
	nrs := nres.NewNResponse(name)
	processChan := make(chan *nats.Msg)
	nrs.Start(processChan)
	// go-routine process message.
	go func() {
		reply := "this is response ==> "
		for {
			select {
			case msg := <-processChan:
				// Process message in here.
				log.Printf("NRes[%s][#%s] Received on QueueNRes[%s]: '%s'", nrs.Group, nrs.Name, nrs.Subject, string(msg.Data))
				datares := reply + string(msg.Data)
				msg.Respond([]byte(datares))
				log.Printf("NRes[%s][#%s] Reply on QueueNRes[%s]: '%s'", nrs.Group, nrs.Name, nrs.Subject, datares)
			}
		}
	}()
	fmt.Printf("SimpleResponse[#%s] start...\n", nrs.Name)
}
```

### File main.go
```go
func InitNConf() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

/** https://github.com/nats-io/nats.go */
/**
 * cd ~/go-projects/src/ntc-gnats
 * go run main.go
 */
func main() {
	// Init NConf
	InitNConf()

	//// Start Simple Subscriber
	for i := 0; i < 2; i++ {
		StartSimpleSubscriber()
	}

	// Start Simple Worker
	for i := 0; i < 2; i++ {
		StartSimpleWorker()
	}

	// Start Simple Response
	for i := 0; i < 2; i++ {
		StartSimpleResponse()
	}

	////// Publish
	//// Case 1: PubSub.
	////// Cach 1.1.
	//name := "notify"
	//subj := "msg.test"
	//for i:=0; i<10; i++ {
	//	msg := "hello " + strconv.Itoa(i)
	//	npub.Publish(name, subj, msg)
	//	log.Printf("Published PubSub[%s] : '%s'\n", subj, msg)
	//}
	//// Cach 1.2.
	name := "notify"
	subj := "msg.test"
	np := npub.GetInstance(name)
	for i := 0; i < 10; i++ {
		msg := "hello " + strconv.Itoa(i)
		np.Publish(subj, msg)
		log.Printf("Published PubSub[%s] : '%s'\n", subj, msg)
	}

	//// Case 2: Queue Group.
	namew := "notify"
	subjw := "worker.email"
	npw := npub.GetInstance(namew)
	for i := 0; i < 10; i++ {
		msg := "hello " + strconv.Itoa(i)
		npw.Publish(subjw, msg)
		log.Printf("Published QueueWorker[%s] : '%s'\n", subjw, msg)
	}

	////// Request
	////// Cach 1.
	//name := "dbreq"
	//subj := "reqres"
	//for i:=0; i<10; i++ {
	//	payload := "this is request " + strconv.Itoa(i)
	//	msg, err := nreq.Request(name, subj, payload)
	//	if err != nil {
	//		log.Fatalf("%v for request", err)
	//	}
	//	log.Printf("NReq Published [%s] : '%s'", subj, payload)
	//	log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	//}
	//// Cach 2.
	namer := "dbreq"
	subjr := "reqres"
	nr := nreq.GetInstance(namer)
	for i := 0; i < 10; i++ {
		payload := "this is request " + strconv.Itoa(i)
		msg, err := nr.Request(subjr, payload)
		if err != nil {
			log.Fatalf("%v for request", err)
		}
		log.Printf("NReq[%s] Published [%s] : '%s'", nr.Name, subjr, payload)
		log.Printf("NReq[%s] Received  [%v] : '%s'", nr.Name, msg.Subject, string(msg.Data))
	}

	// Hang thread Main.
	s := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s, os.Interrupt)
	// Block until we receive our signal.
	<-s
	log.Println("################# End Main #################")
}
```

## License
This code is under the [Apache License v2](https://www.apache.org/licenses/LICENSE-2.0).  
