/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */

package npub

import (
	"errors"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
)

var purl string
var pauth string
var popts []nats.Option

func InitPubConf(name string) error {
	if name == "" {
		return errors.New("name config is not empty.")
	}
	popts = []nats.Option{nats.Name("NPublisher-")}
	c := nconf.GetConfig()
	purl = c.GetString(name+".pub.url")
	//fmt.Println("purl:", purl)
	pauth = c.GetString(name+".pub.auth")
	//fmt.Println("pauth:", pauth)
	if len(pauth) > 0 {
		arrauth := strings.Split(pauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			popts = append(popts, nats.UserInfo(username, password))
		}
	}
	return nil
}

// url = nats://127.0.0.1:4222
// auth = username:password
func InitPubParams(url string, auth string) error {
	popts = []nats.Option{nats.Name("NPublisher")}
	purl = url
	pauth = auth
	if len(pauth) > 0 {
		arrauth := strings.Split(pauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			popts = append(popts, nats.UserInfo(username, password))
		}
	}
	return nil
}

func GetUrl() string {
	return purl
}

func GetAuth() string {
	return pauth
}

func GetOption() []nats.Option {
	return popts
}

func GetConnect() (*nats.Conn, error) {
	// Connect to NATS
	nc, err := nats.Connect(purl, popts...)
	if err != nil {
		log.Println(err)
	}
	return nc, err
}

func Publish(subject string, msg string) error {
	// Connect to NATS
	nc, err := nats.Connect(purl, popts...)
	defer nc.Close()
	if err != nil {
		log.Println(err)
	}
	nc.Publish(subject, []byte(msg))
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	return err
}
