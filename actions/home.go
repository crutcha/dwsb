package actions

import (
	"dwsb/client"
	"dwsb/models"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	// Allow for query param for default selected guild so we can use it
	// on callback redirects. This also gives us quotes which we need to remove.
	defaultSelect := c.Param("default")
	if defaultSelect != "" {
		if defaultSelect[0] == '"' {
			defaultSelect = defaultSelect[1:]
		}
		if i := len(defaultSelect) - 1; defaultSelect[i] == '"' {
			defaultSelect = defaultSelect[:1]
		}
	}

	tx := c.Value("tx").(*pop.Connection)
	userid := c.Session().Get("current_user_id")
	user := models.User{}
	err := tx.Find(&user, userid)
	if err != nil {
		fmt.Println("USER ERROR: ", err)
	}

	selectArr, err := client.CreateGuildArray(user)

	if err != nil {
		// Force re-login to deal with token refresh
		c.Session().Clear()
		return c.Redirect(302, "/login")
	} else {
		c.Set("selectmap", selectArr)
		c.Set("defaultselect", defaultSelect)
		return c.Render(200, r.HTML("board.html"))
	}

}
