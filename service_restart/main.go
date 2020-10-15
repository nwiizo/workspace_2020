package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("start time: ", time.Now().Format("15:04:05"))
	out, err := exec.Command("service", "docker", "status").Output()
	if err != nil {
		fmt.Errorf("ERROR")
	}
	fmt.Printf("out:%s", out)
}
