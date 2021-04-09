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

// blastingCmd represents the blasting command
func blastingCmd() *cobra.Command {
	blastingCmd := &cobra.Command{
		Use:   "brute",
		Short: "口令喷射",
		Long:  "口令喷射工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			// tools.Banner()
			p := cli.NewBlasting(cmd, zap.L())
			if err := p.BlastingMain(); err != nil {
				_ = cmd.Help()
				return err
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			return nil
		},
	}

	blastingCmd.PersistentFlags().StringP("target-host", "i", "", "目标文件")
	blastingCmd.PersistentFlags().StringP("user-file", "u", "", "用户名字典")
	blastingCmd.PersistentFlags().StringP("pass-file", "p", "", "密码字典")
	blastingCmd.PersistentFlags().String("port", "", "服务端口(如目标用的非默认端口,则需自行手工指定)")
	blastingCmd.PersistentFlags().Int("delay", 0, "延迟")
	blastingCmd.PersistentFlags().Int("thread", 0, "线程")
	blastingCmd.PersistentFlags().Int("timeout", 0, "超时")
	blastingCmd.PersistentFlags().Int("retry", 0, "重试次数")
	blastingCmd.PersistentFlags().Bool("scan-port", false, "爆破前是否进行端口扫描")
	blastingCmd.PersistentFlags().StringArrayP("services", "s", []string{},
		`指定要爆破的服务 "ssh","ftp","mssql","mysql","redis","postgresql","http_basic","tomcat","telnet"`)
	blastingCmd.PersistentFlags().String("path", "", `http_basic 路径 "/login"`)
	blastingCmd.PersistentFlags().String("tomcat-path", "", `tomcat路径 "/manager"`)
	blastingCmd.PersistentFlags().StringP("out-file", "o", "", "输出文件,blasting.txt")

	return blastingCmd
}
