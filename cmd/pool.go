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

	"github.com/geeksmy/go-lib/redis"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// poolCmd represents the pool command
func poolCmd() *cobra.Command {
	poolCmd := &cobra.Command{
		Use:   "pool",
		Short: "是否开启代理池",
		Long:  "是否开启代理池,开启代理池需要启动redis",
		RunE: func(cmd *cobra.Command, args []string) error {
			zap.L().Info("待实现.....")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cfgFile)
			if err := redis.Connect(); err != nil {
				return err
			}
			return nil
		},
	}

	redis.BindPflag(poolCmd.Flags(), "redis")

	return poolCmd
}
