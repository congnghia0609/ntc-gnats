/**
 *
 * @author nghiatc
 * @since Dec 17, 2019
 */
package nsub

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
//
//type NSubscriber struct {
//	ID      string
//	Subject string
//	Conn *nats.Conn
//	NSSubt  *nats.Subscription
//}
//
//func (ns *NSubscriber) Run() {
//	log.Printf("Running NSubscriber.ID: %s", ns.ID)
//	// Connect to NATS
//	var err error
//	ns.Conn, err = nats.Connect(surl, sopts...)
//	if err != nil {
//		log.Println(err)
//	}
//	ns.NSSubt, err = ns.Conn.Subscribe(ns.Subject, func(msg *nats.Msg) {
//		log.Printf("NSubscriber[#%s] Received on PubSub [%s]: '%s'", ns.ID, ns.Subject, string(msg.Data))
//	})
//	ns.Conn.Flush()
//	if err := ns.Conn.LastError(); err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("NSubscriber[#%s] is listening on Subject[%s]", ns.ID, ns.Subject)
//	runtime.Goexit()
//	log.Printf("End NSubscriber.ID: %s", ns.ID)
//}
