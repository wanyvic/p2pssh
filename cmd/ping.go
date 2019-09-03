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

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wanyvic/p2pssh/api"
	"github.com/wanyvic/p2pssh/client"
	"github.com/wanyvic/p2pssh/cmd/global"
)

func NewPingCommand(rootCmd cobra.Command) *cobra.Command {
	// pingCmd represents the ping command
	var pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.Errorf("argument error")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Debug("ping called")
			if err := global.ConfigureDaemonLogs(&global.Opt); err != nil {
				logrus.Error(err)
			}
			NodeID, err := peer.IDB58Decode(args[0])
			config := &api.ClientConfig{}
			if err != nil {
				logrus.Error(err)
				return
			} else {
				config.NodeID = NodeID
				if tcpAddr, err := parseConnection(global.Opt.DaemonAddress); err != nil {
					logrus.Error(err)
					return
				} else {
					cli := client.New(context.Background(), tcpAddr, *config)
					cli.ConnHandler = client.PingHandle
					if err := cli.Connect(); err != nil {
						logrus.Error(err)
						return
					}
				}
			}
		},
	}
	pingCmd.PersistentFlags().StringVarP(&global.Opt.DaemonAddress, "daemon-address", "D", api.DefaultDaemonAddress, `connection daemon address`)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(pingCmd)
	return pingCmd
}
