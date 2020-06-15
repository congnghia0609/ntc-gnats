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
	NSSubt  *nats.Subscription
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	//totalWait := 10 * time.Minute
	//reconnectDelay := time.Second
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
	if name == "" {
		return nil
	}
	c := nconf.GetConfig()
	id := name + "_NSubscriber_" + nutil.GetGUUID()
	sopts := []nats.Option{nats.Name(id)}
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
	return &NSubscriber{Name: id, Subject: subject, Url: surl, Auth: sauth, Opts: sopts, Conn: nc, NSSubt: nil}
}

func (ns *NSubscriber) Start(processChan chan *nats.Msg) error {
	log.Printf("Running NSubscriber.Name: %s", ns.Name)
	// Subscribe to NATS
	var err error
	ns.NSSubt, err = ns.Conn.Subscribe(ns.Subject, func(msg *nats.Msg) {
		//log.Printf("NSubscriber[#%s] Received on PubSub [%s]: '%s'", ns.Name, ns.Subject, string(msg.Data))
		processChan <- msg
	})
	ns.Conn.Flush()
	if err = ns.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	log.Printf("NSubscriber[#%s] is listening on Subject[%s]", ns.Name, ns.Subject)
	//runtime.Goexit()
	//log.Printf("End NSubscriber.Name: %s", ns.Name)
	return err
}

// Disconnect safe
func (ns *NSubscriber) Stop() {
	if ns != nil && ns.Conn != nil {
		ns.Conn.Drain()
	}
}

//var surl string
//var sauth string
//var sopts []nats.Option
//
//func setupConnOptions(opts []nats.Option) []nats.Option {
//	//totalWait := 10 * time.Minute
//	reconnectDelay := time.Second
//	opts = append(opts, nats.ReconnectWait(reconnectDelay))
//	opts = append(opts, nats.MaxReconnects(math.MaxInt32))
//	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
//		//log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
//		log.Printf("Disconnected due to:%s, will attempt reconnects for %v times", err, math.MaxInt32)
//	}))
//	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
//		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
//	}))
//	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
//		log.Fatalf("Exiting: %v", nc.LastError())
//	}))
//	return opts
//}
//
//func InitSubConf(name string) error {
//	if name == "" {
//		return errors.New("name config is not empty.")
//	}
//	sopts = []nats.Option{nats.Name("NSubscriber_" + nutil.GetGUUID())}
//	c := nconf.GetConfig()
//	surl = c.GetString(name+".sub.url")
//	log.Printf("surl=%s", surl)
//	sauth = c.GetString(name+".sub.auth")
//	if len(sauth) > 0 {
//		arrauth := strings.Split(sauth, ":")
//		if len(arrauth) == 2 {
//			username := arrauth[0]
//			password := arrauth[1]
//			sopts = append(sopts, nats.UserInfo(username, password))
//		}
//	}
//	sopts = setupConnOptions(sopts)
//	return nil
//}
//
//// url = nats://127.0.0.1:4222
//// auth = username:password
//func InitSubParams(url string, auth string) error {
//	sopts = []nats.Option{nats.Name("NSubscriber_" + nutil.GetGUUID())}
//	surl = url
//	sauth = auth
//	if len(sauth) > 0 {
//		arrauth := strings.Split(sauth, ":")
//		if len(arrauth) == 2 {
//			username := arrauth[0]
//			password := arrauth[1]
//			sopts = append(sopts, nats.UserInfo(username, password))
//		}
//	}
//	sopts = setupConnOptions(sopts)
//	return nil
//}
//
//func GetUrl() string {
//	return surl
//}
//
//func GetAuth() string {
//	return sauth
//}
//
//func GetOption() []nats.Option {
//	return sopts
//}

//func GetConnect() (*nats.Conn, error) {
//	// Connect to NATS
//	nc, err := nats.Connect(surl, sopts...)
//	if err != nil {
//		log.Println(err)
//	}
//	return nc, err
//}
//
///** Pool NSubscriber simple */
//var DefaultQueueSize int64 = 100000
//
//type PoolNSubscriber struct {
//	queueSize int64
//	poolNS *threadpool.ThreadPool
//	listNSub []NSubscriber
//}
//
//func (pns *PoolNSubscriber) AddNSub(ns NSubscriber) {
//	pns.listNSub = append(pns.listNSub, ns)
//}
//
//func (pns *PoolNSubscriber) RunPoolNSub() {
//	if pns.queueSize <= 0 {
//		pns.queueSize = DefaultQueueSize
//	}
//	if len(pns.listNSub) > 0 {
//		numNSub := len(pns.listNSub)
//		pns.poolNS = threadpool.NewThreadPool(numNSub, pns.queueSize)
//		for i:=0; i<numNSub; i++ {
//			ns := &pns.listNSub[i]
//			pns.poolNS.Execute(ns)
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Running PoolNSubscriber size: %d NSubscriber", len(pns.listNSub))
//}
//
//// Disconnect safe
//func (pns *PoolNSubscriber) DrainPoolNSub() {
//	if len(pns.listNSub) > 0 {
//		numNSub := len(pns.listNSub)
//		for i:=0; i<numNSub; i++ {
//			ns := &pns.listNSub[i]
//			ns.Conn.Drain()
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Drain PoolNSubscriber size: %d NSubscriber", len(pns.listNSub))
//}
//
//// Disconnect unsafe
//func (pns *PoolNSubscriber) UnPoolNSub() {
//	if len(pns.listNSub) > 0 {
//		numNSub := len(pns.listNSub)
//		for i:=0; i<numNSub; i++ {
//			ns := &pns.listNSub[i]
//			ns.NSSubt.Unsubscribe()
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Unsubscribe PoolNSubscriber size: %d NSubscriber", len(pns.listNSub))
//}
