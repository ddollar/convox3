package resolver

type TokenAuthenticationRequest struct {
	data string
	id   string
}

func (t *TokenAuthenticationRequest) Data() string {
	return t.data
}

func (t *TokenAuthenticationRequest) Id() string {
	return t.id
}
