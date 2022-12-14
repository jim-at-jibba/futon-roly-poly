package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
	"github.com/slack-go/slack"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func Scrape() {

	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnHTML(".ty-product-block__price-actual", func(e *colly.HTMLElement) {
		current := 95
		fullPrice := e.ChildText(".ty-price-num")
		_, i := utf8.DecodeRuneInString(fullPrice)
		marks, err := strconv.ParseFloat(fullPrice[i:], 64)
		if err != nil {
			fmt.Println("Errot", err)
		}

		fmt.Println(marks)
		if marks != float64(current) {
			SendMsg("Roly Poly has changed price", os.Getenv("ELLA_PHONE_NUMBER"))
		}

	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.Visit("https://www.futoncompany.co.uk/shop-by-product/new/roly-poly-pebble-grey-coast-weave.html")
}

func sendSlackMessage(message string, color string) error {
	token := os.Getenv("SLACK_AUTH_TOKEN")
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	// Create a new client to slack by giving token
	// Set debug to true while developing
	client := slack.New(token, slack.OptionDebug(true))
	attachment := slack.Attachment{
		Pretext: "Roly Poly Futon",
		Text:    message,
		Color:   color,
	}
	_, timestamp, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		return err
	}

	fmt.Printf("Message sent at %s", timestamp)
	return nil

}

func SendMsg(msg string, to string) {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(msg)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("SMS sent successfully!")
	}
}

func main() {
	err := sendSlackMessage("Starting the roly poly scraper", "#36a64f")

	Scrape()

	if err != nil {
		panic(err)
	}
}
