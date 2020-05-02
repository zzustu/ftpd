package ftpd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var (
	delim      = " "
	newline    = "\r\n"
	owner      = "user"
	group      = "group"
	timePatten = "Jan _2 15:04"
)

type FtpFile interface {
	Perm() string
}

type FileFormater interface {
	format([]os.FileInfo) []byte
}

type nlstFileFormater struct{}

func (l nlstFileFormater) format(fs []os.FileInfo) []byte {
	var buf bytes.Buffer
	for _, f := range fs {
		_, _ = fmt.Fprintf(&buf, "%s%s", f.Name(), newline)
	}
	return buf.Bytes()
}

func formatTime(t time.Time, patten string) string {
	return t.Format(patten)
}

type listFileFormater struct{}

func (l listFileFormater) format(fs []os.FileInfo) []byte {
	var buf bytes.Buffer
	for _, f := range fs {
		buf.WriteString(f.Mode().String())
		buf.WriteString(delim)
		buf.WriteString(delim)
		buf.WriteString(delim)
		buf.WriteString(countLink(f))
		buf.WriteString(delim)
		buf.WriteString(owner)
		buf.WriteString(delim)
		buf.WriteString(group)
		buf.WriteString(delim)
		buf.WriteString(formatSize(f.Size()))
		buf.WriteString(delim)
		buf.WriteString(formatTime(f.ModTime(), timePatten))
		buf.WriteString(delim)
		buf.WriteString(f.Name())
		buf.WriteString(newline)
	}
	return buf.Bytes()
}

func formatSize(sz int64) string {
	str := "            "
	size := strconv.FormatInt(sz, 10)
	if len(size) >= len(str) {
		return size
	}

	return str[0:len(str)-len(size)] + size
}

func getFileList(path string, f FileFormater) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		fs := []os.FileInfo{info}
		return f.format(fs), nil
	} else {
		fs, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		return f.format(fs), nil
	}
}
