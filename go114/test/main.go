package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
)

func main() {
	filename := "./file.txt"

	// ファイルオープン
	fp, err := os.Open(filename)
	if err != nil {
		// エラー処理
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	var w sync.WaitGroup
	for scanner.Scan() {
		// ここで一行ずつ処理
		w.Add(1)
		go func() {
			defer w.Done()
			checksum256 := sha256.Sum256(scanner.Bytes())
			fmt.Printf("SHA256 checksum: %x\n", checksum256)
		}()
		w.Wait()
	}

	if err = scanner.Err(); err != nil {
		// エラー処理
	}
}
