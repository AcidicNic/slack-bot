package slack

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"

	"github.com/slack-go/slack"
	"github.com/tidwall/gjson"
)

/*
   TODO: Change @BOT_NAME to the same thing you entered when creating your Slack application.
   NOTE: command_arg_1 and command_arg_2 represent optional parameteras that you define
   in the Slack API UI
*/
const helpMessage = "Hey guy, wondering what commands I can understand?\n\n*Get Random r/ProgrammerHumor Post:*\n\t_@boi <rp/random post/reddit>_\n\n*You can also try:*\n\t_Saying hello or saying you love me_"

const rand_reddit_url = "https://www.reddit.com/r/ProgrammerHumor/random/.json"

// Global slices!
func getHello() []string {
	return []string{"hey", "hello", "hi", "yo"}
}
func getLove() []string {
	return []string{"i love you", "<3", "ily", "love you", "i love u", "love u"}
}
func getReddit() []string {
	return []string{"rand post", "rand", "random", "rp", "random post", "randpost", "random post", "reddit"}
}

/*
   CreateSlackClient sets up the slack RTM (real-timemessaging) client library,
   initiating the socket connection and returning the client.
   DO NOT EDIT THIS FUNCTION. This is a fully complete implementation.
*/
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

/*
   RespondToEvents waits for messages on the Slack client's incomingEvents channel,
   and sends a response when it detects the bot has been tagged in a message with @<botTag>.

   EDIT THIS FUNCTION IN THE SPACE INDICATED ONLY!
*/
func RespondToEvents(slackClient *slack.RTM) {
	for msg := range slackClient.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)

			// TODO: Make your bot do more than respond to a help command. See notes below.
			// Make changes below this line and add additional funcs to support your bot's functionality.
			// sendHelp is provided as a simple example. Your team may want to call a free external API
			// in a function called sendResponse that you'd create below the definition of sendHelp,
			// and call in this context to ensure execution when the bot receives an event.

			// START SLACKBOT CUSTOM CODE
			// ===============================================================
			sendResponse(slackClient, message, ev.Channel)
			sendHelp(slackClient, message, ev.Channel)
			// ===============================================================
			// END SLACKBOT CUSTOM CODE
		default:

		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// sendResponse is NOT unimplemented --- write code in the function body to complete!

func sendResponse(slackClient *slack.RTM, message, slackChannel string) {
	command := strings.ToLower(message)
	println("[RECEIVED] sendResponse:", command)
	// START SLACKBOT CUSTOM CODE
	// ===============================================================
	response := ""
	if contains(command, getHello()) {
		response = "Hey pal!"
	}
	if contains(command, getLove()) {
		response = "I love you too!"
	}
	if contains(command, getReddit()) {
		data, err := getRedditPost()
		if err != nil {
			response = "Oops, something went wrong! Try again."
		} else {
			title := gjson.Get(data, "0.data.children.0.data.title").String()
			text := gjson.Get(data, "0.data.children.0.data.selftext").String()
			postUrl := "www.reddit.com" + gjson.Get(data, "0.data.children.0.data.permalink").String()

			// if [text post], else [image post]
			if text != "" {
				response = fmt.Sprintf("*Random r/ProgrammerHumor post:*\n%s\n```%s```\n> source: %s", title, text, postUrl)
			} else {
				imgUrl := gjson.Get(data, "0.data.children.0.data.url").String()
				response = fmt.Sprintf("*Random r/ProgrammerHumor post:*\n*_<%s|%s>_*\n> source: %s", imgUrl, title, postUrl)
			}
		}
	}

	if response != "" {
		slackClient.SendMessage(slackClient.NewOutgoingMessage(response, slackChannel))
	}
	// ===============================================================
	// END SLACKBOT CUSTOM CODE
}

// Helper func to check if string slice "list" contains the string "msg"!
func contains(msg string, list []string) bool {
   for _, item := range list {
      if item == msg {
         return true
      }
   }
   return false
}

func getRedditPost() (string, error) {
    req, err := http.NewRequest("GET", "https://www.reddit.com/r/ProgrammerHumor/random/.json", nil)
    if err != nil {
        return "", fmt.Errorf("")
    }
    req.Header.Set("User-agent", "boi the Slack Bot")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("")
    }
    return string(body), nil
}
