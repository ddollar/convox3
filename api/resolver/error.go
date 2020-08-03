package resolver

type AuthenticationError struct {
	error
}

func (a AuthenticationError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code": "UNAUTHENTICATED",
	}
}
