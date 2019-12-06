/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"time"
)

/**
 * cd ~/go-projects/src/ntc-gnats
 * go run res.go
 */
func main() {
	// DefaultURL: nats://127.0.0.1:4222
	var urls = nats.DefaultURL
	var showTime = true
	var queueName = "queue-rr"

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Responder")}
	opts = setupConnOptions2(opts)

	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	subj, reply, i := "reqres", "this is response", 0

	nc.QueueSubscribe(subj, queueName, func(msg *nats.Msg) {
		i++
		printMsg2(msg, i)
		msg.Respond([]byte(reply))
		printMsg3(msg, reply, i)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", subj)
	if showTime {
		log.SetFlags(log.LstdFlags)
	}

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")
}

func setupConnOptions2(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

func printMsg2(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
}

func printMsg3(m *nats.Msg, data string, i int) {
	log.Printf("[#%d] Reply on [%s]: '%s'", i, m.Subject, data)
}