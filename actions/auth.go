package actions

import (
	"dwsb/models"
	"fmt"
	"os"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/pkg/errors"
)

var DiscordProvider goth.Provider

func init() {
	gothic.Store = App().SessionStore

	DiscordProvider = discord.New(
		os.Getenv("DISCORD_CLIENT_ID"),
		os.Getenv("DISCORD_CLIENT_SECRET"),
		fmt.Sprintf("%s%s", App().Host, "/auth/discord/callback"),
		"identify",
		"guilds",
	)

	goth.UseProviders(DiscordProvider)
}

func LoginRequired(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userid := c.Session().Get("current_user_id"); userid == nil {
			return c.Redirect(302, "/login")
		}

		// Although buffalo.Context type has method Value() to pull
		// data that's been stashed into a context, it appears this is
		// only temporary and gets flushed at some point, so we can't
		// use it to validate oAuth expiration and need to rely on a Session
		// variable instead.
		if expiry := c.Session().Get("current_user_expiry"); expiry != nil {
			expTime, _ := time.Parse(time.RFC3339, expiry.(string))
			currentTime := time.Now()
			isAvailable := DiscordProvider.RefreshTokenAvailable()
			fmt.Println("AVAILABLE? ", isAvailable)

			// TODO: this logic to pull a user out of database is happening often,
			// abstract this away
			if currentTime.After(expTime) {
				uid := c.Session().Get("current_user_id")
				u := &models.User{}
				tx := c.Value("tx").(*pop.Connection)
				_ = tx.Find(u, uid)

				// TODO: this should be logged
				newToken, _ := DiscordProvider.RefreshToken(u.RefreshToken)

				u.AccessToken = newToken.AccessToken
				u.RefreshToken = newToken.RefreshToken
				u.ExpiresAt = newToken.Expiry
				_ = models.DB.Save(u)
			}
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
	}

	record.Name = user.Name
	record.Provider = user.Provider
	record.ProviderID = user.UserID
	record.Code = c.Param("code")
	record.ExpiresAt = user.ExpiresAt
	record.AccessToken = user.AccessToken
	record.RefreshToken = user.RefreshToken
	dbErr := models.DB.Save(record)
	if dbErr != nil {
		return errors.WithStack(err)
	}

	// Set session
	expiry := record.ExpiresAt.Format(time.RFC3339)
	c.Session().Set("current_user_id", record.ID)
	c.Session().Set("current_user_expiry", expiry)
	fmt.Println(c.Session().Get("current_user_expiry"))
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
