package client

import (
	"dwsb/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/koding/cache"
	"io/ioutil"
	"net/http"
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
func GetGuildsForUser(user models.User) ([]discordgo.UserGuild, error) {
	data, err := guildCache.Get(user.Name)
	if err == nil {
		return data.([]discordgo.UserGuild), nil
	}

	guilds := make([]discordgo.UserGuild, 0)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://discordapp.com/api/v6/users/@me/guilds", nil)
	if err != nil {
		fmt.Println("REG ERR: ", err)
	}
	req.Header.Add("Authorization", "Bearer "+user.AccessToken)
	res, err := client.Do(req)

	// Handle errors or authentication problems
	if err != nil {
		fmt.Println("RESP ERR: ", err)
	}
	if res.StatusCode == 401 {
		// Raise error here to be passed back up through stack, although maybe there's
		// a better way to do this....
		return guilds, errors.New("401 Unauthorized Token")
	}

	body, err := ioutil.ReadAll(res.Body)

	// TODO: auth errors are failing silently but exist within the body of the returned
	// payload. When this is the case JSON will fail to unmarshal, we need to handle this
	// somehow.
	fmt.Println("BODY ", string(body))
	err = json.Unmarshal(body, &guilds)

	guildCache.Set(user.Name, guilds)

	return guilds, nil
}

// Create hashmap of guild choices so we can use CSRF in multiple spots
// TODO: this is used in multiple spots, which could potentially create a lot of API calls.
// Maybe see if this can be cached locally instead of having to save info in database?
func CreateGuildArray(user models.User) ([]map[string]string, error) {
	guilds, guildErr := GetGuildsForUser(user)
	selectArr := make([]map[string]string, 0)

	if guildErr == nil {
		for _, value := range guilds {
			selectMap := make(map[string]string)
			selectMap[value.Name] = value.ID
			selectArr = append(selectArr, selectMap)
		}
	}

	return selectArr, guildErr
}

// Need to also be able to produce guilds as hashmap since templating engine form select
// does not handle arrays well.
func CreateGuildHashmap(user models.User) map[string]string {
	guilds, _ := GetGuildsForUser(user)
	selectMap := make(map[string]string)
	for _, value := range guilds {
		selectMap[value.Name] = value.ID
	}

	return selectMap
}

func EncodeOpus(path string, src *models.Clip) error {

	// Let FFMPEG do all the heavy lifting
	// TODO: use DCA module instead of passing through bash pipe
	inFile := path + "/" + src.Tag
	outFile := path + "/" + src.File
	combined := exec.Command("bash", "-c", "ffmpeg -i "+inFile+" -f s16le -ar 48000 -ac 2 pipe:1 | dca > "+outFile)
	_, err := combined.CombinedOutput()

	return err
}
