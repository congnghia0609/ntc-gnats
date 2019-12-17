/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package nres

import (
	"errors"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"ntc-gnats/nutil"
	"strings"
	"time"
)

var rsurl string
var rsauth string
var rsopts []nats.Option

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

func InitResConf(name string) error {
	if name == "" {
		return errors.New("name config is not empty.")
	}
	rsopts = []nats.Option{nats.Name("NRes_" + nutil.GetGUUID())}
	c := nconf.GetConfig()
	rsurl = c.GetString(name+".res.url")
	rsauth = c.GetString(name+".res.auth")
	if len(rsauth) > 0 {
		arrauth := strings.Split(rsauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			rsopts = append(rsopts, nats.UserInfo(username, password))
		}
	}
	rsopts = setupConnOptions(rsopts)
	return nil
}

// url = nats://127.0.0.1:4222
// auth = username:password
func InitResParams(url string, auth string) error {
	rsopts = []nats.Option{nats.Name("NRes_" + nutil.GetGUUID())}
	rsurl = url
	rsauth = auth
	if len(rsauth) > 0 {
		arrauth := strings.Split(rsauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			rsopts = append(rsopts, nats.UserInfo(username, password))
		}
	}
	rsopts = setupConnOptions(rsopts)
	return nil
}

func GetUrl() string {
	return rsurl
}

func GetAuth() string {
	return rsauth
}

func GetOption() []nats.Option {
	return rsopts
}

func GetConnect() (*nats.Conn, error) {
	// Connect to NATS
	nc, err := nats.Connect(rsurl, rsopts...)
	if err != nil {
		log.Println(err)
	}
	return nc, err
}

///** Pool NRes simple */
//var DefaultQueueSize int64 = 100000
//
//type PoolNRes struct {
//	queueSize   int64
//	poolNW      *threadpool.ThreadPool
//	listNRes []NRes
//}
//
//func (pnw *PoolNRes) AddNRes(nw NRes) {
//	pnw.listNRes = append(pnw.listNRes, nw)
//}
//
//func (pnw *PoolNRes) RunPoolNRes() {
//	if pnw.queueSize <= 0 {
//		pnw.queueSize = DefaultQueueSize
//	}
//	if len(pnw.listNRes) > 0 {
//		numNRes := len(pnw.listNRes)
//		pnw.poolNW = threadpool.NewThreadPool(numNRes, pnw.queueSize)
//		for i:=0; i<numNRes; i++ {
//			nw := &pnw.listNRes[i]
//			pnw.poolNW.Execute(nw)
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Running PoolNRes size: %d NRes", len(pnw.listNRes))
//}
//
//// Disconnect safe
//func (pnw *PoolNRes) DrainPoolNRes() {
//	if len(pnw.listNRes) > 0 {
//		numNRes := len(pnw.listNRes)
//		for i:=0; i<numNRes; i++ {
//			nw := &pnw.listNRes[i]
//			nw.Conn.Drain()
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Drain PoolNRes size: %d NRes", len(pnw.listNRes))
//}
//
//// Disconnect unsafe
//func (pnw *PoolNRes) UnPoolNRes() {
//	if len(pnw.listNRes) > 0 {
//		numNRes := len(pnw.listNRes)
//		for i:=0; i<numNRes; i++ {
//			nw := &pnw.listNRes[i]
//			nw.NRSSubt.Unsubscribe()
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Unsubscribe PoolNRes size: %d NRes", len(pnw.listNRes))
//}
//
//type NRes struct {
//	ID        string
//	Subject   string
//	NameGroup string
//	Conn      *nats.Conn
//	NRSSubt   *nats.Subscription
//}
//
//func (nw *NRes) Run() {
//	log.Printf("Running NRes.ID: %s", nw.ID)
//	// Connect to NATS
//	var err error
//	nw.Conn, err = nats.Connect(rsurl, rsopts...)
//	if err != nil {
//		log.Println(err)
//	}
//	reply := "this is response ==> "
//	nw.NRSSubt, err = nw.Conn.QueueSubscribe(nw.Subject, nw.NameGroup, func(msg *nats.Msg) {
//		log.Printf("NRes[%s][#%s] Received on QueueNRes[%s]: '%s'", nw.NameGroup, nw.ID, nw.Subject, string(msg.Data))
//		datares := reply + string(msg.Data)
//		msg.Respond([]byte(datares))
//		log.Printf("NRes[%s][#%s] Reply on QueueNRes[%s]: '%s'", nw.NameGroup, nw.ID, nw.Subject, datares)
//	})
//	nw.Conn.Flush()
//	if err := nw.Conn.LastError(); err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("NRes[%s][#%s] is listening on Subject[%s]", nw.NameGroup, nw.ID, nw.Subject)
//	runtime.Goexit()
//	log.Printf("End NRes.ID: %s", nw.ID)
//}
