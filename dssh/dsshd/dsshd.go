package dsshd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/pidfile"
	"github.com/sirupsen/logrus"
	"github.com/wanyvic/dssh/dssh/dsshd/server"
)

var (
	PIDFile = "pssh.pid"
)

type DSSHdConfig struct {
	Addr string
	Port int
}

func Start(conf *DSSHdConfig) error {
	logrus.Debug("dsshd start")
	if err := PIDFileCheck(); err != nil {
		return err
	}
	svr, err := server.NewServer(conf.Addr, conf.Port)
	if err != nil {
		return err
	}
	err = svr.Start()
	if err != nil {
		return err
	}
	return nil
}

func PIDFileCheck() error {
	path := filepath.Join(os.TempDir(), PIDFile)
	file, err := pidfile.New(path)
	if err != nil {
		return errors.New("pssh daemon has been started, please stop first")
	}
	logrus.Debug(file)
	return nil
}
