// +build !linux

package client

import (
	"io"
	"os"

	"github.com/shiena/ansicolor"
)

func convertColor(r io.Reader) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	io.Copy(w, r)
}
