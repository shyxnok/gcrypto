/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"
	"xbsrebuild/xbs"

	"github.com/spf13/cobra"
)

var (
	xbsFilePath string
	jsonOutPath string
)

// xbs2jsonCmd represents the xbs2json command
var xbs2jsonCmd = &cobra.Command{
	Use:   "x2j",
	Short: "xbs to json",
	Run: func(cmd *cobra.Command, args []string) {
		buffer, err := xbs.LoadFile(xbsFilePath)
		cobra.CheckErr(err)
		jsonBuffer, err := xbs.XBS2Json(buffer)
		cobra.CheckErr(err)
		os.MkdirAll(filepath.Dir(jsonOutPath), os.ModePerm)
		err = os.WriteFile(jsonOutPath, jsonBuffer, 0644)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(xbs2jsonCmd)
	xbs2jsonCmd.Flags().StringVarP(&xbsFilePath, "xbs", "i", "", "xbs file path")
	xbs2jsonCmd.Flags().StringVarP(&jsonOutPath, "json", "o", "", "json output path")
	xbs2jsonCmd.MarkFlagRequired("xbs")
	xbs2jsonCmd.MarkFlagRequired("json")
}
