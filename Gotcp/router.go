//路由的实现
//具体业务都应该继承该类去重写
package Gotcp

import "gotcp/Igotcp"

type BaseRouter struct {
}

func (b *BaseRouter) BeforeHandle(request Igotcp.IRequest) {}

func (b *BaseRouter) Handle(request Igotcp.IRequest) {}

func (b *BaseRouter) AfterHandle(request Igotcp.IRequest) {}
