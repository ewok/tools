// Package main provides ...
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ewok/tools/wrapper/wrapper"
)

func main() {
	f := wrapper.Forbidden{}
	err := wrapper.ReadForbidden("wrapper.yaml", &f)
	if err != nil {
		err = wrapper.ReadForbidden(os.Getenv("HOME")+"/wrapper.yaml", &f)
		if err != nil {
			err = wrapper.ReadForbidden("/etc/wrapper.yaml", &f)
			if err != nil {
				fmt.Println("Config missed")
				os.Exit(1)
			}
		}
	}

	cmd, ok, err := wrapper.Filter(strings.Join(os.Args, " "), f)
	if err != nil {
		panic(err)
	}

	if !ok {
		fmt.Println("Forbidden cmd")
	} else {
		fmt.Println("Command allowed: " + cmd)
		cc := []string{"-c"}
		cc = append(cc, strings.Split(cmd, " ")...)
		err := exec.Command("sh", cc[0:]...).Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
