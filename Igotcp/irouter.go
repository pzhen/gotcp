//路由 对应 应用层的 控制器 即 Controller
//不同路由去实现不同的业务
package Igotcp

type IRouter interface {
	BeforeHandle(request IRequest)
	Handle(request IRequest)
	AfterHandle(request IRequest)
}
