/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package nres

import (
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"math"
	"ntc-gnats/nutil"
	"strings"
)

type NResponse struct {
	Name      string
	Subject   string
	Group string
	Url     string
	Auth    string
	Opts    []nats.Option
	Conn      *nats.Conn
	Subt   *nats.Subscription
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

func NewNResponse(name string) *NResponse {
	if len(name) == 0 {
		return nil
	}
	c := nconf.GetConfig()
	id := name + "_NResponse_" + nutil.GetGUUID()
	rsopts := []nats.Option{nats.Name(id), nats.MaxReconnects(math.MaxInt32)}
	subject := c.GetString(name + ".res.subject")
	group := c.GetString(name + ".res.group")
	rsurl := c.GetString(name+".res.url")
	log.Printf("rsurl=%s", rsurl)
	rsauth := c.GetString(name+".res.auth")
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
		//log.Printf("NRes[%s][#%s] Received on QueueNRes[%s]: '%s'", nw.NameGroup, nw.ID, nw.Subject, string(msg.Data))
		//datares := reply + string(msg.Data)
		//msg.Respond([]byte(datares))
		//log.Printf("NRes[%s][#%s] Reply on QueueNRes[%s]: '%s'", nw.NameGroup, nw.ID, nw.Subject, datares)
	})
	nrs.Conn.Flush()
	if err = nrs.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("NRes[%s][#%s] is listening on Subject[%s]", nrs.Group, nrs.Name, nrs.Subject)
	//runtime.Goexit()
	//log.Printf("End NRes.ID: %s", nrs.Name)
	return err
}

// Disconnect safe
func (nrs *NResponse) Stop() {
	if nrs != nil && nrs.Conn != nil {
		nrs.Conn.Drain()
	}
}

//var rsurl string
//var rsauth string
//var rsopts []nats.Option
//
//func InitResConf(name string) error {
//	if name == "" {
//		return errors.New("name config is not empty.")
//	}
//	rsopts = []nats.Option{nats.Name("NRes_" + nutil.GetGUUID()), nats.MaxReconnects(math.MaxInt32)}
//	c := nconf.GetConfig()
//	rsurl = c.GetString(name+".res.url")
//	log.Printf("rsurl=%s", rsurl)
//	rsauth = c.GetString(name+".res.auth")
//	if len(rsauth) > 0 {
//		arrauth := strings.Split(rsauth, ":")
//		if len(arrauth) == 2 {
//			username := arrauth[0]
//			password := arrauth[1]
//			rsopts = append(rsopts, nats.UserInfo(username, password))
//		}
//	}
//	rsopts = setupConnOptions(rsopts)
//	return nil
//}
//
//// url = nats://127.0.0.1:4222
//// auth = username:password
//func InitResParams(url string, auth string) error {
//	rsopts = []nats.Option{nats.Name("NRes_" + nutil.GetGUUID())}
//	rsurl = url
//	rsauth = auth
//	if len(rsauth) > 0 {
//		arrauth := strings.Split(rsauth, ":")
//		if len(arrauth) == 2 {
//			username := arrauth[0]
//			password := arrauth[1]
//			rsopts = append(rsopts, nats.UserInfo(username, password))
//		}
//	}
//	rsopts = setupConnOptions(rsopts)
//	return nil
//}
//
//func GetUrl() string {
//	return rsurl
//}
//
//func GetAuth() string {
//	return rsauth
//}
//
//func GetOption() []nats.Option {
//	return rsopts
//}
//
//func GetConnect() (*nats.Conn, error) {
//	// Connect to NATS
//	nc, err := nats.Connect(rsurl, rsopts...)
//	if err != nil {
//		log.Println(err)
//	}
//	return nc, err
//}
//
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
//
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
