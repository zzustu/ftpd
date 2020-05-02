// +build !windows

package ftpd

import (
	"os"
	"strconv"
	"syscall"
)

func countLink(f os.FileInfo) string {

	str := "   "
	sys := f.Sys()
	if sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			sz := strconv.FormatUint(stat.Nlink, 10)
			if len(sz) >= len(str) {
				return sz
			}
			return str[0:len(str)-len(sz)] + sz
		}
	}
	return "  1"
}
