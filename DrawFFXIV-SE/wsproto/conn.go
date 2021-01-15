package wsproto

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Conn : a helper class for serially read and write
type Conn struct {
	conn      *websocket.Conn
	readLock  *sync.Mutex
	writeLock *sync.Mutex
}

// WrapConn : wrap websocket connection
func WrapConn(conn *websocket.Conn) *Conn {
	return &Conn{
		conn:      conn,
		readLock:  &sync.Mutex{},
		writeLock: &sync.Mutex{},
	}
}

// SyncRead : read message from connection serially
func (conn *Conn) SyncRead() (msgType int, msg []byte, err error) {
	conn.readLock.Lock()
	if conn.conn != nil {
		msgType, msg, err = conn.conn.ReadMessage()
	} else {
		msgType, msg, err = websocket.TextMessage, make([]byte, 0), nil
	}
	conn.readLock.Unlock()
	return
}

// SyncWrite : write message to connection serially
func (conn *Conn) SyncWrite(msg string) (err error) {
	conn.writeLock.Lock()
	if conn.conn != nil {
		err = conn.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	} else {
		err = websocket.ErrCloseSent
	}
	conn.writeLock.Lock()
	return
}

// SyncWritef : write message with format to connection serially
func (conn *Conn) SyncWritef(format string, args ...interface{}) (err error) {
	return conn.SyncWrite(fmt.Sprintf(format, args...))
}

// Replace : replace new ws connection to Conn
func (conn *Conn) Replace(newConn *websocket.Conn) {
	conn.readLock.Lock()
	conn.writeLock.Lock()
	conn.conn = newConn
	conn.writeLock.Unlock()
	conn.readLock.Unlock()
}
