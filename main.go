package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	minMinutes = 120
	maxMinutes = 180
)

func init() {
	godotenv.Load(".env")
}

func main() {
	tgToken := os.Getenv("BOT_TOKEN")
	aiToken := os.Getenv("OPENAI_TOKEN")

	if tgToken == "" {
		log.Panic("No telegram token provided")
	}

	if aiToken == "" {
		log.Panic("No open ai token provided")
	}

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	openAI := NewOpenAI(aiToken)

	go scheduler(bot, openAI)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	<-make(chan bool)
}

func scheduler(bot *tgbotapi.BotAPI, openAI *OpenAI) {
	channelIdent := os.Getenv("CHANNEL_ID")
	id, _ := strconv.Atoi(channelIdent)
	channelID := int64(id)

	for {
		randomMinutes := rand.Intn(maxMinutes-minMinutes+1) + minMinutes
		randomDuration := time.Duration(randomMinutes) * time.Minute

		gptText, ok := openAI.GetAnswer("Сгенерируй интересный факт о какой нибудь стране и ее культуре и что бы не только о Японии. В конце текста через разделитель '|' напиши страну о которой идет речь, после чего поставь разделитель '|' и emoji для страны")
		if !ok {
			time.Sleep(10 * time.Second)
			continue
		}

		data := strings.Split(gptText, "|")

		emoji, country, fact := data[2], data[1], data[0]
		text := fmt.Sprintf("%s *%s*\n\n%s", emoji, country, fact)

		msg := tgbotapi.NewMessage(channelID, text)
		msg.ParseMode = "markdown"

		bot.Send(msg)

		fmt.Printf("Next after %d mins", randomMinutes)
		fmt.Println()

		time.Sleep(randomDuration)
	}
}
