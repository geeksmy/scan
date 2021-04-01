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

// webFingerprintCmd represents the webFingerprint command
func webFingerprintCmd() *cobra.Command {
	webFingerprintCmd := &cobra.Command{
		Use:   "web-fingerprint",
		Short: "web指纹识别",
		Long:  "web 指纹识别工具",
		RunE: func(cmd *cobra.Command, args []string) error {
			// tools.Banner()
			p := cli.NewWebFingerprint(cmd, zap.L())
			if err := p.WebFingerprintMain(); err != nil {
				_ = cmd.Help()
				return nil
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			return nil
		},
	}

	webFingerprintCmd.PersistentFlags().StringP("target-urls", "u", "", "目标 url 文件")
	webFingerprintCmd.PersistentFlags().Int("thread", 0, "线程")
	webFingerprintCmd.PersistentFlags().Int("timeout", 0, "超时时间")
	webFingerprintCmd.PersistentFlags().Int("retry", 0, "重试次数 必须>=1")
	webFingerprintCmd.PersistentFlags().StringP("out-file", "o", "", "输出文件 web-fingerprint.txt")
	webFingerprintCmd.PersistentFlags().StringArrayP("target-ports", "p", []string{}, `需要扫描的端口列表 ["80", "443"]`)
	webFingerprintCmd.PersistentFlags().StringP("fingerprint-file", "f", "", "规则或者指纹文件")

	return webFingerprintCmd
}
