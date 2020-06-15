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
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
)

func GetWDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

func InitNConf() {
	wdir := GetWDir()
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

	// GetConfig
	c := nconf.GetConfig()
	fmt.Println(c.GetString("notify.pub.url"))
	fmt.Println(c.GetString("notify.pub.auth"))

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
	////// Cach 1.2.
	//name := "notify"
	//subj := "msg.test"
	//np := npub.GetInstance(name)
	//for i := 0; i < 10; i++ {
	//	msg := "hello "+strconv.Itoa(i)
	//	np.Publish(subj, msg)
	//	log.Printf("Published PubSub[%s] : '%s'\n", subj, msg)
	//}

	//// Case 2: Queue Group.
	name := "notify"
	subj := "worker.email"
	np := npub.GetInstance(name)
	for i := 0; i < 10; i++ {
		msg := "hello " + strconv.Itoa(i)
		np.Publish(subj, msg)
		log.Printf("Published QueueWorker[%s] : '%s'\n", subj, msg)
	}

	// Hang thread Main.
	s := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s, os.Interrupt)
	// Block until we receive our signal.
	<-s
	log.Println("################# End Main #################")
}
