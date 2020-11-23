package messages

type MessageService interface {
	Create(content string) error
	GetAll() ([]Message, error)
}

func NewMessageService(s Store) MessageService {
	return &messageService{store: s}
}

type messageService struct {
	store Store
}

func (c *messageService) Create(content string) error {
	return c.store.Insert(Message{Content: content})
}

func (c *messageService) GetAll() ([]Message, error) {
	return c.store.GetAll()
}
