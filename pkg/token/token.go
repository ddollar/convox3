package token

import "fmt"

type Interface interface {
	AuthenticationRequest(uid string) ([]byte, string, error)
	AuthenticationResponse(uid, chid string, req []byte) error
	RegisterRequest(uid string) ([]byte, string, error)
	RegisterResponse(uid, chid string, req []byte) error
}

var (
	ErrAuthenticate = fmt.Errorf("invalid token authentication")
	ErrTokenInvalid = fmt.Errorf("token invalid")
	ErrOther        = fmt.Errorf("unknown error")
	ErrTimeout      = fmt.Errorf("timeout")
)
