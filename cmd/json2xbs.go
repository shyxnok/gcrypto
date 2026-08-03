package cmd

import (
	"gcrypto/xbs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	iFilePath string
	oFilePath string
)

var json2xbsCmd = &cobra.Command{
	Use:   "j2x",
	Short: "json to xbs",
	Run: func(cmd *cobra.Command, args []string) {
		buffer, err := xbs.LoadFile(iFilePath)
		cobra.CheckErr(err)
		xbsBuffer, err := xbs.Json2XBS(buffer)
		cobra.CheckErr(err)
		os.MkdirAll(filepath.Dir(oFilePath), os.ModePerm)
		err = os.WriteFile(oFilePath, xbsBuffer, 0644)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(json2xbsCmd)
	json2xbsCmd.Flags().StringVarP(&iFilePath, "json", "i", "", "json file path")
	json2xbsCmd.Flags().StringVarP(&oFilePath, "xbs", "o", "", "xbs output path")
	json2xbsCmd.MarkFlagRequired("json")
	json2xbsCmd.MarkFlagRequired("xbs")
}
