// cmd.go 方便被应用 import
// hbhsid.HIDE 的初始化需要应用自己处理
package cmd

import (
	"fmt"

	"scan/pkg/libc/hbhsid"

	"github.com/spf13/cobra"
)

func encodeCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "encode",
		Short: "encode hbhsid",
		RunE: func(cmd *cobra.Command, args []string) error {
			u32, err := cmd.Flags().GetUint32("uint32")
			if err != nil {
				return err
			}

			id := hbhsid.New(u32)
			fmt.Printf("%s\n", id.String())

			return nil
		},
	}

	cmd.Flags().Uint32("uint32", 0, "需要转换的 uint32")
	_ = cmd.MarkFlagRequired("uint32")

	return &cmd
}

func decodeCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "decode",
		Short: "decode hbhsid",
		RunE: func(cmd *cobra.Command, args []string) error {
			_id, err := cmd.Flags().GetString("id")
			if err != nil {
				return err
			}

			if id, e := hbhsid.ParseFromString(_id); e != nil {
				fmt.Printf("解码错误: %s\n", e.Error())
			} else {
				fmt.Println(id.Origin())
			}
			return nil
		},
	}

	cmd.Flags().String("id", "", "需要解码的 hbhsid 字符串")
	_ = cmd.MarkFlagRequired("id")
	return &cmd
}

func hbhsidCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "hbhsid",
		Short: "hbhsid tools",
		Long:  "hbhsid 编解码工具",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hbhSId called")
		},
	}

	cmd.AddCommand(encodeCmd(), decodeCmd())
	return &cmd
}

func GatherCmd(root *cobra.Command) {
	root.AddCommand(hbhsidCmd())
}
