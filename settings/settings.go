package settings

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
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

	// running buffalo test suite changes the current working directory during test
	// loading, so we need to detect whne we're inside a test run to avoid not
	// finding our settings.yml
	var relativePath string
	if flag.Lookup("test.v") != nil {
		relativePath = "../settings.yml"
	} else {
		relativePath = "settings.yml"
	}
	settings, err := ioutil.ReadFile(relativePath)

	if err != nil {
		wd, _ := os.Getwd()
		panic("Couldn't open settings.yml. Does it exist? " + wd)
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
