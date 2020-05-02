package ftpd

import (
	"bufio"
	"context"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	mutex sync.RWMutex

	defaultName           = "Go FTP Server"
	defaultWelcomeMessage = "Welcome to FTP Server"
)

type FtpServerOpt struct {
	Name           string
	Host           string
	Port           int
	WelcomeMessage string
	FtpUserManager FtpUserManager
}

type FtpServer struct {
	ftpListener map[string]FtpListener
	opt         *FtpServerOpt
	listen      net.Listener
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewFtpServer(opt *FtpServerOpt) *FtpServer {
	return &FtpServer{
		ftpListener: nil,
		opt:         opt,
		listen:      nil,
		ctx:         nil,
		cancel:      nil,
	}
}

func (s *FtpServer) ListenAndServe() error {

	var err error

	s.listen, err = net.Listen("tcp", net.JoinHostPort(s.opt.Host, strconv.Itoa(s.opt.Port)))
	if err != nil {
		return err
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.onStart()

	for {
		conn, err := s.listen.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return ErrServerClosed
			default:
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					continue
				}
				return err
			}
		}

		session := s.newFtpSession(conn)
		go session.handler()
	}

}

func (s *FtpServer) Shutdown() error {
	if s.cancel != nil {
		s.cancel()
	}
	if s.listen != nil {
		return s.listen.Close()
	}
	return nil
}

func (s *FtpServer) newFtpSession(conn net.Conn) *FtpSession {

	session := new(FtpSession)
	now := time.Now()

	session.CtrlConn = conn
	session.CtrlReader = bufio.NewReader(conn)
	session.CtrlWriter = bufio.NewWriter(conn)
	session.FtpServer = *s
	session.RemoteAddr = conn.RemoteAddr()
	session.LocalAddr = conn.LocalAddr()
	session.IsLoginedIn = false
	session.ConnectAt = now
	session.LastAccessAt = now
	session.CurrentDir = "/"

	return session
}

func (s *FtpServer) onStart() {
	mutex.RLock()
	defer mutex.RUnlock()
	if s.ftpListener != nil {
		for _, v := range s.ftpListener {
			v.OnStart(*s)
		}
	}
}

func (s *FtpServer) onConnect(session FtpSession) {
	mutex.RLock()
	defer mutex.RUnlock()
	if s.ftpListener != nil {
		for _, v := range s.ftpListener {
			v.OnConnect(session)
		}
	}
}

func (s *FtpServer) beforeCommand(session FtpSession, request FtpRequest) {
	mutex.RLock()
	defer mutex.RUnlock()
	if s.ftpListener != nil {
		for _, v := range s.ftpListener {
			v.BeforeCommand(session, request)
		}
	}
}

func (s *FtpServer) afterCommand(session FtpSession, request FtpRequest, reply int) {
	mutex.RLock()
	defer mutex.RUnlock()
	if s.ftpListener != nil {
		for _, v := range s.ftpListener {
			v.AfterCommand(session, request, reply)
		}
	}
}

func (s *FtpServer) onDisconnect(session FtpSession) {
	mutex.RLock()
	defer mutex.RUnlock()
	if s.ftpListener != nil {
		for _, v := range s.ftpListener {
			v.OnDisconnect(session)
		}
	}
}

func (s *FtpServer) onStop() {
	mutex.RLock()
	defer mutex.RUnlock()
	if s.ftpListener != nil {
		for _, v := range s.ftpListener {
			v.OnStop(*s)
		}
	}
}

func (s *FtpServer) AddListener(name string, lsn FtpListener) {
	mutex.Lock()
	defer mutex.Unlock()
	if s.ftpListener == nil {
		s.ftpListener = make(map[string]FtpListener)
	}
	s.ftpListener[name] = lsn
}

func (s *FtpServer) RemoveListener(name string) {
	mutex.Lock()
	defer mutex.Unlock()
	if s.ftpListener != nil {
		delete(s.ftpListener, name)
	}
}
