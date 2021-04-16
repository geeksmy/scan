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
	"scan/internal/dao"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

// dnsCmd represents the dns command
func dnsCmd() *cobra.Command {
	dnsCmd := &cobra.Command{
		Use:   "dns",
		Short: "子域名扫描",
		Long:  "子域名扫描器，支持[文件，数据库，枚举]这三种扫描模式，默认使用文件模式",
		RunE: func(cmd *cobra.Command, args []string) error {
			zap.L().Info("待实现.....")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			switch cmd.Flags().Lookup("mode").Value.String() {
			case "db":
				dao.InitDB(config.C)
			case "enumerate":
				zap.L().Info("初始化枚举相关参数")
			default:
				zap.L().Info("不需要初始化")
			}
			return nil
		},
	}

	dnsCmd.PersistentFlags().StringP("mode", "m", "file", "支持三种模式, 文件:file, 数据库:db, 枚举: enumerate")

	return dnsCmd
}
