/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package nworker

import (
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"math"
	"ntc-gnats/nutil"
	"strings"
)

type NWorker struct {
	Name    string
	Subject string
	Group   string
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

func NewNWorker(name string) *NWorker {
	if len(name) == 0 {
		return nil
	}
	c := nconf.GetConfig()
	id := name + "_NWorker_" + nutil.GetGUUID()
	wopts := []nats.Option{nats.Name(id), nats.MaxReconnects(math.MaxInt32)}
	subject := c.GetString(name + ".worker.subject")
	group := c.GetString(name + ".worker.group")
	wurl := c.GetString(name + ".worker.url")
	log.Printf("wurl=%s", wurl)
	wauth := c.GetString(name + ".worker.auth")
	if len(wauth) > 0 {
		arrauth := strings.Split(wauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			wopts = append(wopts, nats.UserInfo(username, password))
		}
	}
	wopts = setupConnOptions(wopts)
	nc, err := nats.Connect(wurl, wopts...)
	if err != nil {
		log.Println(err)
	}
	return &NWorker{Name: id, Subject: subject, Group: group, Url: wurl, Auth: wauth, Opts: wopts, Conn: nc, Subt: nil}
}

func (nw *NWorker) Start(processChan chan *nats.Msg) error {
	// Connect to NATS
	var err error
	nw.Subt, err = nw.Conn.QueueSubscribe(nw.Subject, nw.Group, func(msg *nats.Msg) {
		//log.Printf("NWorker[%s][#%s] Received on QueueWorker[%s]: '%s'", nw.Group, nw.Name, nw.Subject, string(msg.Data))
		processChan <- msg
	})
	nw.Conn.Flush()
	if err = nw.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("NWorker[%s][#%s] is listening on Subject[%s]", nw.Group, nw.Name, nw.Subject)
	return err
}

// Disconnect safe
func (nw *NWorker) Stop() {
	if nw != nil && nw.Conn != nil {
		nw.Conn.Drain()
	}
}
