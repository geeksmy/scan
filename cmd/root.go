/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

var RootCmd = &cobra.Command{
	Use:          "scan",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		// tools.Banner()
		_ = cmd.Help()
	},
}

func init() {
	// 全局配置
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径")

	// 新增命令
	RootCmd.AddCommand(versionCmd())
	RootCmd.AddCommand(portCmd())
	RootCmd.AddCommand(blastingCmd())
	RootCmd.AddCommand(webFingerprintCmd())
	// RootCmd.AddCommand(cyberspaceCmd())
	RootCmd.SetHelpCommand(helpCmd())
	// rootCmd.AddCommand(poolCmd())
	// rootCmd.AddCommand(dirCmd())
	// rootCmd.AddCommand(dnsCmd())
}
