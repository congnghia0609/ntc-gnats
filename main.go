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
	"ntc-gnats/nsub"
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

func main() {
	// Init NConf
	InitNConf()

	// GetConfig
	c := nconf.GetConfig()
	fmt.Println(c.GetString("notify.pub.url"))
	fmt.Println(c.GetString("notify.pub.auth"))

	// InitSub
	nsub.InitSubConf("chat")

	var poolnsub nsub.PoolNSubscriber
	for i:=0; i<2; i++ {
		ns := nsub.NSubscriber{int32(i), "msg.test"}
		poolnsub.AddNSub(ns)
	}
	poolnsub.RunPoolNSub()


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

	// uuid
	//for i:=0; i<100; i++ {
	//	fmt.Println(strconv.Itoa(i), nutil.GetUUID())
	//}


	// Hang thread Main.
	s := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s, os.Interrupt)
	// Block until we receive our signal.
	<-s
	log.Println("################# End Main #################")
}
