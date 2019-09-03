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
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wanyvic/p2pssh/cmd/connect"
	"github.com/wanyvic/p2pssh/cmd/global"
	"github.com/wanyvic/p2pssh/version"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "p2pssh",
	Short:            "A Distributed Secure Shell",
	Long:             `A Distributed Secure Shell`,
	TraverseChildren: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: func(cmd *cobra.Command, args []string) error {
		return errors.Errorf("argument error")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := global.ConfigureDaemonLogs(&global.Opt); err != nil {
			logrus.Error(err)
		}
	},
	Version: fmt.Sprintf("%s, build %s", version.Version, version.GitCommit),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&global.Opt.CfgFile, "config", "", `config file (default is $HOME/.p2pssh.yaml)`)
	rootCmd.PersistentFlags().StringVarP(&global.Opt.LogLevel, "log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(NewDaemonCommand(*rootCmd))
	rootCmd.AddCommand(NewLoginCommand(*rootCmd))
	rootCmd.AddCommand(NewPingCommand(*rootCmd))
	rootCmd.AddCommand(connect.NewConnectCommand(*rootCmd))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if global.Opt.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(global.Opt.CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".p2pssh" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".p2pssh")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
