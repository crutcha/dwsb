package client

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, "GIT GUD")
}

func ChannelStateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, _ := s.State.Channel(m.ChannelID)
	g, _ := s.State.Guild(c.GuildID)

	if m.Content == "/join" {
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				_, err := s.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, false)
				if err != nil {
					fmt.Println(err)
				} else {
					Discord.VoiceReady[c.GuildID] = true
				}
			}
		}
	} else if m.Content == "/leave" {
		if vc, ok := s.VoiceConnections[c.GuildID]; ok {
			vc.Disconnect()
			delete(Discord.VoiceReady, c.GuildID)
		}
	}
}
