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
	"ntc-gnats/nreq"
	"path/filepath"
	"runtime"
	"strconv"
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

	//// InitNReq
	////// Cach 1.1.
	//name := "dbreq"
	//for i:=0; i<10; i++ {
	//	subj, payload := "reqres", "this is request " + strconv.Itoa(i)
	//	msg, err := nreq.Request(name, subj, payload)
	//	if err != nil {
	//		log.Fatalf("%v for request", err)
	//	}
	//	log.Printf("NReq Published [%s] : '%s'", subj, payload)
	//	log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	//}

	////// Cach 1.2.
	name := "dbreq"
	nr := nreq.GetInstance(name)
	for i := 0; i < 10; i++ {
		subj, payload := "reqres", "this is request "+strconv.Itoa(i)
		msg, err := nr.Request(subj, payload)
		if err != nil {
			log.Fatalf("%v for request", err)
		}
		log.Printf("NReq[%s] Published [%s] : '%s'", nr.Name, subj, payload)
		log.Printf("NReq[%s] Received  [%v] : '%s'", nr.Name, msg.Subject, string(msg.Data))
	}
}

//func testReq() {
//	// DefaultURL: nats://127.0.0.1:4222
//	var urls = nats.DefaultURL
//
//	// Connect Options.
//	opts := []nats.Option{nats.Name("NATS Sample Requestor")}
//	opts = append(opts, nats.UserInfo("username", "password"))
//
//	// Connect to NATS
//	nc, err := nats.Connect(urls, opts...)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer nc.Close()
//	for i:=0; i<10; i++ {
//		subj, payload := "reqres", []byte("this is request " + strconv.Itoa(i))
//		msg, err := nc.Request(subj, payload, 10*time.Second)
//		if err != nil {
//			if nc.LastError() != nil {
//				log.Fatalf("%v for request", nc.LastError())
//			}
//			log.Fatalf("%v for request", err)
//		}
//
//		log.Printf("Published [%s] : '%s'", subj, payload)
//		log.Printf("Received  [%v] : '%s'", msg.Subject, string(msg.Data))
//	}
//}
