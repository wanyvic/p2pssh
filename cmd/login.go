/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wanyvic/p2pssh/api"
	"github.com/wanyvic/p2pssh/client"
	"github.com/wanyvic/p2pssh/p2pssh/login"
)

const (
	DefaultSSHPrivateKey = "$HOME/.ssh/id_rsa"
	DefaultDaemonAddress = "127.0.0.1:9001"
)

var (
	SSHPrivateKey string
	DaemonAddress string
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Debug("login called")
		if err := configureDaemonLogs(&Opt); err != nil {
			logrus.Error(err)
		}
		if len(args) <= 0 {
			logrus.Error("No connection")
			return
		}
		logrus.Debug(DaemonAddress, SSHPrivateKey)
		config := &api.ClientConfig{}
		var err error
		if config, err = configureClientConfig(args[0]); err != nil {
			logrus.Error(err)
			return
		}
		if tcpAddr, err := parseConnection(DaemonAddress); err != nil {
			logrus.Error(err)
			return
		} else {
			cli := client.New(context.Background(), tcpAddr, *config)
			cli.ConnHandler = client.SSHandle
			if err := cli.Connect(); err != nil {
				logrus.Error(err)
				return
			}
		}
	},
}

func init() {

	rootCmd.AddCommand(loginCmd)

	loginCmd.PersistentFlags().StringVarP(&SSHPrivateKey, "privkey", "P", "", `ssh private key file such as `+DefaultSSHPrivateKey)
	loginCmd.PersistentFlags().StringVarP(&DaemonAddress, "daemon-address", "D", DefaultDaemonAddress, `connection daemon address`)
}
func configureClientConfig(connInfo string) (*api.ClientConfig, error) {
	config, err := login.ParseClientConfig(connInfo, SSHPrivateKey)
	if err != nil {
		return nil, err
	}
	logrus.Debug("UserName: ", config.UserName, " NodeID: ", config.NodeID)

	return &config, nil
}
func parseConnection(valueAddr string) (*net.TCPAddr, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", valueAddr)
	if err != nil {
		return nil, err
	}
	return tcpAddr, nil
}
