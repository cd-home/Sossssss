package zine

import (
	"log"
	"sync"
)

type Handler struct {
	mu      sync.Mutex
	routers map[uint32]Router
}

func NewHandler() *Handler {
	return &Handler{routers: map[uint32]Router{}}
}

func (h *Handler) Handle(request *Request) {
	router, ok := h.routers[request.Msg.MsgID()]
	if !ok {
		log.Printf("Not Found Router Handle By MsgID [%d] ", request.Msg.ID)
		return
	}
	router.Handle(request)
}

func (h *Handler) AddRouter(msgID uint32, router Router) {
	h.mu.Lock()
	h.routers[msgID] = router
	h.mu.Unlock()
}
