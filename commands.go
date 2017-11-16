package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func StartCommand(m *discordgo.MessageCreate) {
	dg.ChannelMessageSend(m.ChannelID, "Simple bot for Overwatch by @Lord Protector#9200\n\n"+
		"**How to use:**\n"+
		"1. Use /save to save your game profile.\n"+
		"2. Use /profile to see your stats.\n"+
		"3. ???\n"+
		"4. PROFIT!\n\n"+
		"**Features:**\n"+
		"— Player profile (/profile command)\n"+
		"— Small summary for heroes\n"+
		"— Reports after every game session\n")

	log.Info("/start command executed successful")
}

func DonateCommand(m *discordgo.MessageCreate) {
	dg.ChannelMessageSend(m.ChannelID, "If you find this bot helpful, "+
		"[you can make small donation](https://paypal.me/krasovsky) to help me pay server bills!")

	log.Info("donate command executed successful")
}

type Hero struct {
	Name                string
	TimePlayedInSeconds int
}

type Heroes []Hero

func (hero Heroes) Len() int {
	return len(hero)
}

func (hero Heroes) Less(i, j int) bool {
	return hero[i].TimePlayedInSeconds < hero[j].TimePlayedInSeconds
}

func (hero Heroes) Swap(i, j int) {
	hero[i], hero[j] = hero[j], hero[i]
}

func SaveCommand(m *discordgo.MessageCreate) {
	info := strings.Split(m.Content, " ")
	var text string

	if len(info) == 3 {
		if info[1] != "psn" && info[1] != "xbl" {
			info[2] = strings.Replace(info[2], "#", "-", -1)
		}

		profile, err := GetOverwatchProfile(info[1], info[2])
		if err != nil {
			log.Warn(err)
			text = "Player not found!"
		} else {
			res, err := dg.UserChannelCreate(m.Author.ID)
			if err != nil {
				log.Warn(err)
				return
			}

			_, err = InsertUser(User{
				Id:      fmt.Sprint(dbPKPrefix, m.Author.ID),
				DMId:    res.ID,
				Profile: profile,
				Region:  info[1],
				Nick:    info[2],
			})
			if err != nil {
				log.Warn(err)
				return
			}

			log.Info("/save command executed successful")
			text = "Saved!"
		}
	} else {
		text = "**Example:** `/save eu|us|kr|psn|xbl BattleTag#1337|ConsoleLogin`"
	}

	dg.ChannelMessageSend(m.ChannelID, text)
}

func ProfileCommand(m *discordgo.MessageCreate) {
	user, err := GetUser(fmt.Sprint(dbPKPrefix, m.Author.ID))
	if err != nil {
		log.Warn(err)
		return
	}

	place, err := GetRatingPlace(fmt.Sprint(dbPKPrefix, m.Author.ID))
	if err != nil {
		log.Warn(err)
		return
	}

	log.Info("/profile command executed successful")

	var text string
	info := strings.Split(m.Content, "_")

	if len(info) == 1 {
		text = MakeSummary(user, place, "CompetitiveStats")
	} else if len(info) == 2 && info[1] == "quick" {
		text = MakeSummary(user, place, "QuickPlayStats")
	}

	dg.ChannelMessageSend(m.ChannelID, text)
}

func HeroCommand(m *discordgo.MessageCreate) {
	user, err := GetUser(fmt.Sprint(dbPKPrefix, m.Author.ID))
	if err != nil {
		log.Warn(err)
		return
	}

	log.Info("/h_ command executed successful")

	var text string
	info := strings.Split(m.Content, "_")
	hero := info[1]

	if len(info) == 2 {
		text = MakeHeroSummary(hero, "CompetitiveStats", user)
	} else if len(info) == 3 && info[2] == "quick" {
		text = MakeHeroSummary(hero, "QuickPlayStats", user)
	}

	dg.ChannelMessageSend(m.ChannelID, text)
}

func RatingTopCommand(m *discordgo.MessageCreate, platform string) {
	top, err := GetRatingTop(platform, 20)
	if err != nil {
		log.Warn(err)
		return
	}

	text := "**Rating Top:**\n"
	for i := range top {
		nick := top[i].Nick
		if top[i].Region != "psn" && top[i].Region != "xbl" {
			nick = strings.Replace(nick, "-", "#", -1)
		}
		text += fmt.Sprintf("%d. %s (%d)\n", i+1, nick, top[i].Profile.Rating)
	}

	dg.ChannelMessageSend(m.ChannelID, text)
}
