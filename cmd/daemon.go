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

	"github.com/wanyvic/p2pssh/cmd/global"
	p2p "github.com/wanyvic/p2pssh/libp2p"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wanyvic/p2pssh/service"
)

type daemonOptions struct {
	PrivateKey string
}

var daemonOpt daemonOptions

func NewDaemonCommand(rootCmd cobra.Command) *cobra.Command {

	// daemonCmd represents the daemon command
	var daemonCmd = &cobra.Command{
		Use:   "daemon",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := global.ConfigureDaemonLogs(&global.Opt); err != nil {
				logrus.Error(err)
			}
			p2p.PrivateKey = daemonOpt.PrivateKey
			cli := p2p.GetLibp2p()
			cli.NewSSHService()
			cli.NewPingService()
			svi := service.New(context.Background(), service.DefaultConnect())
			svi.ConnHandler = service.Handle
			if err := svi.Start(); err != nil {
				logrus.Error(err)
			}
			select {}
		},
	}
	// rootCmd.PersistentFlags().Int16VarP(&Opt.CfgFile, "config", "", `config file (default is $HOME/.p2pssh.yaml)`)
	// Here you will define your flags and configuration settings.

	daemonCmd.PersistentFlags().StringVarP(&daemonOpt.PrivateKey, "privkey", "", "", `daemon with private key for libp2p`)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(daemonCmd)
	return daemonCmd
}
