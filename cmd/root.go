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
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "scan",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		banner()
		_ = cmd.Help()
	},
}

func banner() {
	fmt.Printf(`
____________________________________________________________        
 ____________________________________________________________       
  ____________________________________________________________      
   __/\\\\\\\\\\______/\\\\\\\\___/\\\\\\\\\______/\\/\\\\\\___     
    _\/\\\//////_____/\\\//////___\////////\\\____\/\\\////\\\__    
     _\/\\\\\\\\\\___/\\\____________/\\\\\\\\\\___\/\\\__\//\\\_   
      _\////////\\\__\//\\\__________/\\\/////\\\___\/\\\___\/\\\_  
       __/\\\\\\\\\\___\///\\\\\\\\__\//\\\\\\\\/\\__\/\\\___\/\\\_ 
        _\//////////______\////////____\////////\//___\///____\///__
` + "\n")
}

// Execute is the main cobra method
func Execute() {
	var cancel context.CancelFunc
	mainContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan:
			// caught CTRL+C
			fmt.Println("\n[!] Keyboard interrupt detected, terminating.")
			cancel()
		case <-mainContext.Done():
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		// Leaving this in results in the same error appearing twice
		// Once before and once after the help output. Not sure if
		// this is going to be needed to output other errors that
		// aren't automatically outputted.
		// fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 全局配置
	rootCmd.PersistentFlags().StringP("url", "u", "", "需要扫描的ip或url")
	rootCmd.PersistentFlags().StringP("config", "c", "", "配置文件路径")

	// 新增命令
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(poolCmd())
	rootCmd.AddCommand(dirCmd())
	rootCmd.AddCommand(dnsCmd())
	rootCmd.AddCommand(portCmd())
}
