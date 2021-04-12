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
	"scan/config"
	"scan/controller/cli"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// portCmd represents the port command
func portCmd() *cobra.Command {
	portCmd := &cobra.Command{
		Use:   "port",
		Short: "端口扫描",
		Long:  "注,命令行参数权重大于配置文件参数, --target-file 参数权重大于--target-ips",
		RunE: func(cmd *cobra.Command, args []string) error {
			// tools.Banner()
			p := cli.NewPort(cmd, zap.L())
			if err := p.PortMain(); err != nil {
				// _ = cmd.Help()
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

	portCmd.PersistentFlags().String("protocol", "", "指定协议[tcp/udp]")
	portCmd.PersistentFlags().StringArrayP("target-ips", "i", []string{}, `目标IP列表 "192.168.1.100,192.168.2.0/24"`)
	portCmd.PersistentFlags().StringArrayP("target-ports", "p", []string{}, "要扫描的服务端口列表 22,23")
	portCmd.PersistentFlags().Int("timeout", 0, "超时")
	portCmd.PersistentFlags().Int("thread", 0, "线程")
	portCmd.PersistentFlags().Int("retry", 0, "重试次数")
	portCmd.PersistentFlags().StringP("fingerprint-file", "f", "", "服务指纹库文件( 默认内置 nmap 指纹)")
	portCmd.PersistentFlags().StringP("out-file", "o", "", "将结果输出到指定文件,默认,port.txt")
	portCmd.PersistentFlags().String("target-file", "", "目标IP列表文件 target.txt")

	return portCmd
}
