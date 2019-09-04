// +build !windows

package client

import (
	"io"
	"os"
)

func convertColor(r io.Reader) {
	io.Copy(os.Stdout, r)
}
