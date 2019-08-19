package Gotcp

import (
	"gotcp/Conf"
	"gotcp/Igotcp"
	"log"
	"strconv"
)

type MsgHandle struct {
	MsgRouterMap map[uint32]Igotcp.IRouter

	TaskQueue    []chan Igotcp.IRequest
	WorkPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		MsgRouterMap: make(map[uint32]Igotcp.IRouter),
		TaskQueue:    make([]chan Igotcp.IRequest, Conf.G_Conf.WorkPoolSize),
		WorkPoolSize: Conf.G_Conf.WorkPoolSize,
	}
}

func (mh *MsgHandle) DoMsgHandler(request Igotcp.IRequest) {
	router, ok := mh.MsgRouterMap[request.GetMsgId()]
	if !ok {
		log.Panicf("[Panic] can't find msgId : %d 's router", request.GetMsgId())
	}

	router.BeforeHandle(request)
	router.Handle(request)
	router.AfterHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router Igotcp.IRouter) {
	if _, ok := mh.MsgRouterMap[msgId]; ok {
		log.Panicln("[Panic] mh.MsgRouterMap[" + strconv.Itoa(int(msgId)) + "] : is registered")
	}
	mh.MsgRouterMap[msgId] = router
}

func (mh* MsgHandle) StartWorkerPool() {
	for i:= 0; i < int(mh.WorkPoolSize); i++ {
		mh.TaskQueue[i] = make(chan Igotcp.IRequest, Conf.G_Conf.MaxWorkPoolSize)
		go mh.startWorker(i, mh.TaskQueue[i])
	}
}

func (mh* MsgHandle) startWorker(workerId int, taskQueue chan Igotcp.IRequest) {
	log.Println("[Info] worker id = ", workerId, " is started ...")
	for {
		select {
		case request:= <- taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将请求交给TaskQueue， 由Worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request Igotcp.IRequest) {
	workerID := request.GetConnector().GetConnID() % mh.WorkPoolSize
	mh.TaskQueue[workerID] <- request
}
