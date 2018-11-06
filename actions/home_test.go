package actions

// TODO: need to test both unauthenticated and authenticated attempts
func (as *ActionSuite) Test_HomeHandler() {
	res := as.HTML("/").Get()
	location, _ := res.Result().Location()
	as.Equal(302, res.Code)
	as.Equal("/login", location.String())
}
