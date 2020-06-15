/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package nsub

import (
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"math"
	"ntc-gnats/nutil"
	"strings"
)

type NSubscriber struct {
	Name    string
	Subject string
	Url     string
	Auth    string
	Opts    []nats.Option
	Conn    *nats.Conn
	Subt    *nats.Subscription
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	opts = append(opts, nats.ReconnectWait(nats.DefaultReconnectWait))
	opts = append(opts, nats.MaxReconnects(math.MaxInt32))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %v times", err, math.MaxInt32)
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

func NewNSubscriber(name string) *NSubscriber {
	if len(name) == 0 {
		return nil
	}
	c := nconf.GetConfig()
	id := name + "_NSubscriber_" + nutil.GetGUUID()
	sopts := []nats.Option{nats.Name(id), nats.MaxReconnects(math.MaxInt32)}
	subject := c.GetString(name + ".sub.subject")
	surl := c.GetString(name + ".sub.url")
	log.Printf("surl=%s", surl)
	sauth := c.GetString(name + ".sub.auth")
	if len(sauth) > 0 {
		arrauth := strings.Split(sauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			sopts = append(sopts, nats.UserInfo(username, password))
		}
	}
	sopts = setupConnOptions(sopts)
	nc, err := nats.Connect(surl, sopts...)
	if err != nil {
		log.Println(err)
	}
	return &NSubscriber{Name: id, Subject: subject, Url: surl, Auth: sauth, Opts: sopts, Conn: nc, Subt: nil}
}

func (ns *NSubscriber) Start(processChan chan *nats.Msg) error {
	// Subscribe to NATS
	var err error
	ns.Subt, err = ns.Conn.Subscribe(ns.Subject, func(msg *nats.Msg) {
		//log.Printf("NSubscriber[#%s] Received on PubSub [%s]: '%s'", ns.Name, ns.Subject, string(msg.Data))
		processChan <- msg
	})
	ns.Conn.Flush()
	if err = ns.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("NSubscriber[#%s] is listening on Subject[%s]", ns.Name, ns.Subject)
	return err
}

// Disconnect safe
func (ns *NSubscriber) Stop() {
	if ns != nil && ns.Conn != nil {
		ns.Conn.Drain()
	}
}
