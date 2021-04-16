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
	"os"

	"scan/config"
	"scan/controller/cli"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

// portCmd represents the port command
func portCmd() *cobra.Command {
	portCmd := &cobra.Command{
		Use:   "port",
		Short: "端口扫描",
		RunE: func(cmd *cobra.Command, args []string) error {
			// tools.Banner
			if len(os.Args) == 2 {
				_ = cmd.Help()
				return nil
			}
			p := cli.NewPort(cmd, zap.L())
			if err := p.PortMain(); err != nil {
				return err
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			// dao.InitDB(config.C)
			return nil
		},
	}

	portCmd.PersistentFlags().StringP("protocol", "x", "", "指定协议[tcp/udp]")
	portCmd.PersistentFlags().StringArrayP("target-ips", "i", []string{}, `目标IP列表 "192.168.1.100,192.168.2.0/24"`)
	portCmd.PersistentFlags().StringArrayP("target-ports", "p", []string{}, "要扫描的服务端口列表 22,23")
	portCmd.PersistentFlags().IntP("timeout", "m", 0, "超时,默认1")
	portCmd.PersistentFlags().IntP("thread", "t", 0, "线程,默认20")
	portCmd.PersistentFlags().IntP("retry", "r", 0, "重试次数,默认1")
	portCmd.PersistentFlags().StringP("fingerprint-file", "f", "", "服务指纹库文件( 默认内置 nmap 指纹)")
	portCmd.PersistentFlags().StringP("out-file", "o", "", "将结果输出到指定文件,默认,port.txt")
	portCmd.PersistentFlags().StringP("target-file", "b", "", "目标IP列表文件 target.txt")

	return portCmd
}
