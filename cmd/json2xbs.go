package cmd

import (
	"os"
	"path/filepath"
	"xbsrebuild/xbs"

	"github.com/spf13/cobra"
)

var (
	jsonFilePath string
	xbsOutPath   string
)

var json2xbsCmd = &cobra.Command{
	Use:   "j2x",
	Short: "json to xbs",
	Run: func(cmd *cobra.Command, args []string) {
		buffer, err := xbs.LoadFile(jsonFilePath)
		cobra.CheckErr(err)
		xbsBuffer, err := xbs.Json2XBS(buffer)
		cobra.CheckErr(err)
		os.MkdirAll(filepath.Dir(xbsOutPath), os.ModePerm)
		err = os.WriteFile(xbsOutPath, xbsBuffer, 0644)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(json2xbsCmd)
	json2xbsCmd.Flags().StringVarP(&jsonFilePath, "json", "i", "", "json file path")
	json2xbsCmd.Flags().StringVarP(&xbsOutPath, "xbs", "o", "", "xbs output path")
	json2xbsCmd.MarkFlagRequired("json")
	json2xbsCmd.MarkFlagRequired("xbs")
}
