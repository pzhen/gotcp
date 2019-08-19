package Igotcp

type IServer interface {
	Start()
	Stop()
	Run()
	AddRouter(uint32, IRouter)
	GetManager() IManager


	SetOnConnStart(func(IConnector))
	SetOnConnStop(func(IConnector))

	CallOnConnStart(IConnector)
	CallOnConnStop(IConnector)
}
