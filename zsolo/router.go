package zsolo

type Router interface {
	Handle(request *Request)
}
