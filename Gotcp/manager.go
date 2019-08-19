package Gotcp

import (
	"errors"
	"gotcp/Igotcp"
	"log"
	"sync"
)

type Manager struct {
	connMap  map[uint32]Igotcp.IConnector
	connLock sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		connMap:make(map[uint32] Igotcp.IConnector),
	}
}

func (mgr *Manager) Add(conn Igotcp.IConnector) {
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()
	mgr.connMap[conn.GetConnID()] = conn
	log.Printf("[Info] connection id : %d is added manager map connMap length : %d", conn.GetConnID(), mgr.Len())
}

func (mgr *Manager) Remove(conn Igotcp.IConnector) {
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()
	delete(mgr.connMap, conn.GetConnID())
	log.Printf("[Info] connection id : %d is removed manager map connMap length : %d", conn.GetConnID(), mgr.Len())
}

func (mgr *Manager) Get(connID uint32) (Igotcp.IConnector, error) {
	mgr.connLock.RLock()
	defer mgr.connLock.RUnlock()

	if conn, ok := mgr.connMap[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("[Warning] connection not found")
	}
}

func (mgr *Manager) Len() int {
	return len(mgr.connMap)
}

func (mgr *Manager) ClearConn() {
	mgr.connLock.Lock()
	defer mgr.connLock.Unlock()
	for connID, conn := range mgr.connMap {
		//停止
		conn.Stop()
		//删除
		delete(mgr.connMap, connID)
	}
	log.Println("[Info] Clear All connections success connMap length : ", mgr.Len())
}
