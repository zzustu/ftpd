package ftpd

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type commander interface {
	Execute(*FtpSession, FtpRequest)
}

var (
	// 无需用户权限的命令
	nonAuthenticatedCommands = [6]string{"USER", "PASS", "AUTH", "QUIT", "PROT", "PBSZ"}

	commands = map[string]commander{
		"ABOR": abor{},
		"ACCT": acct{},
		"APPE": appe{},
		"AUTH": auth{},
		"CDUP": cdup{},
		"CWD":  cwd{},
		"DELE": dele{},
		"EPRT": eprt{},
		"EPSV": epsv{},
		"FEAT": feat{},
		"HELP": help{},
		"LANG": lang{},
		"LIST": list{},
		"MD5":  md5{},
		"MFMT": mfmt{},
		"MMD5": md5{},
		"MDTM": mdtm{},
		"MLST": mlst{},
		"MKD":  mkd{},
		"MLSD": mlsd{},
		"MODE": mode{},
		"NLST": nlst{},
		"NOOP": noop{},
		"OPTS": opts{},
		"PASS": pass{},
		//"PASV":          pasv{},
		"PBSZ": pbsz{},
		"PORT": port{},
		"PROT": prot{},
		"PWD":  pwd{},
		"QUIT": quit{},
		"REIN": rein{},
		"REST": rest{},
		"RETR": retr{},
		"RMD":  rmd{},
		"RNFR": rnfr{},
		"RNTO": rnto{},
		//"SITE":          site{},
		"SIZE":          size{},
		"SITE_DESCUSER": siteDescuser{},
		"SITE_HELP":     siteHelp{},
		"SITE_STAT":     siteStat{},
		"SITE_WHO":      siteWho{},
		"SITE_ZONE":     siteZone{},
		"STAT":          stat{},
		"STOR":          stor{},
		"STOU":          stou{},
		"STRU":          stru{},
		"SYST":          syst{},
		"TYPE":          typeCommand{},
		"USER":          user{},
		"XPWD":          pwd{},
	}

	optsMap = map[string]commander{
		"OPTS_MLST": optsMlst{},
		"OPTS_UTF8": optsUTF8{},
	}
)

// 判断某个命令是否无需认证权限
func isWithoutAuthenticationCommand(command string) bool {
	for _, cmd := range nonAuthenticatedCommands {
		if cmd == command {
			return true
		}
	}
	return false
}

type abor struct{}

func (cmd abor) Execute(session *FtpSession, request FtpRequest) {
	session.CloseDataConn()
	session.write(reply226ClosingDataConnection, "ABOR command successful.")
}

type acct struct{}

func (cmd acct) Execute(session *FtpSession, request FtpRequest) {
	session.write(reply202CommandNotImplemented, "Command ACCT not implemented, superfluous at this site.")
}

type appe struct{}

func (cmd appe) Execute(session *FtpSession, request FtpRequest) {

}

type auth struct{}

func (cmd auth) Execute(session *FtpSession, request FtpRequest) {

}

type cdup struct{}

func (cmd cdup) Execute(session *FtpSession, request FtpRequest) {

}

type cwd struct{}

func (cmd cwd) Execute(session *FtpSession, request FtpRequest) {

	path, info, err := session.buildPath(request.Argument)

	if err != nil || !info.IsDir() {
		session.write(reply550RequestedActionNotTaken, "No such directory.")
		return
	}

	session.CurrentDir = path
	session.write(reply250RequestedFileActionOkay, fmt.Sprintf("\"%s\" is current directory.", path))
}

type dele struct{}

func (cmd dele) Execute(session *FtpSession, request FtpRequest) {

}

type eprt struct{}

func (cmd eprt) Execute(session *FtpSession, request FtpRequest) {

}

type epsv struct{}

func (cmd epsv) Execute(session *FtpSession, request FtpRequest) {

}

type feat struct{}

func (cmd feat) Execute(session *FtpSession, request FtpRequest) {

}

type help struct{}

func (cmd help) Execute(session *FtpSession, request FtpRequest) {

}

type lang struct{}

func (cmd lang) Execute(session *FtpSession, request FtpRequest) {

}

type list struct{}

func (cmd list) Execute(session *FtpSession, request FtpRequest) {

	argument := request.Argument
	root := session.FtpUser.HomeDir
	path := filepath.Join(root, filepath.Join(session.CurrentDir, argument))

	files, err := getFileList(path, new(listFileFormater))
	if err != nil {
		session.write(reply550RequestedActionNotTaken, "No such directory.")
		return
	}

	session.write(reply150FileStatusOkay, "Opening ASCII mode data connection for file list")

	session.writeData(files)
}

type md5 struct{}

func (cmd md5) Execute(session *FtpSession, request FtpRequest) {

}

type mfmt struct{}

func (cmd mfmt) Execute(session *FtpSession, request FtpRequest) {

}

type mmd5 struct{}

func (cmd mmd5) Execute(session *FtpSession, request FtpRequest) {

}

type mdtm struct{}

func (cmd mdtm) Execute(session *FtpSession, request FtpRequest) {

}

type mlst struct{}

func (cmd mlst) Execute(session *FtpSession, request FtpRequest) {

}

type mkd struct{}

func (cmd mkd) Execute(session *FtpSession, request FtpRequest) {

}

type mlsd struct{}

func (cmd mlsd) Execute(session *FtpSession, request FtpRequest) {

}

type mode struct{}

func (cmd mode) Execute(session *FtpSession, request FtpRequest) {

}

type nlst struct{}

func (cmd nlst) Execute(session *FtpSession, request FtpRequest) {

	argument := request.Argument
	root := session.FtpUser.HomeDir
	path := filepath.Join(root, filepath.Join(session.CurrentDir, argument))

	files, err := getFileList(path, new(nlstFileFormater))
	if err != nil {
		session.write(reply503BadSequenceOfCommands, "POR121T or PASV must be issued first.")
		return
	}

	session.write(reply150FileStatusOkay, "Opening ASCII mode data connection for file list")

	session.writeData(files)
}

type noop struct{}

func (cmd noop) Execute(session *FtpSession, request FtpRequest) {
	session.write(reply200CommandOkay, "Command NOOP okay.")
}

type opts struct{}

func (cmd opts) Execute(session *FtpSession, request FtpRequest) {
	argument := request.Argument
	if argument == "" {
		session.write(reply501SyntaxErrorInParametersOrArguments, "Syntax error in parameters or arguments.")
		return
	}

	args := strings.Split(argument, " ")
	if len(args) == 0 {
		session.write(reply501SyntaxErrorInParametersOrArguments, "Syntax error in parameters or arguments.")
		return
	}

	code := "OPTS_" + strings.ToUpper(args[0])
	c := optsMap[code]
	if c == nil {
		session.write(reply502CommandNotImplemented, "OPTS not implemented.")
		return
	}

	c.Execute(session, request)
}

type pass struct{}

func (cmd pass) Execute(session *FtpSession, request FtpRequest) {
	password := request.Argument
	username := session.getAttribute(attributeUserArgument)
	if username == "" && session.FtpUser == nil {
		session.write(reply503BadSequenceOfCommands, "Login with USER first.")
		return
	}

	if session.IsLoginedIn {
		session.write(reply202CommandNotImplemented, "Already logged-in.")
		return
	}

	ftpUser, err := session.FtpServer.opt.FtpUserManager.Authenticate(username, password)
	if err != nil {
		session.write(reply530NotLoggedIn, "Authentication failed.")
		return
	}

	session.IsLoginedIn = true
	session.FtpUser = ftpUser
	session.write(reply230UserLoggedIn, "User logged in, proceed.")
}

type pasv struct{}

func (cmd pasv) Execute(session *FtpSession, request FtpRequest) {

}

type pbsz struct{}

func (cmd pbsz) Execute(session *FtpSession, request FtpRequest) {

}

type port struct{}

func (cmd port) Execute(session *FtpSession, request FtpRequest) {
	argument := request.Argument
	if argument == "" {
		session.write(reply501SyntaxErrorInParametersOrArguments, "Syntax error in parameters or arguments.")
		return
	}
	// TODO 判断是否开启主动模式
	addr, err := decoderSocket(argument)
	if err != nil {
		session.write(reply501SyntaxErrorInParametersOrArguments, "Syntax error in parameters or arguments.")
		return
	}

	// TODO 判断是否开启被动模式IP检查再决定是否检查IP地址
	if n, ok := session.RemoteAddr.(*net.TCPAddr); ok {
		if !addr.IP.Equal(n.IP) {
			session.write(reply501SyntaxErrorInParametersOrArguments, "Syntax error in parameters or arguments.")
			return
		}
	}

	// TODO 绑定数据通道
	conn, err := newPortModeConn(addr)
	if err != nil {
		log.Print(err)
		session.write(reply501SyntaxErrorInParametersOrArguments, "Syntax error in parameters or arguments.")
		return
	}

	if session.DataConn != nil {
		_ = session.DataConn.Close()
		session.DataConn = nil
	}

	session.DataConn = conn

	log.Printf("[%s] Enable PORT Mode, Destination: %s", session.RemoteAddr, addr)

	session.write(reply200CommandOkay, "Command PORT okay.")
}

type prot struct{}

func (cmd prot) Execute(session *FtpSession, request FtpRequest) {

}

type pwd struct{}

func (cmd pwd) Execute(session *FtpSession, _ FtpRequest) {
	session.write(reply257PathNameCreated, fmt.Sprintf("\"%s\" is current directory.", session.CurrentDir))
}

type quit struct{}

func (cmd quit) Execute(session *FtpSession, request FtpRequest) {
	session.write(reply221ClosingControlConnection, "Goodbye.")
	session.Close()
}

type rein struct{}

func (cmd rein) Execute(session *FtpSession, request FtpRequest) {

}

type rest struct{}

func (cmd rest) Execute(session *FtpSession, request FtpRequest) {

}

type retr struct{}

func (cmd retr) Execute(session *FtpSession, request FtpRequest) {
	abspath, _ := session.getFilePath(request.Argument)

	f, err := os.Open(abspath)
	if err != nil {
		session.write(reply550RequestedActionNotTaken, "No such file or directory.")
		return
	}
	defer func() {
		_ = f.Close()
	}()
	if fi, err := os.Stat(abspath); err != nil || fi.IsDir() {
		session.write(reply550RequestedActionNotTaken, "Not a plain file.")
		return
	}

	session.write(reply150FileStatusOkay, "Data transfer starting.")
	session.writeFile(f)
}

type rmd struct{}

func (cmd rmd) Execute(session *FtpSession, request FtpRequest) {

}

type rnfr struct{}

func (cmd rnfr) Execute(session *FtpSession, request FtpRequest) {

}

type rnto struct{}

func (cmd rnto) Execute(session *FtpSession, request FtpRequest) {

}

type site struct{}

func (cmd site) Execute(session *FtpSession, request FtpRequest) {
	argument := request.Argument
	if argument == "" {
		session.write(reply200CommandOkay, "Command SITE okay. Use SITE HELP to get more information.")
		return
	}

	code := "SITE_" + strings.ToUpper(argument)

	c := commands[code]
	if c == nil {
		session.write(reply502CommandNotImplemented, "Command SITE not implemented for "+argument)
		return
	}

	c.Execute(session, request)
}

type size struct{}

func (cmd size) Execute(session *FtpSession, request FtpRequest) {
	_, info, err := session.buildPath(request.Argument)
	if err != nil {
		session.write(reply550RequestedActionNotTaken, "No such file or directory.")
		return
	}

	if info.IsDir() {
		session.write(reply550RequestedActionNotTaken, "Not a plain file.")
		return
	}

	session.write(reply213FileStatus, strconv.FormatInt(info.Size(), 10))
}

type siteDescuser struct{}

func (cmd siteDescuser) Execute(session *FtpSession, _ FtpRequest) {
	u := session.FtpUser
	message := fmt.Sprintf("\nusername : %s\npassword : ******\nhome dir : %s", u.Username, u.HomeDir)
	session.write(reply200CommandOkay, message)
}

type siteHelp struct{}

func (cmd siteHelp) Execute(session *FtpSession, _ FtpRequest) {
	message := "\nDESCUSER : display user information." +
		"\nHELP     : display this message." +
		"\nSTAT     : show statistics." +
		"\nWHO      : display all connected users." +
		"\nZONE     : display timezone."
	session.write(reply200CommandOkay, message)
}

type siteStat struct{}

func (cmd siteStat) Execute(session *FtpSession, _ FtpRequest) {
	message := "\nwill todo"
	session.write(reply200CommandOkay, message)
}

type siteWho struct{}

func (cmd siteWho) Execute(session *FtpSession, _ FtpRequest) {
	message := "\nwill todo"
	session.write(reply200CommandOkay, message)
}

type siteZone struct{}

func (cmd siteZone) Execute(session *FtpSession, _ FtpRequest) {
	s := time.Now()
	session.write(reply200CommandOkay, s.String())
}

//
type stat struct{}

func (cmd stat) Execute(session *FtpSession, request FtpRequest) {

}

type stor struct{}

func (cmd stor) Execute(session *FtpSession, request FtpRequest) {

}

type stou struct{}

func (cmd stou) Execute(session *FtpSession, request FtpRequest) {

}

type stru struct{}

func (cmd stru) Execute(session *FtpSession, request FtpRequest) {

}

type syst struct{}

func (cmd syst) Execute(session *FtpSession, request FtpRequest) {
	session.write(reply215NameSystemType, fmt.Sprintf("UNIX Type: %s", session.FtpServer.opt.Name))
}

type typeCommand struct{}

func (cmd typeCommand) Execute(session *FtpSession, request FtpRequest) {
	if request.Argument == "" {
		session.write(reply501SyntaxErrorInParametersOrArguments, "Syntax error in parameters or arguments.")
		return
	}

	t := strings.ToUpper(request.Argument[:1])
	switch t {
	case "A":
		session.setAttribute(attributeDataType, dataTypeAscii)
	case "I":
		session.setAttribute(attributeDataType, dataTypeBinary)
	default:
		session.write(reply504CommandNotImplementedForThatParameter, fmt.Sprintf("Command TYPE not implemented for the parameter %s.", request.Argument))
		return
	}

	session.write(reply200CommandOkay, "TYPE Command Okay.")
}

type user struct{}

func (cmd user) Execute(session *FtpSession, request FtpRequest) {
	username := request.Argument
	if session.IsLoginedIn {
		if session.FtpUser.Username == username {
			session.write(reply230UserLoggedIn, "Already logged-in.")
		} else {
			session.write(reply530NotLoggedIn, "Invalid user name.")
		}
		return
	}

	session.setAttribute(attributeUserArgument, username)
	session.write(reply331UserNameOkayNeedPassword, "User name okay, need password.")
}

type optsMlst struct{}

func (cmd optsMlst) Execute(session *FtpSession, request FtpRequest) {

}

type optsUTF8 struct{}

func (cmd optsUTF8) Execute(session *FtpSession, request FtpRequest) {
	session.write(reply200CommandOkay, "Command OPTS okay.")
}
