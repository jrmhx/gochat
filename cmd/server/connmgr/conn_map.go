package connmgr

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type ConnMgr struct {
	mu   sync.Mutex
	sem  chan struct{} // max active conn
	ctab map[string]net.Conn
}

func NewConnMgr(mact uint) *ConnMgr {
	return &ConnMgr{
		sem:  make(chan struct{}, mact),
		ctab: make(map[string]net.Conn, mact),
	}
}

func (cm *ConnMgr) Add(addr net.Addr, conn net.Conn) error {
	key := addr.String()
	select {
	case cm.sem <- struct{}{}:
		// reserver a slot
	case <-time.After(500 * time.Millisecond):
		conn.Write([]byte("connection timeout"))
		conn.Close()
		return fmt.Errorf("connection pool full, %s timeout", key)
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if _, inmap := cm.ctab[key]; inmap == false {
		cm.ctab[key] = conn
		return nil
	} else {
		<-cm.sem // release the slot
		return fmt.Errorf("connection %s already exist", key)
	}
}

func (cm *ConnMgr) Delete(addr net.Addr) {
	key := addr.String()
	cm.mu.Lock()
	conn, prs := cm.ctab[key]
	if prs {
		delete(cm.ctab, key)
	}
	cm.mu.Unlock()
	if !prs {
		return
	}
	if err := conn.Close(); err != nil {
		fmt.Println("close conn error:", err)
	}
	<-cm.sem
}

func (cm *ConnMgr) Clear() {
	cm.mu.Lock()
	conns := make([]net.Conn, 0, len(cm.ctab))
	for k, conn := range cm.ctab {
		conns = append(conns, conn)
		delete(cm.ctab, k)
	}
	cm.mu.Unlock()

	for _, conn := range conns {
		conn.Close()
		<-cm.sem
	}
}
