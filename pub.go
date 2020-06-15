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
}
