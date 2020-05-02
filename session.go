package ftpd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	attributeUserArgument = "user-argument"
	attributeDataType     = "data-type"
	dataTypeAscii         = "ASCII"
	dataTypeBinary        = "Binary"
)

type FtpSession struct {
	CtrlConn   net.Conn
	CtrlReader *bufio.Reader
	CtrlWriter *bufio.Writer
	DataConn   DataConn

	FtpServer  FtpServer
	FtpUser    *FtpUser
	RemoteAddr net.Addr
	LocalAddr  net.Addr
	ConnectAt  time.Time

	IsLoginedIn  bool
	LastAccessAt time.Time

	CurrentDir string

	Attribute map[string]string
}

func (session *FtpSession) handler() {

	session.write(reply220ServiceReady, defaultWelcomeMessage)

	log.Printf("[%s] =====[+++]=====", session.RemoteAddr)

	for {
		line, err := session.CtrlReader.ReadString('\n')
		if err != nil {
			break
		}
		session.interpreter(line)
	}

	// 关闭FTP连接
	session.Close()

	// 断开连接时通知监听器
	session.FtpServer.onDisconnect(*session)

	log.Printf("[%s] =====[!!!]=====", session.RemoteAddr)
}

func (session *FtpSession) interpreter(line string) {

	session.LastAccessAt = time.Now()

	request := parseLine(line)

	log.Printf("[%s] >>> %s", session.RemoteAddr, request.Line)

	// 判断该命令是否是需要权限认证的命令
	if !session.IsLoginedIn && !isWithoutAuthenticationCommand(request.Command) {
		session.write(reply530NotLoggedIn, "Access denied")
		return
	}

	// 判断命令是否存在
	c := commands[request.Command]
	if c == nil {
		session.write(reply502CommandNotImplemented, "Command not implemented")
		return
	}

	// 在执行命令前触发beforeCommand监听器
	session.FtpServer.beforeCommand(*session, request)

	// 开始执行命令
	c.Execute(session, request)

	// 在执行命令前触发afterCommand监听器
	session.FtpServer.afterCommand(*session, request, 0)
}

func (session *FtpSession) getAttribute(key string) string {
	if session.Attribute == nil {
		return ""
	}
	return session.Attribute[key]
}

func (session *FtpSession) setAttribute(key, value string) {
	if session.Attribute == nil {
		session.Attribute = make(map[string]string)
	}
	session.Attribute[key] = value
}

func (session *FtpSession) removeAttribute(key string) {
	if session.Attribute != nil {
		delete(session.Attribute, key)
	}
}

func (session *FtpSession) Close() {
	if err := session.CtrlConn.Close(); err != nil {
		if _, ok := err.(*net.OpError); !ok {
			log.Print(err)
		}
	}
	session.CloseDataConn()
}

func (session *FtpSession) CloseDataConn() {
	if session.DataConn != nil {
		// 关闭数据通道
		if err := session.DataConn.Close(); err != nil {
			log.Print(err)
		}
		// 将FTP Session的数据通道置空
		session.DataConn = nil
	}
}

func (session *FtpSession) buildPath(path string) (string, os.FileInfo, error) {

	abspath, sandpath := session.getFilePath(path)

	info, err := os.Stat(abspath)
	return sandpath, info, err
}

func (session *FtpSession) getFilePath(path string) (string, string) {
	// 逻辑路径(即: FTP用户所看到的绝对路径)
	sandpath := session.CurrentDir
	if len(path) > 0 {
		if path[:1] == "/" {
			sandpath = filepath.Clean(path)
		} else {
			sandpath = filepath.Join(sandpath, path)
		}
	}

	// Windows下的路径分割符是 '\' 要转义成 '/'
	abspath := filepath.Join(session.FtpUser.HomeDir, sandpath)
	sandpath = strings.Replace(sandpath, string(filepath.Separator), "/", -1)

	return abspath, sandpath
}

// 向控制通道写入返回信息
func (session FtpSession) write(reply int, message string) {
	msg := fmt.Sprintf("%d %s\n", reply, message)
	_, err := session.CtrlWriter.WriteString(msg)
	if err == nil {
		_ = session.CtrlWriter.Flush()
	}
}

// 往数据通道写入数据
func (session *FtpSession) writeData(data []byte) {

	// 检查数据通道是否开启
	if session.DataConn == nil {
		session.write(reply503BadSequenceOfCommands, "PORT or PASV must be issued first.")
		return
	}

	// 向数据通道写入数据
	if _, err := session.DataConn.Write(data); err != nil {
		log.Print(err)
	}

	message := "Closing data connection, sent " + strconv.Itoa(len(data)) + " bytes"
	session.write(reply226ClosingDataConnection, message)

	// 完毕后关闭数据通道
	session.CloseDataConn()
}

func (session *FtpSession) writeFile(data io.ReadCloser) {

	// 检查数据通道是否开启
	if session.DataConn == nil {
		session.write(reply503BadSequenceOfCommands, "PORT or PASV must be issued first.")
		return
	}

	sz, _ := io.Copy(session.DataConn, data)

	message := "Closing data connection, sent " + strconv.FormatInt(sz, 10) + " bytes"
	session.write(reply226ClosingDataConnection, message)

	// 完毕后关闭数据通道
	session.CloseDataConn()
}
