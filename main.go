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

func GetWDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

func InitNConf() {
	wdir := GetWDir()
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

func main() {
	// Init NConf
	InitNConf()

	// GetConfig
	c := nconf.GetConfig()
	fmt.Println(c.GetString("notify.pub.url"))
	fmt.Println(c.GetString("notify.pub.auth"))

	// InitPub
	npub.InitPubConf("notify")

	// Case 1: PubSub.
	for i:=0; i<10; i++ {
		subj, msg := "msg.test", "hello " + strconv.Itoa(i)
		npub.Publish(subj, msg)
		log.Printf("Published PubSub [%s] : '%s'\n", subj, msg)
	}

	//// Case 2: Queue Group.
	//for i:=0; i<10; i++ {
	//	subj, msg := "job", []byte("hello " + strconv.Itoa(i))
	//	nc.Publish(subj, msg)
	//	log.Printf("Published Queue [%s] : '%s'\n", subj, msg)
	//}

}
