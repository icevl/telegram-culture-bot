package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const ConfigFile = "config.json"

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetOutput(os.Stdout)
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

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	channels := loadConfig()
	openAI := NewOpenAI(aiToken)
	saveDataChan := make(chan SaveData)

	for _, channel := range channels {
		go scheduler(bot, openAI, channel, saveDataChan)
	}

	go func() {
		for data := range saveDataChan {
			saveChannelNextTime(data.Channel, data.Next)
		}
	}()

	go eventListener(bot)

	<-make(chan bool)
}

func loadConfig() []Channel {
	data := []Channel{}
	file, err := os.ReadFile(ConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func scheduler(bot *tgbotapi.BotAPI, openAI *OpenAI, channel Channel, saveNextChan chan<- SaveData) {
	log.Printf("Starting scheduler for: '%s'", channel.Title)

	for {
		unixTime := time.Now().Unix()

		if unixTime < channel.NextTime {
			time.Sleep(1 * time.Minute)
			continue
		}

		randomMinutes := rand.Intn(channel.MaxMins-channel.MinMins+1) + channel.MinMins
		randomDuration := time.Duration(randomMinutes) * time.Minute

		gptText, ok := openAI.GetAnswer(channel.Prompt)
		if !ok {
			time.Sleep(10 * time.Second)
			continue
		}

		data := strings.Split(gptText, "|")
		text := gptText

		if len(data) == 3 {
			emoji, title, body := data[2], escapeQuotes(data[1]), escapeQuotes(data[0])
			text = fmt.Sprintf("%s *%s*\n\n%s", emoji, title, body)
		}

		if channel.Image != "" {
			ok := sendWithPhoto(bot, openAI, channel, text)
			if !ok {
				continue
			}
		} else {
			ok := sendTextOnly(bot, openAI, channel, text)
			if !ok {
				continue
			}
		}

		channel.NextTime = time.Now().Add(randomDuration).Unix()

		log.Printf("Next message for '%s' will be sent in %d minutes", channel.Title, randomMinutes)

		saveData := SaveData{
			Channel: &channel,
			Next:    channel.NextTime,
		}

		saveNextChan <- saveData
	}
}

func eventListener(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
	}
}

func sendTextOnly(bot *tgbotapi.BotAPI, openAI *OpenAI, channel Channel, text string) bool {
	msg := tgbotapi.NewMessage(channel.ChatID, text)
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)

	if err != nil {
		log.Printf("Error sending message: %s", err)
		return false
	}

	return true
}

func sendWithPhoto(bot *tgbotapi.BotAPI, openAI *OpenAI, channel Channel, text string) bool {
	picturePrompt := channel.Image

	if channel.Image == "from_prompt_result" {
		picturePrompt = text
	}

	if picturePrompt == "" {
		log.Println("No picture prompt provided")
		return false
	}

	file, ok := openAI.GetImage(picturePrompt)
	if !ok {
		return false
	}

	imageData, _ := os.ReadFile(file)
	blob := tgbotapi.FileBytes{Name: "image.png", Bytes: imageData}

	msg := tgbotapi.NewPhoto(channel.ChatID, blob)
	msg.ParseMode = "markdown"
	msg.Caption = text

	_, err := bot.Send(msg)
	os.Remove(file)

	if err != nil {
		log.Printf("Error sending message: %s", err)
		return false
	}

	return true
}

func saveChannelNextTime(channel *Channel, next int64) {
	channels := loadConfig()

	for i, ch := range channels {
		if ch.Prompt == channel.Prompt {
			channels[i].NextTime = next
		}
	}

	file, _ := json.MarshalIndent(channels, "", " ")
	_ = os.WriteFile(ConfigFile, file, 0644)
}

func escapeQuotes(text string) string {
	re := regexp.MustCompile(`^"(.*)"$`)
	return re.ReplaceAllString(text, "$1")
}
