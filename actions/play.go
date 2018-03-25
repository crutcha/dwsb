package actions

import (
	"dwsb/client"
	"dwsb/models"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// Use a struct for JSON parsing since it will automatically do CSRF for us
type PlayRequest struct {
	Name string `json:name`
}

func PlayHandler(c buffalo.Context) error {
	request := &PlayRequest{}
	err := c.Bind(request)

	if err != nil {
		fmt.Println("PLAY BIND: ", err)
	}

	tx := c.Value("tx").(*pop.Connection)
	clip := models.Clip{}
	err = tx.Find(&clip, request.Name)

	client.PlaySound(clip)
	return err
}
