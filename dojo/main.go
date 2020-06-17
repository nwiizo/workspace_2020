package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"sort"
	"strings"
)

type Result struct {
	Pref  string
	Count int
}

type Results []Result

func (r Result) String() string {
	return fmt.Sprintf("%d %s\n", r.Count, r.Pref)
}

func (r Results) Less(i, j int) bool {
	return r[i].Pref < r[j].Pref
}

func (r Results) Len() int {
	return len(r)
}

func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func main() {
	if terminal.IsTerminal(0) {
		fmt.Println("not input")
	} else {
		b := true
		scanner := bufio.NewScanner(os.Stdin)
		for b {
			b = scanner.Scan()
			t := scanner.Text()
			if t == "" {
				break
			}
			// 文字入力の処理,分割
			t1 := strings.ReplaceAll(strings.ReplaceAll(t, ".", ""), ",", "")
			slice := strings.Split(t1, " ")
			prefCounts := make(map[string]int)
			for _, v := range slice {
				prefCounts[v] += 1
			}
			newThreeLine(bestThree(prefCounts))
		}
	}
}

func strSorted(v1, v2 string) bool {
	vsort := []string{v1, v2}
	vool := sort.StringsAreSorted(vsort)
	return vool
}

func newThreeLine(v1, v2, v3 string) {
	fmt.Println(v1)
	fmt.Println(v2)
	fmt.Println(v3)
}

func bestThree(numbers map[string]int) (string, string, string) {
	var maxNumber, tNumber, sNumber int
	var maxStr, tStr, sStr string
	for a, n := range numbers {
		if (n > maxNumber && n >= tNumber && n > sNumber) || (n == maxNumber && strSorted(a, maxStr)) {
			sNumber = tNumber
			tNumber = maxNumber
			maxNumber = n
			sStr = tStr
			tStr = maxStr
			maxStr = a
			continue
		}
		if n > tNumber && n > sNumber || (n == tNumber && strSorted(a, tStr)) {
			sNumber = tNumber
			tNumber = n
			sStr = tStr
			tStr = a
			continue
		}
		if n > sNumber || (n == tNumber && strSorted(a, sStr)) {
			sNumber = n
			sStr = a
			continue
		}
	}
	return maxStr, tStr, sStr
}
