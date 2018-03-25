package client

import (
	"dwsb/models"
	"dwsb/settings"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

var (
	speakers    map[uint32]*gopus.Decoder
	opusEncoder *gopus.Encoder
	Discord     *DiscordClient
)

// discordgo.VoiceConnection contains a RWMutex as well as Ready bool property but both are used internally
// to track opussend/opusrecv channels. Since the desired state is to simply drop play requests if another clip is
// already playing instead of waiting for a lock and playing after, discordgo.Session must be extended
// to allow for this.
type DiscordClient struct {
	Session    *discordgo.Session
	VoiceReady map[string]bool // Whether or not voice is being sent, guild as key
}

//Create exportable client session
func init() {
	session, err := discordgo.New(os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("DISCORD ERR: ", err)
	}
	Discord = &DiscordClient{Session: session, VoiceReady: make(map[string]bool)}
}

// TODO: get rid of this trash
func OnError(str string, err error) {
	prefix := "dgVoice: " + str

	if err != nil {
		os.Stderr.WriteString(prefix + ": " + err.Error())

	} else {
		os.Stderr.WriteString(prefix)

	}

}

func (c *DiscordClient) Connect() {
	c.Session.AddHandler(ReadyHandler)
	c.Session.AddHandler(ChannelStateHandler)
	err := Discord.Session.Open()
	if err != nil {
		panic(err)
	}

	// Anonymous goroutine to catch sigterm and close cleanly
	s := make(chan os.Signal, 2)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-s
		fmt.Println("Closing Discord Websocket...")
		err := c.Session.Close()
		if err != nil {
			fmt.Println("CLOSE ERROR: ", err)
		}
	}()
}

func changeVoiceReady(guild string) {
	Discord.VoiceReady[guild] = true
}

func PlaySound(file models.Clip) {
	if vc, ok := Discord.Session.VoiceConnections[file.Guild]; ok {
		// Create file based on path in settings
		fileDir := filepath.Join(settings.LocalSettings.Media.Location, file.Guild)
		filePath := fileDir + "/" + file.File
		if Discord.VoiceReady[file.Guild] {
			Discord.VoiceReady[file.Guild] = false
			defer changeVoiceReady(file.Guild)

			//playAudioFile(vc, filePath, make(chan bool))
			loadSoundBuffer(vc, filePath)
		} else {
			fmt.Println("Guild Voice not ready")
		}
	} else {
		fmt.Println("not playing cause not in channel")
	}

}

func loadSoundBuffer(vc *discordgo.VoiceConnection, filename string) error {
	file, err := os.Open(filename)
	var opuslen uint16
	buffer := make([][]byte, 0)

	for err == nil {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, send buffer to opus receiver
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err

			}
			sendToOpus(vc, buffer)

		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err

		}

		// Append encoded pcm data to the buffer.
		// TODO: make this into a channel that can be consumed by sendToOpus instead of
		// holding the entire file in memory
		buffer = append(buffer, InBuf)

	}

	return err
}

func sendToOpus(vc *discordgo.VoiceConnection, buffer [][]byte) {
	// Send "speaking" packet over the voice websocket
	err := vc.Speaking(true)
	if err != nil {
		OnError("Couldn't set speaking", err)

	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := vc.Speaking(false)
		if err != nil {
			OnError("Couldn't stop speaking", err)

		}

	}()

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}
}
