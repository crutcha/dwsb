package actions

import (
	"dwsb/client"
	"dwsb/models"
	"dwsb/settings"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func UploadHandler(c buffalo.Context) error {

	// Grab guilds for template
	// TODO: error handling here
	tx := c.Value("tx").(*pop.Connection)
	uid := c.Session().Get("current_user_id")
	user := models.User{}
	err := tx.Find(&user, uid)
	if err != nil {
		return errors.WithStack(err)
	}

	clip := &models.Clip{}
	c.Set("selectmap", client.CreateGuildHashmap(user))
	c.Set("clip", clip)
	fmt.Println(settings.LocalSettings.Media.MaxSize)
	c.Set("max_size", settings.LocalSettings.Media.MaxSize)
	return c.Render(200, r.HTML("upload.html"))
}

func UploadCreate(c buffalo.Context) error {
	clip := &models.Clip{}

	// Populate fields not in form that are needed for model
	upload, _ := c.File("file")
	clip.File = upload.Filename
	clip.Tag = upload.Filename
	if err := c.Bind(clip); err != nil {
		return errors.WithStack(err)
	}

	guildpath := filepath.Join(settings.LocalSettings.Media.Location, clip.Guild)
	if _, err := os.Stat(guildpath); os.IsNotExist(err) {
		os.Mkdir(guildpath, os.ModePerm)
	}
	f, err := os.Create(guildpath + "/" + clip.File)

	if err != nil {
		return errors.WithStack(err)
	}

	_, err = io.Copy(f, upload.File)

	//Convert MP3 to Opus
	opusErr := client.EncodeOpus(f)

	// If encoding succeeded, save to database, else alert
	if opusErr == nil {
		tx := c.Value("tx").(*pop.Connection)

		// Update with DCA file after MP3 conversion
		clip.File = upload.Filename + ".dca"
		err := tx.Create(clip)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		return errors.WithStack(err)
	}

	return c.Redirect(302, "/?default="+clip.Guild)
}
