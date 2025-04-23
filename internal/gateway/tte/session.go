package tte

import (
	"encoding/json"
	"fmt"

	"github.com/dan-frohlich/tabetopevents/internal/logging"
)

type Session struct {
	ID       string `json:"id"`
	UID      string `json:"user_id"`
	userName string
	client   Client
	log      logging.Logger
}

type ApiError struct {
	Message string `json:"message"`
	Data    string `json:"data"`
	Code    int    `json:"code"`
}

func (ae *ApiError) String() string {
	if ae == nil {
		return ""
	}
	return fmt.Sprintf("[%d]: (%s) %s", ae.Code, ae.Data, ae.Message)
}

func (ae *ApiError) Error() string {
	return ae.String()
}

type UserRespose struct {
	Result Users     `json:"result"`
	Err    *ApiError `json:"error"`
}

type Users struct {
	Items  []User  `json:"items"`
	Paging *Paging `json:"paging"`
}

type User map[string]any

func (s Session) TestConnection() (err error) {
	params := map[string]string{
		"session_id": s.ID,
	}
	var b []byte
	b, err = s.client.httpGet("/api/user", params, nil)
	if err != nil {
		return err
	}

	var user UserRespose
	err = json.Unmarshal(b, &user)

	const (
		badSessionCode    = 401
		noSessionCode     = 441
		adminRequiredCode = 450
	)
	if user.Err != nil {
		switch user.Err.Code {
		case badSessionCode, noSessionCode:
			return user.Err
		case adminRequiredCode: //[expetcted] "You must be an admin to do that."
			fallthrough
		default:
		}
	}

	return err

}
