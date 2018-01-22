package communication

const UpdateMessage = "Update"

type Session interface {
	Publish(message string, data string)
	Subscribe(message string) chan string
	Close()
}

func NewSession() Session {
	return newMemSession()
}
