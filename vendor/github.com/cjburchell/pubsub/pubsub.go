package pubsub

// Subscription interface
type Subscription interface {
	Close()
}

// MsgHandler is a callback function that processes messages delivered to
// asynchronous subscribers.
type MsgHandler func(msg []byte)

type subscription struct {
	subject string
	channel chan []byte
	handler MsgHandler
}

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// Close the subscription
func (sub *subscription) Close() {
	session := sessions[sub.subject]
	session.subscriptions = removeSubscription(session.subscriptions, sliceIndex(len(session.subscriptions), func(i int) bool { return session.subscriptions[i] == sub }))
	if len(session.subscriptions) == 0 {
		delete(sessions, sub.subject)
	}
}

type session struct {
	subscriptions []*subscription
}

var sessions = make(map[string]*session, 0)

// SubscribeChan subscribes to the subject using a channel
func SubscribeChan(subject string, channel chan []byte) Subscription {
	var sub = &subscription{
		subject: subject,
		channel: channel,
	}
	if s, ok := sessions[subject]; ok {
		s.subscriptions = append(s.subscriptions, sub)
	} else {
		var s session
		s.subscriptions = append(s.subscriptions, sub)
		sessions[subject] = &s
	}

	return sub
}

// Subscribe to a subject
func Subscribe(subject string, handler MsgHandler) Subscription {
	var sub = &subscription{
		subject: subject,
		handler: handler,
	}
	if s, ok := sessions[subject]; ok {
		s.subscriptions = append(s.subscriptions, sub)
	} else {
		var s session
		s.subscriptions = append(s.subscriptions, sub)
		sessions[subject] = &s
	}

	return sub
}

// Publish a message to all subscribers of a subject
func Publish(subject string, data []byte) {
	if session, ok := sessions[subject]; ok {
		for _, sub := range session.subscriptions {
			if sub.channel != nil {
				sub.channel <- data
			} else {
				go sub.handler(data)
			}
		}
	}
}

func removeSubscription(s []*subscription, i int) []*subscription {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
