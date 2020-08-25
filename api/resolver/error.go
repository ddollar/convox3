package resolver

type AuthenticationError struct {
	error
}

func (e AuthenticationError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code": "UNAUTHENTICATED",
	}
}

type TimeoutError struct {
	error
}

func (e TimeoutError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code": "TIMEOUT",
	}
}

type TokenRequiredError struct {
	error
}

func (e TokenRequiredError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code": "TOKEN_REQUIRED",
	}
}
