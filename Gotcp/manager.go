package Gotcp

import (
	"github.com/pkg/errors"
	"gotcp/Igotcp"
	"sync"
)

type Manager struct {
	connMap  map[string]Igotcp.IConnector
	connLock sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		connMap:make(map[string] Igotcp.IConnector),
	}
}

func (mgr *Manager) Add(conn Igotcp.IConnector) {
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()
	mgr.connMap[conn.GetUUID()] = conn
}

func (mgr *Manager) Remove(conn Igotcp.IConnector) {
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()
	delete(mgr.connMap, conn.GetUUID())
}

func (mgr *Manager) Get(uuid string) (Igotcp.IConnector, error) {
	mgr.connLock.RLock()
	defer mgr.connLock.RUnlock()

	if conn, ok := mgr.connMap[uuid]; ok {
		return conn, nil
	} else {
		return nil, errors.New("Manager UUID=" + uuid + " not found")
	}
}

func (mgr *Manager) Len() int {
	return len(mgr.connMap)
}

func (mgr *Manager) ClearConn() {
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()
	for uuid, conn := range mgr.connMap {
		conn.Stop()
		delete(mgr.connMap, uuid)
	}
	debugPrint("Manager cleared all UUIDS", mgr.Len())
}
