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

// Execute is the main cobra method
// func Execute() {
// 	var cancel context.CancelFunc
// 	mainContext, cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	signalChan := make(chan os.Signal, 1)
// 	signal.Notify(signalChan, os.Interrupt)
// 	defer func() {
// 		signal.Stop(signalChan)
// 		cancel()
// 	}()
// 	go func() {
// 		select {
// 		case <-signalChan:
// 			// caught CTRL+C
// 			fmt.Println("\n[!] 检测到键盘中断，正在终止。")
// 			cancel()
// 		case <-mainContext.Done():
// 		}
// 	}()
//
// 	if err := rootCmd.Execute(); err != nil {
// 		// Leaving this in results in the same error appearing twice
// 		// Once before and once after the help output. Not sure if
// 		// this is going to be needed to output other errors that
// 		// aren't automatically outputted.
// 		// fmt.Println(err)
// 		os.Exit(1)
// 	}
// }

func init() {
	// 全局配置
	// rootCmd.PersistentFlags().StringP("url", "u", "", "需要扫描的ip或url")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径")

	// 新增命令
	RootCmd.AddCommand(versionCmd())
	RootCmd.AddCommand(portCmd())
	RootCmd.AddCommand(blastingCmd())
	RootCmd.AddCommand(webFingerprintCmd())
	// RootCmd.AddCommand(cyberspaceCmd())
	// rootCmd.AddCommand(poolCmd())
	// rootCmd.AddCommand(dirCmd())
	// rootCmd.AddCommand(dnsCmd())
}
