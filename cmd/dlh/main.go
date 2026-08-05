// dlh 是 .lh 解密命令：dlh <srcPath>
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gcrypto/dlh"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: dlh <srcPath>")
		os.Exit(1)
	}
	if filepath.Ext(os.Args[1]) != ".lh" {
		fmt.Fprintf(os.Stderr, "dlh: %s does not end with .lh\n", os.Args[1])
		os.Exit(1)
	}
	dst, err := dlh.DecryptFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "dlh:", err)
		os.Exit(1)
	}
	fmt.Printf("decrypted %s -> %s\n", os.Args[1], dst)
}
