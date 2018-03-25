package actions

import (
	"dwsb/models"
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/pkg/errors"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		discord.New(
			os.Getenv("DISCORD_CLIENT_ID"),
			os.Getenv("DISCORD_CLIENT_SECRET"),
			fmt.Sprintf("%s%s", App().Host, "/auth/discord/callback"),
			"identify",
			"guilds",
		),
	)
}

func LoginRequired(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userid := c.Session().Get("current_user_id"); userid == nil {
			return c.Redirect(302, "/login")
		}
		return next(c)
	}
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}

	// Query Database
	query := models.DB.Where("provider = ? and provider_id = ?", user.Provider, user.UserID)
	exists, err := query.Exists("users")
	if err != nil {
		return errors.WithStack(err)
	}

	// Create record if it doesn't exist
	record := &models.User{}
	if exists {
		err := query.First(record)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		record.Name = user.Name
		record.Provider = user.Provider
		record.ProviderID = user.UserID
		record.Code = c.Param("code")
		record.ExpiresAt = user.ExpiresAt
		record.AccessToken = user.AccessToken
		record.RefreshToken = user.RefreshToken
		err := models.DB.Save(record)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// Set session
	c.Session().Set("current_user_id", record.ID)
	err = c.Session().Save()
	if err != nil {
		return errors.WithStack(err)
	}

	//return c.Render(200, r.JSON(user))
	return c.Redirect(302, "/")
}

func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				fmt.Println(err)
				//	return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

func LogoutHandler(c buffalo.Context) error {
	c.Session().Clear()
	c.Session().Save()
	return c.Redirect(302, "/login")
}
