// elh 是 .lh 加密命令：elh <srcPath>
package main

import (
	"fmt"
	"os"

	"gcrypto/elh"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: elh <srcPath>")
		os.Exit(1)
	}
	dst, err := elh.EncryptFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "elh:", err)
		os.Exit(1)
	}
	fmt.Printf("encrypted %s -> %s\n", os.Args[1], dst)
}
