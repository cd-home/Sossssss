package message

// Message
// 消息格式设计
// Target 私聊[好友、陌生人]、群聊
type Message struct {
	Action string
	From   uint32
	To     uint32
	Target string
	Type   string
	Body   string
}
