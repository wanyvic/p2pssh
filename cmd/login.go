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
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wanyvic/p2pssh/client"
	"github.com/wanyvic/p2pssh/p2pssh/login"
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
		fmt.Println("login called")
		if err := configureDaemonLogs(&Opt); err != nil {
			logrus.Error(err)
		}
		if len(args) > 0 {
			userAuth, err := login.ParseAuth(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			userAuth.Password = "lrh19950815"
			logrus.Debug(userAuth)
			c := client.New(context.Background(), client.DefaultConnect(), userAuth)
			if err = c.Connect(); err != nil {
				logrus.Error(err)
			}
			select {}
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	rootCmd.PersistentFlags().StringVarP(&Opt.CfgFile, "Pubkey", "P", "", `login with pubkey`)
}
