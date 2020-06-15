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

	////// Request
	////// Cach 1.
	//name := "dbreq"
	//subj := "reqres"
	//for i:=0; i<10; i++ {
	//	payload := "this is request " + strconv.Itoa(i)
	//	msg, err := nreq.Request(name, subj, payload)
	//	if err != nil {
	//		log.Fatalf("%v for request", err)
	//	}
	//	log.Printf("NReq Published [%s] : '%s'", subj, payload)
	//	log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	//}

	//// Cach 2.
	name := "dbreq"
	subj := "reqres"
	nr := nreq.GetInstance(name)
	for i := 0; i < 10; i++ {
		payload := "this is request "+strconv.Itoa(i)
		msg, err := nr.Request(subj, payload)
		if err != nil {
			log.Fatalf("%v for request", err)
		}
		log.Printf("NReq[%s] Published [%s] : '%s'", nr.Name, subj, payload)
		log.Printf("NReq[%s] Received  [%v] : '%s'", nr.Name, msg.Subject, string(msg.Data))
	}
}
