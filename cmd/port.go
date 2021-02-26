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
		Long:  "端口扫描器",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := cli.NewPort(cmd, zap.L())
			if err := p.PortMain(); err != nil {
				zap.L().Error("端口扫描失败", zap.Error(err))
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

	portCmd.PersistentFlags().String("protocol", "tcp", "扫描协议[tcp/udp]")
	portCmd.PersistentFlags().StringArrayP("target-ips", "i", []string{}, `服务器IP ["192.168.1.100", "192.168.1.11"]`)
	portCmd.PersistentFlags().StringArrayP("target-ports", "p", []string{}, `需要扫描的端口列表 ["22", "23"]`)
	portCmd.PersistentFlags().Int("timeout", 0, "超时时间")
	portCmd.PersistentFlags().Int("thread", 0, "扫描线程")
	portCmd.PersistentFlags().Int("retry", 0, "重试次数")
	portCmd.PersistentFlags().StringP("fingerprint-file", "f", "", "规则或者指纹文件")

	return portCmd
}
