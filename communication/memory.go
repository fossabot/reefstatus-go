package communication

import "github.com/cjburchell/reefstatus-go/common"

type memSession struct {
	subscriptions map[string]chan string
}

var sessions = make([]*memSession, 0)

func newMemSession() *memSession {
	var session *memSession
	session.subscriptions = make(map[string]chan string)
	sessions = append(sessions, session)
	return session
}

func (memSession) Publish(message string, data string) {
	for _, session := range sessions {
		if session == nil {
			continue
		}

		if ch, ok := session.subscriptions[message]; ok {
			ch <- message
		}
	}
}

func (session *memSession) Subscribe(message string) chan string {
	ch := make(chan string)
	session.subscriptions[message] = ch
	return ch
}

func remove(s []*memSession, i int) []*memSession {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (session *memSession) Close() {
	sessions = remove(sessions, common.SliceIndex(len(sessions), func(i int) bool { return sessions[i] == session }))
}
