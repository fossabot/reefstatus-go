package communication

import "github.com/cjburchell/pubsub"

const UpdateMessage = "Update"

type Session interface {
	Publish(subj string, data []byte)
	Subscribe(subj string) chan []byte
	Close()
}

type memSession struct {
}

func NewSession() Session {
	var session memSession
	return &session
}

func (memSession) Publish(subj string, data []byte) {
	pubsub.Publish(subj, data)
}

func (memSession) Subscribe(subj string) chan []byte {
	channel := make(chan []byte)
	pubsub.SubscribeChan(subj, channel)
	return channel
}

func (memSession) Close() {
}
