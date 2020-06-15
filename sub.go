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
	"ntc-gnats/nsub"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
)

func InitNConf5() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

/**
 * cd ~/go-projects/src/ntc-gnats
 * go run sub.go
 */
func main() {
	// Init NConf
	InitNConf5()

	//// Start Simple Subscriber
	for i:=0; i<2; i++ {
		StartSimpleSubscriber()
	}

	// Hang thread Main.
	s := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s, os.Interrupt)
	// Block until we receive our signal.
	<-s
	log.Println("################# End Main #################")
}

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
				log.Printf("ChatNSubscriber[#%s] Received on PubSub [%s]: '%s'", ns.Name, ns.Subject, string(msg.Data))
			}
		}
	}()
	fmt.Printf("ChatSubscriber[#%s] start...\n", ns.Name)
}

//func testSub() {
//	// DefaultURL: nats://127.0.0.1:4222
//	var urls = nats.DefaultURL
//	var showTime = true
//
//	// Connect Options.
//	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
//	opts = append(opts, nats.UserInfo("username", "password"))
//	opts = setupConnOptions(opts)
//
//	// Connect to NATS
//	nc, err := nats.Connect(urls, opts...)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//subj, i := "msg.test", 0
//	//nc.Subscribe(subj, func(msg *nats.Msg) {
//	//	i += 1
//	//	printMsg(msg, i)
//	//})
//	go doTask(nc)
//	fmt.Println("=========== Out Subscribe ===========")
//	nc.Flush()
//
//	if err := nc.LastError(); err != nil {
//		log.Fatal(err)
//	}
//
//	//log.Printf("Listening on [%s]", subj)
//	if showTime {
//		log.SetFlags(log.LstdFlags)
//	}
//
//	fmt.Println("=========== Out Subscribe 1 ===========")
//	runtime.Goexit()
//	fmt.Println("=========== Out Subscribe 2 ===========")
//}
//
//func doTask(nc *nats.Conn) {
//	subj, i := "msg.test", 0
//	nc.Subscribe(subj, func(msg *nats.Msg) {
//		i += 1
//		printMsg(msg, i)
//	})
//	log.Printf("Listening on [%s]", subj)
//}
//
//func setupConnOptions(opts []nats.Option) []nats.Option {
//	totalWait := 10 * time.Minute
//	reconnectDelay := time.Second
//
//	opts = append(opts, nats.ReconnectWait(reconnectDelay))
//	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
//	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
//		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
//	}))
//	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
//		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
//	}))
//	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
//		log.Fatalf("Exiting: %v", nc.LastError())
//	}))
//	return opts
//}
//
//func printMsg(m *nats.Msg, i int) {
//	log.Printf("[#%d] Received on PubSub [%s]: '%s'", i, m.Subject, string(m.Data))
//}


