package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"os"
)

const (
	OutputPin = "40"
	BotName   = "DiscordLedController"
)

func checkError(err error, exit bool) {
	if err != nil && exit == false {
		log.Println(fmt.Sprintf("ERROR : %s", err))
	}

	if err != nil && exit == true {
		log.Fatal(fmt.Sprintf("ERROR : %s", err))
	}
}

func createBot(discord *discordgo.Session) *gobot.Robot {
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, OutputPin)

	robot := gobot.NewRobot(BotName,
		[]gobot.Connection{r},
		[]gobot.Device{led},
		func() {
			discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
				messageContent := m.Content

				if messageContent == "start" {
					err := led.On()
					checkError(err, false)
				}

				if messageContent == "stop" {
					err := led.Off()
					checkError(err, false)
				}
			})
		},
	)

	return robot
}

func main() {
	if err := godotenv.Load(); err != nil {
		checkError(err, true)
	}

	discordToken := os.Getenv("DISCORD_TOKEN")

	if discordToken == "" {
		emptyDiscordTokenError := fmt.Errorf("discord token invalid")
		checkError(emptyDiscordTokenError, true)
	}

	discord, err := discordgo.New("Bot " + discordToken)
	checkError(err, true)

	err = discord.Open()
	checkError(err, true)

	robot := createBot(discord)

	err = robot.Start()
	checkError(err, true)
}
