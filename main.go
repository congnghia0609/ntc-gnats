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
	"ntc-gnats/nres"
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
func main() {
	// Init NConf
	InitNConf()

	// GetConfig
	c := nconf.GetConfig()
	fmt.Println(c.GetString("notify.pub.url"))
	fmt.Println(c.GetString("notify.pub.auth"))

	////// InitSub
	//nsub.InitSubConf("chat")
	//// Init PoolNSubscriber
	//var poolnsub nsub.PoolNSubscriber
	//for i:=0; i<2; i++ {
	//	ns := nsub.NSubscriber{strconv.Itoa(i), "msg.test", nil, nil}
	//	poolnsub.AddNSub(ns)
	//}
	//poolnsub.RunPoolNSub()

	//// InitWorker
	//nworker.InitWorkerConf("email")
	//// Init PoolNWorker
	//var poolnworker nworker.PoolNWorker
	//for i:=0; i<2; i++ {
	//	nw := nworker.NWorker{strconv.Itoa(i), "worker.email", "worker.email", nil, nil}
	//	poolnworker.AddNWorker(nw)
	//}
	//poolnworker.RunPoolNWorker()

	// InitNRes
	nres.InitResConf("proccess")
	// Init PoolNRes
	var poolnres nres.PoolNRes
	for i:=0; i<2; i++ {
		nrs := nres.NRes{strconv.Itoa(i), "reqres", "proccess", nil, nil}
		poolnres.AddNRes(nrs)
	}
	poolnres.RunPoolNRes()



	////// InitPub
	//npub.InitPubConf("notify")
	//// Case 1: PubSub.
	//for i:=0; i<10; i++ {
	//	subj, msg := "msg.test", "hello " + strconv.Itoa(i)
	//	npub.Publish(subj, msg)
	//	log.Printf("Published PubSub [%s] : '%s'\n", subj, msg)
	//}

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
