/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */

package npub

import (
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"math"
	"ntc-gnats/nutil"
	"strings"
	"sync"
)

var mNP sync.Mutex
var mapInstanceNP = map[string]*NPublisher{}

type NPublisher struct {
	Name string
	Url  string
	Auth string
	Opts []nats.Option
	Conn *nats.Conn
}

func NewNPublisher(name string) *NPublisher {
	if len(name) == 0 {
		fmt.Errorf("name config is not empty.")
		return nil
	}
	c := nconf.GetConfig()
	id := name + "_NPublisher_" + nutil.GetGUUID()
	popts := []nats.Option{nats.Name(id), nats.MaxReconnects(math.MaxInt32)}
	purl := c.GetString(name + ".pub.url")
	fmt.Printf("purl=%s\n", purl)
	pauth := c.GetString(name + ".pub.auth")
	//log.Printf("pauth: %s\n", pauth)
	if len(pauth) > 0 {
		arrauth := strings.Split(pauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			popts = append(popts, nats.UserInfo(username, password))
		}
	}
	// Connect to NATS
	nc, err := nats.Connect(purl, popts...)
	if err != nil {
		fmt.Printf("Connect to NATS Fail: %v\n", err)
	}
	return &NPublisher{Name: id, Url: purl, Auth: pauth, Opts: popts, Conn: nc}
}

func GetInstance(name string) *NPublisher {
	instance := mapInstanceNP[name]
	if instance == nil || instance.Conn.IsClosed() {
		mNP.Lock()
		defer mNP.Unlock()
		instance = mapInstanceNP[name]
		if instance == nil || instance.Conn.IsClosed() {
			instance = NewNPublisher(name)
			mapInstanceNP[name] = instance
		}
	}
	return instance
}

func (np *NPublisher) Publish(subject string, msg string) error {
	if len(subject) > 0 && len(msg) > 0 {
		np.Conn.Publish(subject, []byte(msg))
		np.Conn.Flush()
		err := np.Conn.LastError()
		if err != nil {
			log.Fatalf("NPublisher publish message Error: %v\n", err)
		}
		return err
	}
	return nil
}

func (np *NPublisher) PublishByte(subject string, msg []byte) error {
	if len(subject) > 0 && len(msg) > 0 {
		np.Conn.Publish(subject, msg)
		np.Conn.Flush()
		err := np.Conn.LastError()
		if err != nil {
			log.Fatalf("NPublisher publish message Error: %v\n", err)
		}
		return err
	}
	return nil
}

func (np *NPublisher) Close() {
	if np != nil {
		np.Conn.Close()
	}
}

func Publish(name string, subject string, msg string) error {
	np := GetInstance(name)
	return np.Publish(subject, msg)
}

func PublishByte(name string, subject string, msg []byte) error {
	np := GetInstance(name)
	return np.PublishByte(subject, msg)
}
