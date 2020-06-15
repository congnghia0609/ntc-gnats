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
	for i := 0; i < 2; i++ {
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
				log.Printf("SimpleSubscriber[#%s] Received on PubSub [%s]: '%s'", ns.Name, ns.Subject, string(msg.Data))
			}
		}
	}()
	fmt.Printf("SimpleSubscriber[#%s] start...\n", ns.Name)
}
