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
	"ntc-gnats/nres"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
)

func InitNConf4() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

/**
 * cd ~/go-projects/src/ntc-gnats
 * go run res.go
 */
func main() {
	// Init NConf
	InitNConf4()

	// Start Simple Response
	for i := 0; i < 2; i++ {
		StartSimpleResponse()
	}

	// Hang thread Main.
	s := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s, os.Interrupt)
	// Block until we receive our signal.
	<-s
	log.Println("################# End Main #################")
}

func StartSimpleResponse() {
	name := "dbres"
	nrs := nres.NewNResponse(name)
	processChan := make(chan *nats.Msg)
	nrs.Start(processChan)
	// go-routine process message.
	go func() {
		reply := "this is response ==> "
		for {
			select {
			case msg := <-processChan:
				// Process message in here.
				log.Printf("NRes[%s][#%s] Received on QueueNRes[%s]: '%s'", nrs.Group, nrs.Name, nrs.Subject, string(msg.Data))
				datares := reply + string(msg.Data)
				msg.Respond([]byte(datares))
				log.Printf("NRes[%s][#%s] Reply on QueueNRes[%s]: '%s'", nrs.Group, nrs.Name, nrs.Subject, datares)
			}
		}
	}()
	fmt.Printf("SimpleResponse[#%s] start...\n", nrs.Name)
}
