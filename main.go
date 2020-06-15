/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */

package main

import (
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"log"
	"ntc-gnats/npub"
	"ntc-gnats/nreq"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
)

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
