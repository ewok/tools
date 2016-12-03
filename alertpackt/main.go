// Package main provides ...
package main

import "github.com/ewok/tools/packfreebook"

import "github.com/xconstruct/go-pushbullet"

func main() {

	pb := pushbullet.New("<TOKEN>")
	devs, err := pb.Devices()
	if err != nil {
		panic(err)
	}

	for n := range devs {
		if devs[n].Nickname != "" {
			err = pb.PushNote(devs[n].Iden, "PacktFreeBook", packfreebook.PackFreeBook())
			if err != nil {
				panic(err)
			}
		}
	}
}
