package nsub

import (
	"errors"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
	"time"
)

var surl string
var sauth string
var sopts []nats.Option

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

func InitSubConf(name string) error {
	if name == "" {
		return errors.New("name config is not empty.")
	}
	sopts = []nats.Option{nats.Name("NSubscriber")}
	c := nconf.GetConfig()
	surl = c.GetString(name+".sub.url")
	sauth = c.GetString(name+".sub.auth")
	if len(sauth) > 0 {
		arrauth := strings.Split(sauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			sopts = append(sopts, nats.UserInfo(username, password))
		}
	}
	sopts = setupConnOptions(sopts)
	return nil
}

// url = nats://127.0.0.1:4222
// auth = username:password
func InitSubParams(url string, auth string) error {
	sopts = []nats.Option{nats.Name("NSubscriber")}
	surl = url
	sauth = auth
	if len(sauth) > 0 {
		arrauth := strings.Split(sauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			sopts = append(sopts, nats.UserInfo(username, password))
		}
	}
	sopts = setupConnOptions(sopts)
	return nil
}
