package Igotcp

type IHandle interface {
	DoMsgHandler(IRequest)
	AddRouter(uint32, IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(IRequest)
}
