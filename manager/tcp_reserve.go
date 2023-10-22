package manager

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func NewTCPReserve(config ReserveConfig) (Reserve, error) {
	return &tcpReserve{
		config:       config,
		listSock:     nil,
		errorCounter: 0,
		nextTryTime:  time.Time{},
		lock:         &sync.RWMutex{},
		addrOrigin:   fmt.Sprintf("%s:%d", config.BaseHost, config.BasePort),
		addrReserve:  fmt.Sprintf("%s:%d", config.ReserveHost, config.ReservePort),
	}, nil
}

var _ Reserve = &tcpReserve{}

type tcpReserve struct {
	config       ReserveConfig
	listSock     net.Listener
	errorCounter int
	nextTryTime  time.Time
	lock         *sync.RWMutex

	//
	addrOrigin  string
	addrReserve string
}

func (tcpreserve *tcpReserve) Start(ctx context.Context) error {
	sock, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", tcpreserve.config.ListenPort))
	if err != nil {
		return err
	}
	tcpreserve.listSock = sock

	go tcpreserve.handleListener()

	return nil

}

func (tcpreserve *tcpReserve) Stop(ctx context.Context) error {
	return tcpreserve.listSock.Close()
}

func (tcpreserve *tcpReserve) Error() error {
	panic("not implemented") // TODO: Implement
}

func (tcpreserve *tcpReserve) handleListener() {
	for {
		conn, err := tcpreserve.listSock.Accept()
		if err != nil {
			return
		}
		go tcpreserve.handleConnection(conn)
	}
}

func (tcpreserve *tcpReserve) handleConnection(conn net.Conn) {
	remote, err := tcpreserve.getTarget()
	if err != nil {
		return
	}

	go tcpreserve.copyIO(conn, remote)
	tcpreserve.copyIO(remote, conn)
	conn.Close()
	remote.Close()
}

func (tcpreserve *tcpReserve) copyIO(from net.Conn, to net.Conn) {
	_, err := io.Copy(from, to)
	if err != nil {
		return
	}
}

func (tcpreserve *tcpReserve) getTarget() (remote net.Conn, err error) {
	timeout := time.Duration(tcpreserve.config.ConnectTimeout) * time.Second
	shouldTryOrigin := tcpreserve.shouldTryOrigin()
	logrus.Infof("new req: %s, shouldTryOrigin: %v details: %s", time.Now().UTC(), shouldTryOrigin, tcpreserve)
	if shouldTryOrigin {
		remote, err = tcpreserve.getOriging(timeout)
	} else {
		remote, err = tcpreserve.getRemote(timeout)
	}
	if err != nil {
		return
	}

	return remote, nil
}

func (tcpreserve *tcpReserve) getRemote(timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", tcpreserve.addrReserve, timeout)
}

func (tcpreserve *tcpReserve) getOriging(timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", tcpreserve.addrOrigin, timeout)
	if err != nil {
		tcpreserve.addError()
		// try to connect to remote
		return tcpreserve.getRemote(timeout)
	}
	if tcpreserve.lastErrored() {
		tcpreserve.resetError()
	}
	return conn, nil
}

func (tcpreserve *tcpReserve) getNextCheckTime() time.Time {
	n := time.Now().UTC()

	var secondsToWait int
	secondsToWaitFloat := math.Pow(2, float64(tcpreserve.errorCounter))

	if secondsToWaitFloat > float64(tcpreserve.config.MaxIdleSeconds) {
		secondsToWait = tcpreserve.config.MaxIdleSeconds
	} else {
		secondsToWait = int(secondsToWaitFloat)
	}
	return n.Add(time.Duration(secondsToWait) * time.Second)
}

func (tcpreserve *tcpReserve) lastErrored() bool {
	tcpreserve.RLock()
	defer tcpreserve.RUnLock()
	return tcpreserve.errorCounter > 0
}

func (tcpreserve *tcpReserve) shouldTryOrigin() bool {
	tcpreserve.RLock()
	defer tcpreserve.RUnLock()

	if tcpreserve.errorCounter == 0 {
		return true
	}
	nowNs := time.Now().UTC().UnixNano()
	tryTime := tcpreserve.nextTryTime.UnixNano()
	return nowNs > tryTime
}

func (tcpreserve *tcpReserve) resetError() {
	tcpreserve.Lock()
	defer tcpreserve.UnLock()
	tcpreserve.errorCounter = 0
	tcpreserve.nextTryTime = time.Now().UTC()
}

func (tcpreserve *tcpReserve) addError() {
	tcpreserve.Lock()
	defer tcpreserve.UnLock()
	tcpreserve.errorCounter += 1
	tcpreserve.nextTryTime = tcpreserve.getNextCheckTime()
	logrus.Infof("added error: after: %d : %s %s", tcpreserve.errorCounter, tcpreserve.nextTryTime, tcpreserve)
}

func (tcpreserve *tcpReserve) Lock() {
	tcpreserve.lock.Lock()
}

func (tcpreserve *tcpReserve) RLock() {
	tcpreserve.lock.RLock()
}

func (tcpreserve *tcpReserve) RUnLock() {
	tcpreserve.lock.RUnlock()
}

func (tcpreserve *tcpReserve) UnLock() {
	tcpreserve.lock.Unlock()
}

func (tcpreserve *tcpReserve) String() string {
	return fmt.Sprintf("<Reserve [%d] -> (%s {%d:%s}| %s)",
		tcpreserve.config.ListenPort,
		tcpreserve.addrOrigin, tcpreserve.errorCounter, tcpreserve.nextTryTime, tcpreserve.addrReserve)
}
