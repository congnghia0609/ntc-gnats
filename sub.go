/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
	"github.com/nats-io/nats.go"
)

/**
 * cd ~/go-projects/src/ntc-gnats
 * go run sub.go
 */
func main() {
	// DefaultURL: nats://127.0.0.1:4222
	var urls = nats.DefaultURL
	var showTime = true

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
	opts = append(opts, nats.UserInfo("username", "password"))
	opts = setupConnOptions(opts)

	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	subj, i := "msg.test", 0

	nc.Subscribe(subj, func(msg *nats.Msg) {
		i += 1
		printMsg(msg, i)
	})
	fmt.Println("=========== Out Subscribe ===========")
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", subj)
	if showTime {
		log.SetFlags(log.LstdFlags)
	}

	fmt.Println("=========== Out Subscribe 1 ===========")
	runtime.Goexit()
	fmt.Println("=========== Out Subscribe 2 ===========")
}

func setupConnOptions(opts []nats.Option) []nats.Option {
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

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on PubSub [%s]: '%s'", i, m.Subject, string(m.Data))
}


