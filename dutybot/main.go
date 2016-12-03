package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nlopes/slack"
)

var (
	onDuty = slack.User{}
	botID  = ""

	configFile = flag.String("config", "config.json", "Config file")

	dutyMembers []string
)

type configStruct struct {
	Token                 string
	DutyGroup             string
	Debug                 bool
	MessageOnDuty         string
	MessageOffDuty        string
	MessageNoDuty         string
	MessageDuty           string
	MessageNotInDutyGroup string
	MessageNotDuty        string
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func updateDutyGroup(api *slack.Client, groupName string) {

	allChanenels, err := api.GetChannels(true)
	if err != nil {
		panic("Cannot get all channels")
	}

	allGroups, err := api.GetGroups(true)
	if err != nil {
		panic("Cannot get all groups")
	}

	var dutyChannel *slack.Channel
	for _, v := range allChanenels {

		if v.Name == groupName {

			dutyChannel, err = api.GetChannelInfo(v.ID)
			if err != nil {
				panic(err)
			}
			dutyMembers = dutyChannel.Members
		}
	}

	var dutyGroup *slack.Group
	for _, v := range allGroups {

		if v.Name == groupName {

			dutyGroup, err = api.GetGroupInfo(v.ID)
			if err != nil {
				panic(err)
			}

			dutyMembers = dutyGroup.Members

		}
	}
	if dutyChannel == nil && dutyGroup == nil {
		panic("Duty group not found")
	}

	fmt.Printf("Alld duties %s\n", dutyMembers)

}

func main() {
	flag.Parse()

	// Read config block
	config := &configStruct{}
	userJSON, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic("ReadFile json failed")
	}
	if err = json.Unmarshal(userJSON, &config); err != nil {
		panic("Unmarshal json failed")
	}

	api := slack.New(config.Token)
	api.SetDebug(config.Debug)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	postParams := slack.PostMessageParameters{}
	postParams.Username = "dutybot"
	postParams.Parse = "full"
	postParams.IconEmoji = ":runner:"

	// Get Duty users

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:

			case *slack.ConnectedEvent:
				fmt.Println("Infos:", ev.Info)
				fmt.Println("Connection counter:", ev.ConnectionCount)

				botID = ev.Info.User.ID

				updateDutyGroup(api, config.DutyGroup)

			case *slack.MessageEvent:

				if ev.SubType == "channel_leave" || ev.SubType == "channel_join" {
					updateDutyGroup(api, config.DutyGroup)
					fmt.Printf("%s %s\n", ev.SubType, ev.User)
				}

				if strings.Contains(ev.Text, botID) {
					if onDuty.Name != "" {
						_, _, err := rtm.PostMessage(ev.Channel, fmt.Sprintf(config.MessageDuty, onDuty.Name), postParams)
						if err != nil {
							fmt.Printf("%s\n", err)
						}
					} else {
						_, _, err := rtm.PostMessage(ev.Channel, config.MessageNoDuty, postParams)
						if err != nil {
							fmt.Printf("%s\n", err)
						}
					}

				} else if strings.Contains(ev.Text, "@onduty") {

					slackUser, err := rtm.GetUserInfo(ev.User)
					if err != nil {
						fmt.Printf("%s\n", err)
					}

					if stringInSlice(slackUser.ID, dutyMembers) {

						onDuty = *slackUser
						_, _, err = rtm.PostMessage(ev.Channel, fmt.Sprintf(config.MessageOnDuty, onDuty.Name), postParams)
						if err != nil {
							fmt.Printf("%s\n", err)
						}

					} else {

						_, _, err = rtm.PostMessage(ev.Channel, fmt.Sprintf(config.MessageNotInDutyGroup, slackUser.Name), postParams)
						if err != nil {
							fmt.Printf("%s\n", err)
						}
					}

				} else if strings.Contains(ev.Text, "@offduty") {

					slackUser, err := rtm.GetUserInfo(ev.User)
					if err != nil {
						fmt.Printf("%s\n", err)
					}

					if onDuty.ID == ev.User {
						_, _, err := rtm.PostMessage(ev.Channel, fmt.Sprintf(config.MessageOffDuty, onDuty.Name), postParams)
						if err != nil {
							fmt.Printf("%s\n", err)
						}
						onDuty = slack.User{}
					} else {
						_, _, err := rtm.PostMessage(ev.Channel, fmt.Sprintf(config.MessageNotDuty, slackUser.Name), postParams)
						if err != nil {
							fmt.Printf("%s\n", err)
						}
					}
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
			}
		}
	}
}
