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

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// portCmd represents the port command
func portCmd() *cobra.Command {
	portCmd := &cobra.Command{
		Use:   "dns",
		Short: "端口扫描",
		Long:  "端口扫描器",
		RunE: func(cmd *cobra.Command, args []string) error {
			zap.L().Info("待实现.....")
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Init(cmd.Flags().Lookup("config").Value.String())
			return nil
		},
	}

	return portCmd
}
