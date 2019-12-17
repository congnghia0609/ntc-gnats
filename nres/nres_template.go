package nres

import (
	"github.com/nats-io/nats.go"
	"github.com/shettyh/threadpool"
	"log"
	"runtime"
	"time"
)

/** Pool NRes simple */
var DefaultQueueSize int64 = 100000

type PoolNRes struct {
	queueSize   int64
	poolNW      *threadpool.ThreadPool
	listNRes []NRes
}

func (pnw *PoolNRes) AddNRes(nw NRes) {
	pnw.listNRes = append(pnw.listNRes, nw)
}

func (pnw *PoolNRes) RunPoolNRes() {
	if pnw.queueSize <= 0 {
		pnw.queueSize = DefaultQueueSize
	}
	if len(pnw.listNRes) > 0 {
		numNRes := len(pnw.listNRes)
		pnw.poolNW = threadpool.NewThreadPool(numNRes, pnw.queueSize)
		for i:=0; i<numNRes; i++ {
			nw := &pnw.listNRes[i]
			pnw.poolNW.Execute(nw)
		}
	}
	time.Sleep(200 * time.Millisecond)
	log.Printf("Running PoolNRes size: %d NRes", len(pnw.listNRes))
}

// Disconnect safe
func (pnw *PoolNRes) DrainPoolNRes() {
	if len(pnw.listNRes) > 0 {
		numNRes := len(pnw.listNRes)
		for i:=0; i<numNRes; i++ {
			nw := &pnw.listNRes[i]
			nw.Conn.Drain()
		}
	}
	time.Sleep(200 * time.Millisecond)
	log.Printf("Drain PoolNRes size: %d NRes", len(pnw.listNRes))
}

// Disconnect unsafe
func (pnw *PoolNRes) UnPoolNRes() {
	if len(pnw.listNRes) > 0 {
		numNRes := len(pnw.listNRes)
		for i:=0; i<numNRes; i++ {
			nw := &pnw.listNRes[i]
			nw.NRSSubt.Unsubscribe()
		}
	}
	time.Sleep(200 * time.Millisecond)
	log.Printf("Unsubscribe PoolNRes size: %d NRes", len(pnw.listNRes))
}

type NRes struct {
	ID        string
	Subject   string
	NameGroup string
	Conn      *nats.Conn
	NRSSubt   *nats.Subscription
}

func (nw *NRes) Run() {
	log.Printf("Running NRes.ID: %s", nw.ID)
	// Connect to NATS
	var err error
	nw.Conn, err = nats.Connect(rsurl, rsopts...)
	if err != nil {
		log.Println(err)
	}
	reply := "this is response ==> "
	nw.NRSSubt, err = nw.Conn.QueueSubscribe(nw.Subject, nw.NameGroup, func(msg *nats.Msg) {
		log.Printf("NRes[%s][#%s] Received on QueueNRes[%s]: '%s'", nw.NameGroup, nw.ID, nw.Subject, string(msg.Data))
		datares := reply + string(msg.Data)
		msg.Respond([]byte(datares))
		log.Printf("NRes[%s][#%s] Reply on QueueNRes[%s]: '%s'", nw.NameGroup, nw.ID, nw.Subject, datares)
	})
	nw.Conn.Flush()
	if err := nw.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	log.Printf("NRes[%s][#%s] is listening on Subject[%s]", nw.NameGroup, nw.ID, nw.Subject)
	runtime.Goexit()
	log.Printf("End NRes.ID: %s", nw.ID)
}
