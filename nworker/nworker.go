package nworker

import (
	"errors"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"ntc-gnats/nutil"
	"strings"
	"time"
)

var wurl string
var wauth string
var wopts []nats.Option

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

func InitWorkerConf(name string) error {
	if name == "" {
		return errors.New("name config is not empty.")
	}
	wopts = []nats.Option{nats.Name("NWorker_" + nutil.GetGUUID())}
	c := nconf.GetConfig()
	wurl = c.GetString(name+".worker.url")
	wauth = c.GetString(name+".worker.auth")
	if len(wauth) > 0 {
		arrauth := strings.Split(wauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			wopts = append(wopts, nats.UserInfo(username, password))
		}
	}
	wopts = setupConnOptions(wopts)
	return nil
}

// url = nats://127.0.0.1:4222
// auth = username:password
func InitWorkerParams(url string, auth string) error {
	wopts = []nats.Option{nats.Name("NWorker_" + nutil.GetGUUID())}
	wurl = url
	wauth = auth
	if len(wauth) > 0 {
		arrauth := strings.Split(wauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			wopts = append(wopts, nats.UserInfo(username, password))
		}
	}
	wopts = setupConnOptions(wopts)
	return nil
}

func GetUrl() string {
	return wurl
}

func GetAuth() string {
	return wauth
}

func GetOption() []nats.Option {
	return wopts
}

func GetConnect() (*nats.Conn, error) {
	// Connect to NATS
	nc, err := nats.Connect(wurl, wopts...)
	if err != nil {
		log.Println(err)
	}
	return nc, err
}

///** Pool NWorker simple */
//var DefaultQueueSize int64 = 100000
//
//type PoolNWorker struct {
//	queueSize   int64
//	poolNW      *threadpool.ThreadPool
//	listNWorker []NWorker
//}
//
//func (pnw *PoolNWorker) AddNWorker(nw NWorker) {
//	pnw.listNWorker = append(pnw.listNWorker, nw)
//}
//
//func (pnw *PoolNWorker) RunPoolNWorker() {
//	if pnw.queueSize <= 0 {
//		pnw.queueSize = DefaultQueueSize
//	}
//	if len(pnw.listNWorker) > 0 {
//		numNWorker := len(pnw.listNWorker)
//		pnw.poolNW = threadpool.NewThreadPool(numNWorker, pnw.queueSize)
//		for i:=0; i<numNWorker; i++ {
//			nw := &pnw.listNWorker[i]
//			pnw.poolNW.Execute(nw)
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Running PoolNWorker size: %d NWorker", len(pnw.listNWorker))
//}
//
//// Disconnect safe
//func (pnw *PoolNWorker) DrainPoolNWorker() {
//	if len(pnw.listNWorker) > 0 {
//		numNWorker := len(pnw.listNWorker)
//		for i:=0; i<numNWorker; i++ {
//			nw := &pnw.listNWorker[i]
//			nw.Conn.Drain()
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Drain PoolNWorker size: %d NWorker", len(pnw.listNWorker))
//}
//
//// Disconnect unsafe
//func (pnw *PoolNWorker) UnPoolNWorker() {
//	if len(pnw.listNWorker) > 0 {
//		numNWorker := len(pnw.listNWorker)
//		for i:=0; i<numNWorker; i++ {
//			nw := &pnw.listNWorker[i]
//			nw.NWSubt.Unsubscribe()
//		}
//	}
//	time.Sleep(200 * time.Millisecond)
//	log.Printf("Unsubscribe PoolNWorker size: %d NWorker", len(pnw.listNWorker))
//}
//
//type NWorker struct {
//	ID        string
//	Subject   string
//	NameGroup string
//	Conn *nats.Conn
//	NWSubt    *nats.Subscription
//}
//
//func (nw *NWorker) Run() {
//	log.Printf("Running NWorker.ID: %s", nw.ID)
//	// Connect to NATS
//	var err error
//	nw.Conn, err = nats.Connect(wurl, wopts...)
//	if err != nil {
//		log.Println(err)
//	}
//	nw.NWSubt, err = nw.Conn.QueueSubscribe(nw.Subject, nw.NameGroup, func(msg *nats.Msg) {
//		log.Printf("NWorker[%s][#%s] Received on QueueWorker[%s]: '%s'", nw.NameGroup, nw.ID, nw.Subject, string(msg.Data))
//	})
//	nw.Conn.Flush()
//	if err := nw.Conn.LastError(); err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("NWorker[%s][#%s] is listening on Subject[%s]", nw.NameGroup, nw.ID, nw.Subject)
//	runtime.Goexit()
//	log.Printf("End NWorker.ID: %s", nw.ID)
//}
