package actions

import (
	"dwsb/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func UserHandler(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	userid := c.Session().Get("current_user_id")
	user := models.User{}
	err := tx.Find(&user, userid)
	loadErr := tx.Load(&user, "Clips")

	if err != nil {
		return errors.WithStack(err)
	}
	if loadErr != nil {
		return errors.WithStack(loadErr)
	}
	return c.Render(200, r.JSON(user))
}
