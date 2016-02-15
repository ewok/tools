// Package main provides ...
package wrapper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	f := Forbidden{}
	err := ReadForbidden("wrapper.yaml", &f)
	if err != nil {
		err = ReadForbidden(os.Getenv("HOME")+"/wrapper.yaml", &f)
		if err != nil {
			err = ReadForbidden("/etc/wrapper.yaml", &f)
			if err != nil {
				fmt.Println("Config missed")
				os.Exit(1)
			}
		}
	}

	cmd, ok, err := Filter(strings.Join(os.Args, " "), f)
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
		// fmt.Printf("%s", out)
		// shellCmd.Args = cc[1:]
		// shellCmd.Stdin = &os.Stdin
		// shellCmd.Stdout = os.Stdout
		// shellCmd.Stderr = &os.Stderr
		// shellCmd.Run()
	}
}
