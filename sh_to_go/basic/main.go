package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmdstr := "ip route | grep default"
	out, err := exec.Command("sh", "-c", cmdstr).Output()
	if err != nil {
		return
	}
	fmt.Printf("%s", out)
}
