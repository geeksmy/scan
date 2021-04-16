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

// webFingerprintCmd represents the webFingerprint command
func webFingerprintCmd() *cobra.Command {
	webFingerprintCmd := &cobra.Command{
		Use:   "web",
		Short: "web服务探测",
		RunE: func(cmd *cobra.Command, args []string) error {
			// tools.Banner()
			if len(os.Args) == 2 {
				_ = cmd.Help()
				return nil
			}
			p := cli.NewWebFingerprint(cmd, zap.L())
			if err := p.WebFingerprintMain(); err != nil {
				return err
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			return nil
		},
	}

	webFingerprintCmd.PersistentFlags().StringP("target-urls", "u", "", "目标IP列表文件")
	webFingerprintCmd.PersistentFlags().IntP("thread", "t", 0, "线程,默认20")
	webFingerprintCmd.PersistentFlags().IntP("timeout", "m", 0, "超时,默认1")
	webFingerprintCmd.PersistentFlags().IntP("retry", "r", 0, "重试次数,默认1")
	webFingerprintCmd.PersistentFlags().StringP("out-file", "o", "", "输出文件 web-fingerprint.txt")
	webFingerprintCmd.PersistentFlags().StringArrayP("target-ports", "p", []string{}, "要扫描的Web服务端口列表 80,443")
	webFingerprintCmd.PersistentFlags().StringP("fingerprint-file", "f", "", "Web指纹文件")

	return webFingerprintCmd
}
