package messages

type Store interface {
	Insert(Message) error
	GetAll() ([]Message, error)
}

func NewInMemoryStore() Store {
	messages := make([]Message, 0)

	return &inMemoryStore{messages: messages}
}

type inMemoryStore struct {
	messages []Message
}

func (i *inMemoryStore) Insert(msg Message) error {
	i.messages = append(i.messages, msg)

	return nil
}

func (i *inMemoryStore) GetAll() ([]Message, error) {
	return i.messages, nil
}
