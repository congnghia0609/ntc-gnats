package nworker

import (
	"github.com/nats-io/nats.go"
	"github.com/shettyh/threadpool"
	"log"
	"runtime"
	"time"
)

/** Pool NWorker simple */
var DefaultQueueSize int64 = 100000

type PoolNWorker struct {
	queueSize   int64
	poolNW      *threadpool.ThreadPool
	listNWorker []NWorker
}

func (pnw *PoolNWorker) AddNWorker(nw NWorker) {
	pnw.listNWorker = append(pnw.listNWorker, nw)
}

func (pnw *PoolNWorker) RunPoolNWorker() {
	if pnw.queueSize <= 0 {
		pnw.queueSize = DefaultQueueSize
	}
	if len(pnw.listNWorker) > 0 {
		numNWorker := len(pnw.listNWorker)
		pnw.poolNW = threadpool.NewThreadPool(numNWorker, pnw.queueSize)
		for i:=0; i<numNWorker; i++ {
			nw := &pnw.listNWorker[i]
			pnw.poolNW.Execute(nw)
		}
	}
	time.Sleep(200 * time.Millisecond)
	log.Printf("Running PoolNWorker size: %d NWorker", len(pnw.listNWorker))
}

// Disconnect safe
func (pnw *PoolNWorker) DrainPoolNWorker() {
	if len(pnw.listNWorker) > 0 {
		numNWorker := len(pnw.listNWorker)
		for i:=0; i<numNWorker; i++ {
			nw := &pnw.listNWorker[i]
			nw.Conn.Drain()
		}
	}
	time.Sleep(200 * time.Millisecond)
	log.Printf("Drain PoolNWorker size: %d NWorker", len(pnw.listNWorker))
}

// Disconnect unsafe
func (pnw *PoolNWorker) UnPoolNWorker() {
	if len(pnw.listNWorker) > 0 {
		numNWorker := len(pnw.listNWorker)
		for i:=0; i<numNWorker; i++ {
			nw := &pnw.listNWorker[i]
			nw.NWSubt.Unsubscribe()
		}
	}
	time.Sleep(200 * time.Millisecond)
	log.Printf("Unsubscribe PoolNWorker size: %d NWorker", len(pnw.listNWorker))
}

type NWorker struct {
	ID        string
	Subject   string
	NameGroup string
	Conn *nats.Conn
	NWSubt    *nats.Subscription
}

func (nw *NWorker) Run() {
	log.Printf("Running NWorker.ID: %s", nw.ID)
	// Connect to NATS
	var err error
	nw.Conn, err = nats.Connect(wurl, wopts...)
	if err != nil {
		log.Println(err)
	}
	nw.NWSubt, err = nw.Conn.QueueSubscribe(nw.Subject, nw.NameGroup, func(msg *nats.Msg) {
		log.Printf("NWorker[%s][#%s] Received on QueueWorker[%s]: '%s'", nw.NameGroup, nw.ID, nw.Subject, string(msg.Data))
	})
	nw.Conn.Flush()
	if err := nw.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	log.Printf("NWorker[%s][#%s] is listening on Subject[%s]", nw.NameGroup, nw.ID, nw.Subject)
	runtime.Goexit()
	log.Printf("End NWorker.ID: %s", nw.ID)
}
