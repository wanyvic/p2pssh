// +build !linux

package client

import (
	"io"
	"os"

	"github.com/bitbored/go-ansicon"
)

func convertColor(r io.Reader) {
	w := ansicon.Convert(os.Stdout)
	io.Copy(w, r)
}
