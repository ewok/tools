package wrapper

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Arg struct {
	Arg []string `arg,omitempty`
}

type Cmd struct {
	Name string `name`
	Run  string `run`
	Args []Arg  `args,omitempty`
}
type Forbidden struct {
	Cmd []Cmd `cmd`
}

func ParseShort(args string) (out []string) {

	args = strings.Trim(args, "-")
	out = strings.Split(args, "")
	return
}

func ParseLong(args string) string {
	return strings.Trim(args, "--")
}

func Filter(cmdString string, forbidden Forbidden) (filtered string, ok bool, err error) {

	params := strings.Split(cmdString, " ")
	filtered = params[0]

	found := []string{}

	for _, value := range params[1:] {
		if strings.HasPrefix(value, "--") {
			found = append(found, ParseLong(value))
		} else if strings.HasPrefix(value, "-") {
			for _, v := range ParseShort(value) {
				found = append(found, v)
			}
		} else {
			found = append(found, value)
		}
	}

	isGood := true

Loop:
	for _, cmd := range forbidden.Cmd {
		if cmd.Name == params[0] {
			filtered = cmd.Run
			for _, arg := range cmd.Args {
				isGood = true
				for _, syn := range arg.Arg {
					for _, value := range found {
						if value == syn {
							isGood = false
							break
						}
					}
				}
				if isGood {
					break Loop
				}
			}
		}
	}

	if isGood {
		filtered = filtered + " " + strings.Join(params[1:], " ")
		return filtered, isGood, nil
	}

	return "", isGood, nil
}

func ReadForbidden(fileName string, forbidden *Forbidden) (err error) {

	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bs, forbidden)
	if err != nil {
		return err
	}

	return nil
}
