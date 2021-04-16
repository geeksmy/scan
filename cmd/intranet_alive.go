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

// intranetAliveCmd represents the intranetAlive command
func intranetAliveCmd() *cobra.Command {
	intranetAliveCmd := &cobra.Command{
		Use:   "survive",
		Short: "存活段搜集[管理权限]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(os.Args) == 2 {
				_ = cmd.Help()
				return nil
			}
			p := cli.NewIntranetAlive(cmd, zap.L())
			if err := p.IntranetAliveMain(); err != nil {
				return err
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			return nil
		},
	}

	intranetAliveCmd.PersistentFlags().StringP("target", "i", "", "目标 (可以多个, 以逗号隔开, 入1,254)")
	intranetAliveCmd.PersistentFlags().IntP("timeout", "m", 0, "超时,默认1")
	intranetAliveCmd.PersistentFlags().IntP("thread", "t", 0, "线程,默认20")
	intranetAliveCmd.PersistentFlags().IntP("retry", "r", 0, "重试次数,默认1")
	intranetAliveCmd.PersistentFlags().Float32P("delay", "d", 0, "延迟,默认0")
	intranetAliveCmd.PersistentFlags().StringP("out-file", "o", "", "将结果输出到指定文件")
	return intranetAliveCmd
}
