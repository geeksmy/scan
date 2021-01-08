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
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = cobra.Command{
	Use: "scan",
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

func init() {
	// 全局配置
	RootCmd.PersistentFlags().StringP("config", "c", "", "配置文件路径")
	RootCmd.PersistentFlags().StringP("url", "u", "", "需要扫描的ip或url")

	// 新增命令
	RootCmd.AddCommand(versionCmd())
}
