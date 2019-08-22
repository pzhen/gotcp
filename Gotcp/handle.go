package Gotcp

import (
	"github.com/pkg/errors"
	"gotcp/Conf"
	"gotcp/Igotcp"
)

type MsgHandle struct {
	MsgRouterMap map[uint32]Igotcp.IRouter

	TaskQueue    []chan Igotcp.IRequest
	WorkPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		MsgRouterMap: make(map[uint32]Igotcp.IRouter),
		TaskQueue:    make([]chan Igotcp.IRequest, Conf.SrvConf.WorkPoolSize),
		WorkPoolSize: Conf.SrvConf.WorkPoolSize,
	}
}

func (mh *MsgHandle) DoMsgHandler(request Igotcp.IRequest) {
	router, ok := mh.MsgRouterMap[request.GetMsgId()]
	if !ok {
		debugPrintError("%+v", errors.WithStack(errors.Errorf("Handle router=%d not found", request.GetMsgId())))
		return
	}

	router.BeforeHandle(request)
	router.Handle(request)
	router.AfterHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router Igotcp.IRouter) {
	if _, ok := mh.MsgRouterMap[msgId]; ok {
		debugPrintError("%+v", errors.WithStack(errors.Errorf("Handle router=%d is registered", msgId)))
		return
	}
	mh.MsgRouterMap[msgId] = router
}

func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		mh.TaskQueue[i] = make(chan Igotcp.IRequest, Conf.SrvConf.MaxWorkPoolSize)
		go mh.startWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) startWorker(workerId int, taskQueue chan Igotcp.IRequest) {
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将请求交给TaskQueue， 由Worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request Igotcp.IRequest) {
	workerID := request.GetConnector().GetUUIDHashCode() % mh.WorkPoolSize
	mh.TaskQueue[workerID] <- request
}
