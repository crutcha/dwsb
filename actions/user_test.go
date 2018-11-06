package actions

import (
	"dwsb/models"
	"time"
)

// TODO: need to test both unauthenticated and authenticated attempts
func (as *ActionSuite) Test_UserInfoHandler() {
	testUser := &models.User{
		Name:         "test user",
		Provider:     "test provider",
		ProviderID:   "111111111",
		AccessToken:  "44ubn1bvs",
		RefreshToken: "adsfj42ybv12",
		DiscordID:    "",
		Code:         "asdfasdfasdfasdf",
		ExpiresAt:    time.Now(),
	}

	testUser2 := &models.User{
		Name:         "test user 2",
		Provider:     "test provider",
		ProviderID:   "111111111",
		AccessToken:  "44ubn1bvssbda",
		RefreshToken: "adsfj42ybv12dbdb",
		DiscordID:    "",
		Code:         "asdfasdfasdfasdf",
		ExpiresAt:    time.Now(),
	}

	t1err := as.DB.Create(testUser)
	t2err := as.DB.Create(testUser2)
	as.NoError(t1err)
	as.NoError(t2err)

	clip1 := &models.Clip{
		Name:   "test clip 1",
		Tag:    "BS Tag",
		File:   "BS FILE",
		Guild:  "BS GUILD",
		UserID: testUser.ID,
	}
	clip2 := &models.Clip{
		Name:   "test clip 2",
		Tag:    "BS Tag 2",
		File:   "BS FILE 2",
		Guild:  "BS GUILD 2",
		UserID: testUser.ID,
	}
	clip3 := &models.Clip{
		Name:   "test clip 3",
		Tag:    "BS Tag 3",
		File:   "BS FILE 3",
		Guild:  "BS GUILD 3",
		UserID: testUser2.ID,
	}

	c1err := as.DB.Create(clip1)
	c2err := as.DB.Create(clip2)
	c3err := as.DB.Create(clip3)
	as.NoError(c1err)
	as.NoError(c2err)
	as.NoError(c3err)

	as.Session.Set("current_user_id", testUser.ID)

	user := &models.User{}
	res := as.JSON("/user/").Get()
	res.Bind(user)

	as.Equal(200, res.Code)
	as.Exactly(len(user.Clips), 2)
	for _, clip := range user.Clips {
		as.Equal(clip.UserID, testUser.ID)
	}
}
