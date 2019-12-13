package nsub

import (
	"errors"
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"github.com/shettyh/threadpool"
	"log"
	"ntc-gnats/nutil"
	"runtime"
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
	sopts = []nats.Option{nats.Name("NSubscriber_" + nutil.GetGUUID())}
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
	sopts = []nats.Option{nats.Name("NSubscriber_" + nutil.GetGUUID())}
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

func GetUrl() string {
	return surl
}

func GetAuth() string {
	return sauth
}

func GetOption() []nats.Option {
	return sopts
}

func GetConnect() (*nats.Conn, error) {
	// Connect to NATS
	nc, err := nats.Connect(surl, sopts...)
	if err != nil {
		log.Println(err)
	}
	return nc, err
}

/** Pool NSubscriber simple */
var DefaultQueueSize int64 = 100000

type PoolNSubscriber struct {
	queueSize int64
	poolNS *threadpool.ThreadPool
	listNSub []NSubscriber
}

func (pns *PoolNSubscriber) AddNSub(ns NSubscriber) {
	pns.listNSub = append(pns.listNSub, ns)
}

func (pns *PoolNSubscriber) RunPoolNSub() {
	if pns.queueSize <= 0 {
		pns.queueSize = DefaultQueueSize
	}
	if len(pns.listNSub) > 0 {
		numNSub := len(pns.listNSub)
		pns.poolNS = threadpool.NewThreadPool(numNSub, pns.queueSize)
		for i:=0; i<numNSub; i++ {
			ns := &pns.listNSub[i]
			pns.poolNS.Execute(ns)
		}
	}
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Running PoolNSubscriber size:", len(pns.listNSub), "NSubscriber")
}

func (pns *PoolNSubscriber) UnPoolNSub() {
	if len(pns.listNSub) > 0 {
		numNSub := len(pns.listNSub)
		pns.poolNS = threadpool.NewThreadPool(numNSub, pns.queueSize)
		for i:=0; i<numNSub; i++ {
			ns := &pns.listNSub[i]
			ns.NSSubt.Unsubscribe()
		}
	}
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Unsubscribe PoolNSubscriber size:", len(pns.listNSub), "NSubscriber")
}

type NSubscriber struct {
	ID      int32
	Subject string
	NSSubt  *nats.Subscription
}

func (ns *NSubscriber) Run() {
	fmt.Println("Running NSubscriber.ID:", ns.ID)
	// Connect to NATS
	nc, err := nats.Connect(surl, sopts...)
	if err != nil {
		log.Println(err)
	}
	ns.NSSubt, err = nc.Subscribe(ns.Subject, func(msg *nats.Msg) {
		log.Printf("NSubscriber[#%d] Received on PubSub [%s]: '%s'", ns.ID, ns.Subject, string(msg.Data))
	})
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	log.Printf("NSubscriber[#%d] is listening on Subject[%s]", ns.ID, ns.Subject)
	runtime.Goexit()
	fmt.Println("End NSubscriber.ID:", ns.ID)
}
