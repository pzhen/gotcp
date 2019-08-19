package Igotcp

type IManager interface {
	//添加链接
	Add(conn IConnector)
	//删除连接
	Remove(conn IConnector)
	//根据connID获取链接
	Get(connID uint32) (IConnector, error)
	//得到当前连接总数
	Len() int
	//清除并终止所有d连接
	ClearConn()
}
