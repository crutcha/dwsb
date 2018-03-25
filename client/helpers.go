package client

import (
	"dwsb/models"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/koding/cache"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// 5 minute cache guild API calls
var guildCache *cache.MemoryTTL

func init() {
	guildCache = cache.NewMemoryWithTTL(300 * time.Second)
	guildCache.StartGC(1 * time.Second)
}

// Discordgo uses bot token instead of bearer and also use websocket API. Since this is the only
// RESTful call, simple use built-in HTTP library to do GET and grab guild JSON then map it to
// guild struct.
func GetGuildsForUser(user models.User) []discordgo.UserGuild {
	data, err := guildCache.Get(user.Name)
	if err == nil {
		return data.([]discordgo.UserGuild)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://discordapp.com/api/v6/users/@me/guilds", nil)
	if err != nil {
		fmt.Println("REG ERR: ", err)
	}
	req.Header.Add("Authorization", "Bearer "+user.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("RESP ERR: ", err)
	}

	guilds := make([]discordgo.UserGuild, 0)
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &guilds)
	guildCache.Set(user.Name, guilds)

	return guilds
}

// Create hashmap of guild choices so we can use CSRF in multiple spots
// TODO: this is used in multiple spots, which could potentially create a lot of API calls.
// Maybe see if this can be cached locally instead of having to save info in database?
func CreateGuildArray(user models.User) []map[string]string {
	guilds := GetGuildsForUser(user)
	selectArr := make([]map[string]string, 0)
	for _, value := range guilds {
		selectMap := make(map[string]string)
		selectMap[value.Name] = value.ID
		selectArr = append(selectArr, selectMap)
	}

	return selectArr
}

// Need to also be able to produce guilds as hashmap since templating engine form select
// does not handle arrays well.
func CreateGuildHashmap(user models.User) map[string]string {
	guilds := GetGuildsForUser(user)
	selectMap := make(map[string]string)
	for _, value := range guilds {
		selectMap[value.Name] = value.ID
	}

	return selectMap
}

func EncodeOpus(src *os.File) error {

	// Let FFMPEG do all the heavy lifting
	// TODO: use DCA module instead of passing through bash pipe
	name := src.Name() + ".dca"
	combined := exec.Command("bash", "-c", "ffmpeg -i "+src.Name()+" -f s16le -ar 48000 -ac 2 pipe:1 | dca > "+name)
	fmt.Println(combined)
	_, err := combined.CombinedOutput()

	return err
}
