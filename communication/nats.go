package communication

import (
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/nats-io/go-nats"
)

type natsSession struct {
	nc            *nats.Conn
	subscriptions map[string]*nats.Subscription
}

func newNatsSession() *natsSession {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("unable to connect to nats %s", err.Error())
		return nil
	}

	var session *natsSession
	session.nc = nc
	session.subscriptions = make(map[string]*nats.Subscription)
	return session
}

func (session natsSession) Publish(message string, data string) {
	session.nc.Publish(message, []byte(data))
}

func (session *natsSession) Subscribe(message string) chan string {
	ch := make(chan string)
	sub, err := session.nc.Subscribe(message, func(msg *nats.Msg) {
		ch <- string(msg.Data)
	})
	if err != nil {
		log.Fatalf("unable to subscribe to %s", message)
	}
	session.subscriptions[message] = sub
	return ch
}

func (session natsSession) Close() {
	session.nc.Close()
	for _, sub := range session.subscriptions {
		sub.Unsubscribe()
	}
}
