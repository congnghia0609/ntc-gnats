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
	"strconv"
	"time"
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

	//// InitNRes
	nres.InitResConf("dbres")
	// Init PoolNRes
	var poolnres nres.PoolNRes
	for i:=0; i<2; i++ {
		nrs := nres.NRes{strconv.Itoa(i), "reqres", "dbquery", nil, nil}
		poolnres.AddNRes(nrs)
	}
	poolnres.RunPoolNRes()

	// Hang thread Main.
	s := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s, os.Interrupt)
	// Block until we receive our signal.
	<-s
	log.Println("################# End Main #################")
}

func testRes() {
	// DefaultURL: nats://127.0.0.1:4222
	var urls = nats.DefaultURL
	var showTime = true
	var queueName = "queue-rr"

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Responder")}
	opts = append(opts, nats.UserInfo("username", "password"))
	opts = setupConnOptions2(opts)

	// Connect to NATS
	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	subj, reply, i := "reqres", "this is response ==> ", 0

	nc.QueueSubscribe(subj, queueName, func(msg *nats.Msg) {
		i++
		printMsg2(msg, i)
		msg.Respond([]byte(reply + string(msg.Data)))
		printMsg3(msg, reply + string(msg.Data), i)
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
