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
package connect

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/wanyvic/p2pssh/api"
	"github.com/wanyvic/p2pssh/client"
	"github.com/wanyvic/p2pssh/cmd/global"
)

func NewConnectRmCommand(rootCmd cobra.Command) *cobra.Command {
	// rmCmd represents the rm command
	var rmCmd = &cobra.Command{
		Use:   "rm",
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
			if err := global.ConfigureDaemonLogs(&global.Opt); err != nil {
				fmt.Println(err)
				return
			}
			peer, err := peer.IDB58Decode(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			req := api.ConnectRmRequests{Peer: peer}
			res := api.ConnectRmResponses{}
			err = client.JsonRPConnect(global.Opt.DaemonAddress, "Server.ConnectRm", &req, &res)
			if err != nil {
				fmt.Println(err)
				return
			}
			if res.Err != "" {
				fmt.Println(res.Err)
			} else {
				fmt.Println(res.Result)
			}

		},
	}

	rmCmd.PersistentFlags().StringVarP(&global.Opt.DaemonAddress, "daemon-address", "D", api.DefaultDaemonAddress, `connection daemon address`)
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return rmCmd
}
