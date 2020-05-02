package ftpd

import (
	"errors"
	"log"
	"testing"
)

type fum struct {
}

func (fum) Authenticate(username, password string) (*FtpUser, error) {
	fu := &FtpUser{
		Username: "admin",
		Password: "123",
		HomeDir:  "D:\\Software",
	}

	if fu.Username == username && fu.Password == password {
		return fu, nil
	} else {
		return nil, errors.New("")
	}
}

func TestFtpServer(t *testing.T) {

	log.Print("Start...")

	var m fum
	opt := FtpServerOpt{
		Name:           defaultName,
		Host:           "",
		Port:           2121,
		WelcomeMessage: defaultWelcomeMessage,
		FtpUserManager: m,
	}

	server := NewFtpServer(&opt)

	if err := server.ListenAndServe(); err != nil {
		t.Error(err)
	}

}
