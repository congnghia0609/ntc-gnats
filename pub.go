/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
)

/**
 * cd ~/go-projects/src/ntc-gnats
 * go run pub.go
 */
func main() {
	// DefaultURL: nats://127.0.0.1:4222
	var urls = nats.DefaultURL
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Publisher")}

	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Case 1: PubSub.
	for i:=0; i<10; i++ {
		subj, msg := "msg.test", []byte("hello " + strconv.Itoa(i))
		nc.Publish(subj, msg)
		log.Printf("Published PubSub [%s] : '%s'\n", subj, msg)
	}

	//// Case 2: Queue Group.
	//for i:=0; i<10; i++ {
	//	subj, msg := "job", []byte("hello " + strconv.Itoa(i))
	//	nc.Publish(subj, msg)
	//	log.Printf("Published Queue [%s] : '%s'\n", subj, msg)
	//}


	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
}