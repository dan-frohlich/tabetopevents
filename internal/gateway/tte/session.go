package tte

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dan-frohlich/tabetopevents/internal/logging"
)

type Session struct {
	ID       string `json:"id"`
	UID      string `json:"user_id"`
	userName string
	client   Client
	log      logging.Log
}

type ConventionRespose struct {
	Result Conventions `json:"result"`
	Err    *ApiError   `json:"error"`
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

type Conventions struct {
	Items  []Convention `json:"items"`
	Paging *Paging      `json:"paging"`
}

func (s Session) GetCachedActiveConventions() (cache ConventionCache, err error) {
	var b []byte
	b, err = s.client.db.Read("conventions", "conventions", "json")
	if err != nil {
		s.log.Error("read", "error", err)
		return cache, err
	}
	s.log.Debug("read", "bytes", len(b))
	c := &Conventions{}
	err = json.Unmarshal(b, c)
	if err != nil {
		return cache, err
	}
	if len(c.Items) == 0 {
		return cache, fmt.Errorf("cache contained zero items")
	}
	cache.Conventions = c.Items
	cache.Age, err = s.client.db.CacheAge("conventions", "conventions", "json")

	return cache, err
}

func (s Session) GetActiveConventions() (cz []Convention, err error) {
	cr, err := s.getConventionsByPage(1)
	if err != nil {
		return cz, err
	}
	if cr.Err != nil {
		return cz, cr.Err
	}
	cz = append(cz, cr.Result.Items...)
	if cr.Result.Paging == nil {
		return cz, err
	}

	nextPage := cr.Result.Paging.NextPageNumber
	pageCount := cr.Result.Paging.TotalPages
	if pageCount == 0 {
		return cz, err
	}
	for i := nextPage; i <= pageCount; i++ {
		cr, _ = s.getConventionsByPage(int(i))
		if cr.Err != nil {
			return cz, fmt.Errorf("[%d]: (%s) %s", cr.Err.Code, cr.Err.Data, cr.Err.Message)
		}
		cz = append(cz, cr.Result.Items...)
	}
	c := &Conventions{Items: cz}
	if b, e := json.Marshal(c); e == nil {
		s.client.db.Store("conventions", "conventions", "json", b)
	}
	return cz, err
}

func (s Session) getConventionsByPage(page int) (cr ConventionRespose, err error) {
	params := map[string]string{
		"session_id":      s.ID,
		"_page_number":    fmt.Sprintf("%d", page),
		"_items_per_page": "100",
	}
	var b []byte
	b, err = s.client.httpGet("/api/convention", params, nil)
	if err != nil {
		return cr, err
	}

	err = json.Unmarshal(b, &cr)
	if err != nil {
		return cr, err
	}
	if cr.Err != nil {
		return cr, cr.Err
	}
	return cr, err
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

type ConventionCache struct {
	Conventions []Convention
	Age         time.Duration
}

type ConventionEventCache struct {
	ConventionEvents []ConventionEvent
	Age              time.Duration
}

func (s Session) GetCachedConventionEvents(con Convention) (cache ConventionEventCache, err error) {
	c := &ConventionEvents{}
	var b []byte
	b, err = s.client.db.Read("events", con.ViewURI, "json")
	if err != nil {
		return cache, err
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		return cache, err
	}
	cache.ConventionEvents = c.Items
	cache.Age, err = s.client.db.CacheAge("events", con.ViewURI, "json")
	return cache, err
}

func (s Session) GetConventionEvents(con Convention) (ez []ConventionEvent, err error) {
	conID := con.ID
	var resp ConventionEventsRespose
	resp, err = s.getConventionEventsByPage(conID, 1)
	if err != nil {
		return ez, err
	}
	ez = append(ez, resp.Result.Items...)

	nextPage := resp.Result.Paging.NextPageNumber
	pageCount := resp.Result.Paging.TotalPages
	if pageCount == 0 {
		return ez, err
	}
	for i := nextPage; i <= pageCount; i++ {
		s.log.Debug("getting pages", "current", i, "last", resp.Result.Paging.TotalPages)
		resp, _ = s.getConventionEventsByPage(conID, int(i))
		if resp.Err != nil {
			return ez, fmt.Errorf("[%d]: (%s) %s", resp.Err.Code, resp.Err.Data, resp.Err.Message)
		}
		ez = append(ez, resp.Result.Items...)
	}
	c := &ConventionEvents{Items: ez}
	if b, e := json.Marshal(c); e == nil {
		s.client.db.Store("events", con.ViewURI, "json", b)
	}
	return ez, err

}

func (s Session) getConventionEventsByPage(conID string, page int) (cr ConventionEventsRespose, err error) {
	params := map[string]string{
		"session_id":             s.ID,
		"_page_number":           fmt.Sprintf("%d", page),
		"_items_per_page":        "100",
		"_include_relationships": "1",
	}
	uri := fmt.Sprintf("/api/convention/%s/events", conID)
	var b []byte
	b, err = s.client.httpGet(uri, params, nil)

	cer := ConventionEventsRespose{}
	err = json.Unmarshal(b, &cer)
	return cer, err
}
