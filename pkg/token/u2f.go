package token

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/settings"
	"github.com/pkg/errors"
	"github.com/tstranex/u2f"
)

type U2FToken struct {
	m model.Interface
}

func NewU2F(m model.Interface) Interface {
	return &U2FToken{m: m}
}

func (t *U2FToken) AuthenticationRequest(uid string) ([]byte, string, error) {
	rs, err := t.registrations(uid)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	ch, chid, err := t.challengeCreate()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	req := ch.SignRequest(rs)

	data, err := json.Marshal(req)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	return data, chid, nil
}

func (t *U2FToken) AuthenticationResponse(uid, chid string, data []byte) error {
	if errr := t.tokenError(data); errr != nil {
		return errr
	}

	var res u2f.SignResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return errors.WithStack(err)
	}

	ch, err := t.challengeGet(chid)
	if err != nil {
		return errors.WithStack(err)
	}

	tks, err := t.m.UserTokens(uid)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, tk := range tks {
		var reg u2f.Registration

		if err := reg.UnmarshalBinary(tk.Data); err != nil {
			return errors.WithStack(err)
		}

		if counter, err := reg.Authenticate(res, *ch, uint32(tk.Counter)); err == nil {
			tk.Counter = int(counter)
			tk.Used = time.Now().UTC()

			if err := t.m.TokenSave(&tk); err != nil {
				return errors.WithStack(err)
			}

			return nil
		}
	}

	return ErrAuthenticate
}

func (t *U2FToken) RegisterRequest(uid string) ([]byte, string, error) {
	rs, err := t.registrations(uid)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	ch, chid, err := t.challengeCreate()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	// req

	// if err := c.SessionSet("token-challenge", string(data)); err != nil {
	//   return err
	// }

	req := u2f.NewWebRegisterRequest(ch, rs)

	data, err := json.Marshal(req)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	return data, chid, nil
}

func (t *U2FToken) RegisterResponse(uid, chid string, data []byte) error {
	if errr := t.tokenError(data); errr != nil {
		return errr
	}

	ch, err := t.challengeGet(chid)
	if err != nil {
		return errors.WithStack(err)
	}

	var res u2f.RegisterResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return errors.WithStack(err)
	}

	reg, err := u2f.Register(res, *ch, &u2f.Config{SkipAttestationVerify: true})
	if err != nil {
		return errors.WithStack(err)
	}

	data, err = reg.MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	token := &model.Token{
		ID:     fmt.Sprintf("%x", reg.KeyHandle),
		UserID: uid,
		Kind:   "u2f",
		Data:   data,
	}

	if err := t.m.TokenSave(token); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *U2FToken) challengeCreate() (*u2f.Challenge, string, error) {
	host := fmt.Sprintf("https://%s", settings.Host)

	ch, err := u2f.NewChallenge(host, []string{host})
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	data, err := json.Marshal(ch)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	mch := &model.Challenge{
		Data: data,
	}

	if err := t.m.ChallengeSave(mch); err != nil {
		return nil, "", errors.WithStack(err)
	}

	return ch, mch.ID, nil
}

func (t *U2FToken) challengeGet(id string) (*u2f.Challenge, error) {
	mch, err := t.m.ChallengeGet(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var ch u2f.Challenge

	if err := json.Unmarshal(mch.Data, &ch); err != nil {
		return nil, errors.WithStack(err)
	}

	return &ch, nil
}

func (t *U2FToken) registrations(uid string) ([]u2f.Registration, error) {
	rs := []u2f.Registration{}

	ts, err := t.m.UserTokensByKind(uid, "u2f")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, t := range ts {
		var r u2f.Registration

		if err := r.UnmarshalBinary(t.Data); err != nil {
			return nil, errors.WithStack(err)
		}

		rs = append(rs, r)
	}

	return rs, nil
}

func (t *U2FToken) tokenError(data []byte) error {
	var te struct {
		ErrorCode int `json:"errorCode"`
	}

	if err := json.Unmarshal(data, &te); err == nil {
		switch te.ErrorCode {
		case 0:
			return nil
		case 1:
			return ErrOther
		case 4:
			return ErrTokenInvalid
		case 5:
			return ErrTimeout
		default:
			return errors.WithStack(fmt.Errorf("token error: %d", te.ErrorCode))
		}
	}

	return nil
}
