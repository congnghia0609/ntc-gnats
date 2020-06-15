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
	"ntc-gnats/nworker"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
)

func InitNConf6() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

/**
 * cd ~/go-projects/src/ntc-gnats
 * go run worker.go
 */
func main() {
	// Init NConf
	InitNConf6()

	//// InitWorker
	//nworker.InitWorkerConf("email")
	//// Init PoolNWorker
	//var poolnworker nworker.PoolNWorker
	//for i:=0; i<2; i++ {
	//	nw := nworker.NWorker{strconv.Itoa(i), "worker.email", "worker.email1", nil, nil}
	//	poolnworker.AddNWorker(nw)
	//}
	//poolnworker.RunPoolNWorker()

	// Start Simple Worker
	for i := 0; i < 2; i++ {
		StartSimpleWorker()
	}

	// Hang thread Main.
	s := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s, os.Interrupt)
	// Block until we receive our signal.
	<-s
	log.Println("################# End Main #################")
}

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
