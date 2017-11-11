package main

import (
	"os"
	"fmt"
	"strings"
	"os/signal"
	"syscall"
	"github.com/sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	r "gopkg.in/gorethink/gorethink.v3"
)

var log = logrus.New()

var (
	session *r.Session
	discordSession *discordgo.Session
)

func main() {
	log.Formatter = new(logrus.TextFormatter)
	log.Info("OverStatsDiscord 1.0 started!")

	var err error

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN env variable not specified!")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session, ", err)
	}

	// Database pool init
	go InitConnectionPool()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection, ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	defer dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	discordSession = s

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// userId for logger
	commandLogger := log.WithFields(logrus.Fields{"user_id": m.Author.ID})

	if strings.HasPrefix(m.Content, "/start") {
		commandLogger.Info("command /start triggered")
		go StartCommand(s, m)
	}

	if strings.HasPrefix(m.Content, "/donate") {
		commandLogger.Info("command /donate triggered")
		go DonateCommand(s, m)
	}

	if strings.HasPrefix(m.Content, "/save") {
		commandLogger.Info("command /save triggered")
		go SaveCommand(s, m)
	}

	if strings.HasPrefix(m.Content, "/profile") {
		commandLogger.Info("command /profile triggered")
		go ProfileCommand(s, m)
	}

	if strings.HasPrefix(m.Content, "/h_") {
		commandLogger.Info("command /h_ triggered")
		go HeroCommand(s, m)
	}

	if strings.HasPrefix(m.Content, "/ratingtop") {
		commandLogger.Info("command /ratingtop triggered")
		if strings.HasSuffix(m.Content, "console") {
			go RatingTopCommand(s, m, "console")
		} else {
			go RatingTopCommand(s, m, "pc")
		}
	}
}
