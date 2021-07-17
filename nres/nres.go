/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package nres

import (
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/congnghia0609/ntc-gnats/nutil"
	"github.com/nats-io/nats.go"
	"log"
	"math"
	"strings"
)

type NResponse struct {
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

func NewNResponse(name string) *NResponse {
	if len(name) == 0 {
		return nil
	}
	c := nconf.GetConfig()
	id := name + "_NResponse_" + nutil.GetGUUID()
	rsopts := []nats.Option{nats.Name(id), nats.MaxReconnects(math.MaxInt32)}
	subject := c.GetString(name + ".res.subject")
	group := c.GetString(name + ".res.group")
	rsurl := c.GetString(name + ".res.url")
	log.Printf("rsurl=%s", rsurl)
	rsauth := c.GetString(name + ".res.auth")
	if len(rsauth) > 0 {
		arrauth := strings.Split(rsauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			rsopts = append(rsopts, nats.UserInfo(username, password))
		}
	}
	rsopts = setupConnOptions(rsopts)
	nc, err := nats.Connect(rsurl, rsopts...)
	if err != nil {
		log.Println(err)
	}
	return &NResponse{Name: id, Subject: subject, Group: group, Url: rsurl, Auth: rsauth, Opts: rsopts, Conn: nc, Subt: nil}
}

func (nrs *NResponse) Start(processChan chan *nats.Msg) error {
	// Connect to NATS
	var err error
	nrs.Subt, err = nrs.Conn.QueueSubscribe(nrs.Subject, nrs.Group, func(msg *nats.Msg) {
		processChan <- msg
	})
	nrs.Conn.Flush()
	if err = nrs.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("NResponse[%s][#%s] is listening on Subject[%s]", nrs.Group, nrs.Name, nrs.Subject)
	return err
}

// Disconnect safe
func (nrs *NResponse) Stop() {
	if nrs != nil && nrs.Conn != nil {
		nrs.Conn.Drain()
	}
}
