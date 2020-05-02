//+build windows

package ftpd

import "os"

func countLink(_ os.FileInfo) string {
	return "1"
}
