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
	"ntc-gnats/nreq"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

func InitNConf2() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

/**
 * cd ~/go-projects/src/ntc-gnats
 * go run req.go
 */
func main() {
	// Init NConf
	InitNConf2()

	// InitNReq
	nreq.InitReqConf("db")

	for i:=0; i<10; i++ {
		subj, payload := "reqres", "this is request " + strconv.Itoa(i)
		msg, err := nreq.Request(subj, payload)
		if err != nil {
			log.Fatalf("%v for request", err)
		}
		log.Printf("Published [%s] : '%s'", subj, payload)
		log.Printf("Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	}
}

func test1() {
	// DefaultURL: nats://127.0.0.1:4222
	var urls = nats.DefaultURL

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Requestor")}
	opts = append(opts, nats.UserInfo("username", "password"))

	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	for i:=0; i<10; i++ {
		subj, payload := "reqres", []byte("this is request " + strconv.Itoa(i))
		msg, err := nc.Request(subj, payload, 10*time.Second)
		if err != nil {
			if nc.LastError() != nil {
				log.Fatalf("%v for request", nc.LastError())
			}
			log.Fatalf("%v for request", err)
		}

		log.Printf("Published [%s] : '%s'", subj, payload)
		log.Printf("Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	}
}
