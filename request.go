package ftpd

import (
	"strings"
	"time"
)

type FtpRequest struct {
	Line       string
	Command    string
	Argument   string
	ReceivedAt time.Time
}

func parseLine(line string) *FtpRequest {
	params := strings.SplitN(strings.TrimSpace(line), " ", 2)
	request := &FtpRequest{
		Line:       line,
		ReceivedAt: time.Now(),
	}
	sz := len(params)
	if sz > 0 {
		request.Command = strings.ToUpper(params[0])
		if sz == 2 {
			request.Argument = params[1]
		}
	}
	return request
}
