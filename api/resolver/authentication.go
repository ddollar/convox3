package resolver

type Authentication struct {
	user User
}

func (a Authentication) Token() (string, error) {
	return a.user.Token()
}

func (a Authentication) User() User {
	return a.user
}
