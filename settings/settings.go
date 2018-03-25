package settings

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var LocalSettings Settings

type Settings struct {
	Media struct {
		Location string `yaml:location`
		MaxSize  int    `yaml:maxsize`
	}
}

func init() {

	// Load user defined settings
	LocalSettings = Settings{}
	settings, err := ioutil.ReadFile("settings.yml")
	if err != nil {
		panic("Couldn't open settings.yml. Does it exist?")
	}
	err = yaml.Unmarshal(settings, &LocalSettings)
	if err != nil {
		panic("Couldn't parse settings.yml for config. Please check syntax.")
	}

	fmt.Println("SETTINGS LOADED:", LocalSettings)
	// Correct for missing trailing slash if applicable
	/*
		pathlen := len(LocalSettings.Media.Location)
		if LocalSettings.Media.Location[pathlen-1] != '/' {
			LocalSettings.Media.Location += "/"
		}
	*/
}
