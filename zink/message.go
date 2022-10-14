package zink

type Message struct {
	Len  uint32
	ID   uint32
	Data []byte
}

func NewMessage(ID uint32, data []byte) *Message {
	return &Message{
		Len:  uint32(len(data)),
		ID:   ID,
		Data: data,
	}
}

func (m *Message) MsgLen() uint32 {
	return m.Len
}

func (m *Message) MsgID() uint32 {
	return m.ID
}

func (m *Message) MsgData() []byte {
	return m.Data
}
