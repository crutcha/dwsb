package actions

import (
	"dwsb/models"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// Use a struct for JSON parsing since it will automatically do CSRF for us
type GuildRequest struct {
	Guild string `json:guild`
}

// Board handler is JSON endpoint for refreshing soundboard clips based
// on the guild they are in.
func BoardHandler(c buffalo.Context) error {

	request := &GuildRequest{}
	err := c.Bind(request)
	if err != nil {
		fmt.Println("GUILD BIND: ", err)
	}
	tx := c.Value("tx").(*pop.Connection)

	clips := []models.Clip{}
	query := tx.Where("guild = ?", request.Guild)
	err = query.All(&clips)
	if err != nil {
		fmt.Println("BOARD ERR: ", err)
	}
	return c.Render(200, r.JSON(clips))

}
