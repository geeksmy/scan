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

// passgenCmd represents the passgen command
func passgenCmd() *cobra.Command {
	passgenCmd := &cobra.Command{
		Use:   "passgen",
		Short: "密码生成",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := cli.NewPassGen(cmd, zap.L())
			if err := p.PassGenMain(); err != nil {
				return err
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			config.Init(cfgFile)
		},
	}

	passgenCmd.PersistentFlags().StringP("year", "y", "", "年份(可以同时有多个,以逗号隔开,如,2019,2020,2021)")
	passgenCmd.PersistentFlags().String("domain-name", "", "目标域名简称,或者谐音变形(可以有多个,以逗号隔开,如,sangfor,sinfor)")
	passgenCmd.PersistentFlags().String("domain", "", "完整域名 (如,www.baidu.com)")
	passgenCmd.PersistentFlags().StringP("device", "d", "", "设备名(可以有多个,以逗号隔开,如,dell,hp)")
	passgenCmd.PersistentFlags().IntP("length", "l", 0, "生成的密码在前中后分别包含几个特殊字符")
	passgenCmd.PersistentFlags().StringP("out-file", "o", "", "将结果输出到指定文件,默认,pass.txt")
	return passgenCmd
}
