package ftpd

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	ErrSocketFormat = errors.New("socket format error")
	ErrServerClosed = errors.New("FTP Server Closed")
)

type FtpUser struct {
	Username   string
	Password   string
	HomeDir    string
	currentDir string
}

type FtpUserManager interface {
	Authenticate(string, string) (*FtpUser, error)
}

type FtpListener interface {
	OnStart(*FtpServer)
	OnConnect(*FtpSession)
	BeforeCommand(*FtpSession, *FtpRequest)
	AfterCommand(*FtpSession, *FtpRequest)
	OnDisconnect(*FtpSession)
	OnStop(*FtpServer)
}

func decoderSocket(arg string) (*net.TCPAddr, error) {
	// 127,0,0,1,50,199 12999
	args := strings.Split(arg, ",")
	if len(args) != 6 {
		return nil, ErrSocketFormat
	}

	ipv4 := args[0] + "." + args[1] + "." + args[2] + "." + args[3]

	high, err := strconv.Atoi(args[4])
	if err != nil {
		return nil, err
	}
	low, err := strconv.Atoi(args[5])
	if err != nil {
		return nil, err
	}

	port := (high << 8) | low

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ipv4, port))
	if err != nil {
		return nil, err
	}

	return addr, err
}

func encoderSocket(session FtpSession) {
}
