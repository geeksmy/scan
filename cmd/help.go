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
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// helpCmd represents the help command
func helpCmd() *cobra.Command {
	helpCommand := &cobra.Command{
		Use:               "help [command]",
		Short:             "模块使用帮助",
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(c *cobra.Command, args []string) error {
			cmd, args, e := c.Root().Find(args)
			if cmd == nil || e != nil || len(args) > 0 {
				return errors.Errorf("未知的模板使用帮助: %v", strings.Join(args, " "))
			}

			helpFunc := cmd.HelpFunc()
			helpFunc(cmd, args)
			return nil
		},
	}

	return helpCommand
}
