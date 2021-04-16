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

// blastingCmd represents the blasting command
func blastingCmd() *cobra.Command {
	blastingCmd := &cobra.Command{
		Use:   "brute",
		Short: "口令喷射",
		RunE: func(cmd *cobra.Command, args []string) error {
			// tools.Banner()
			if len(os.Args) == 2 {
				_ = cmd.Help()
				return nil
			}
			p := cli.NewBlasting(cmd, zap.L())
			if err := p.BlastingMain(); err != nil {
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
	blastingCmd.PersistentFlags().StringP("pass-file", "l", "", "密码字典")
	blastingCmd.PersistentFlags().StringP("port", "p", "", "服务端口(如非默认,请自行手工指定)")
	blastingCmd.PersistentFlags().IntP("delay", "d", 0, "延迟,默认1")
	blastingCmd.PersistentFlags().IntP("thread", "t", 0, "线程,默认20")
	blastingCmd.PersistentFlags().IntP("timeout", "m", 0, "超时,默认1")
	blastingCmd.PersistentFlags().IntP("retry", "r", 0, "重试次数,默认1")
	blastingCmd.PersistentFlags().BoolP("scan-port", "n", false, "爆破前是否进行端口扫描")
	blastingCmd.PersistentFlags().StringArrayP("services", "s", []string{},
		`指定要爆破的服务 "ssh,ftp,mssql,mysql,redis,postgresql,http_basic,tomcat,telnet"`)
	blastingCmd.PersistentFlags().StringP("path", "b", "", `http_basic 路径 "/login"`)
	blastingCmd.PersistentFlags().StringP("tomcat-path", "a", "", `tomcat路径 "/manager"`)
	blastingCmd.PersistentFlags().StringP("out-file", "o", "", "输出文件,blasting.txt")

	return blastingCmd
}
