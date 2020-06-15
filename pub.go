/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package main

import (
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"ntc-gnats/npub"
	"path/filepath"
	"runtime"
	"strconv"
)

func InitNConf3() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

/**
 * cd ~/go-projects/src/ntc-gnats
 * go run pub.go
 */
func main() {
	// Init NConf
	InitNConf3()

	////// Publish
	//// Case 1: PubSub.
	////// Cach 1.1.
	//name := "notify"
	//for i:=0; i<10; i++ {
	//	subj, msg := "msg.test", "hello " + strconv.Itoa(i)
	//	npub.Publish(name, subj, msg)
	//	log.Printf("Published PubSub[%s] : '%s'\n", subj, msg)
	//}
	////// Cach 1.2.
	name := "notify"
	np := npub.GetInstance(name)
	for i:=0; i<10; i++ {
		subj, msg := "msg.test", "hello " + strconv.Itoa(i)
		np.Publish(subj, msg)
		log.Printf("Published PubSub[%s] : '%s'\n", subj, msg)
	}


	//// Case 2: Queue Group.
	//npub.InitPubConf("notify")
	//for i:=0; i<10; i++ {
	//	subj, msg := "worker.email", "hello " + strconv.Itoa(i)
	//	npub.Publish(subj, msg)
	//	log.Printf("Published QueueWorker[%s] : '%s'\n", subj, msg)
	//}
}

func testPub() {
	// DefaultURL: nats://127.0.0.1:4222
	var urls = nats.DefaultURL
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Publisher")}
	opts = append(opts, nats.UserInfo("username", "password"))

	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	//// Case 1: PubSub.
	//for i:=0; i<10; i++ {
	//	subj, msg := "msg.test", []byte("hello " + strconv.Itoa(i))
	//	nc.Publish(subj, msg)
	//	log.Printf("Published PubSub [%s] : '%s'\n", subj, msg)
	//}

	// Case 2: Queue Group.
	for i:=0; i<10; i++ {
		subj, msg := "worker.email", []byte("hello " + strconv.Itoa(i))
		nc.Publish(subj, msg)
		log.Printf("Published QueueWorker[%s] : '%s'\n", subj, msg)
	}

	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
}
