package tte

import (
	"encoding/json"
	"fmt"
	"time"
)

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
	var b []byte
	if b, err = json.Marshal(c); err == nil {
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

type ConventionEventCache struct {
	ConventionEvents []ConventionEvent
	Age              time.Duration
}

type ConventionEventsRespose struct {
	Result ConventionEvents `json:"result"`
	Err    *ApiError        `json:"error"`
}

type ConventionEvents struct {
	Items  []ConventionEvent `json:"items"`
	Paging Paging            `json:"paging"`
}

type ConventionEvent struct {
	EventNumber         int    `json:"event_number"`
	Price               int    `json:"price"`
	StartDate           string `json:"start_date"`
	Private             int    `json:"private"`
	SpaceName           string `json:"space_name"`
	IsScheduled         int    `json:"is_scheduled"`
	UnreservedQuantity  int    `json:"unreserved_quantity"`
	ScheduledDate       string `json:"scheduled_date"`
	StartdaypartID      string `json:"startdaypart_id"`
	DateUpdated         string `json:"date_updated"`
	TypeID              string `json:"type_id"`
	AlternatedaypartID  string `json:"alternatedaypart_id"`
	EndDate             string `json:"end_date"`
	LongDescriptionHTML string `json:"long_description_html"`
	MaxQuantity         int    `json:"max_quantity"`
	AttendeeHeadCount   int    `json:"attendee_head_count"`
	RoomID              string `json:"room_id"`
	PreferreddaypartID  string `json:"preferreddaypart_id"`
	SpaceID             string `json:"space_id"`
	AutoschedulerFailed int    `json:"autoscheduler_failed"`
	ObjectName          string `json:"object_name"`
	WaitCount           int    `json:"wait_count"`
	TakenCount          int    `json:"taken_count"`
	CustomFields        struct {
		Publisher       string `json:"Publisher"`
		Tournament      string `json:"Tournament?"`
		SubCategory     string `json:"SubCategory"`
		HostingGroup    string `json:"HostingGroup"`
		Complexity      string `json:"Complexity"`
		GM              string `json:"GM"`
		TournamentStyle string `json:"TournamentStyle"`
		RulesTaught     string `json:"RulesTaught"`
		Edition         string `json:"Edition"`
		PlayerExp       string `json:"PlayerExp"`
		TournamentStage string `json:"TournamentStage"`
	} `json:"custom_fields"`
	IsTournament           int                          `json:"is_tournament"`
	HostShowedUp           int                          `json:"host_showed_up"`
	StartdaypartName       string                       `json:"startdaypart_name"`
	HostsAlsoPlay          int                          `json:"hosts_also_play"`
	Duration               int                          `json:"duration"`
	HostCount              int                          `json:"host_count"`
	ConventionID           string                       `json:"convention_id"`
	SpacesNeeded           int                          `json:"spaces_needed"`
	Claimable              int                          `json:"claimable"`
	MaxHosts               int                          `json:"max_hosts"`
	IsCancelled            int                          `json:"is_cancelled"`
	SessionCount           int                          `json:"session_count"`
	AgeRange               string                       `json:"age_range"`
	ViewURI                string                       `json:"view_uri"`
	IsOnline               int                          `json:"is_online"`
	ID                     string                       `json:"id"`
	SoldCount              int                          `json:"sold_count"`
	SessionName            interface{}                  `json:"session_name"`
	SpecialRequests        interface{}                  `json:"special_requests"`
	Trashed                int                          `json:"trashed"`
	SessionSeats           int                          `json:"session_seats"`
	ReservedTickets        int                          `json:"reserved_tickets"`
	Relationships          ConventionEventRelationships `json:"_relationships"`
	RoomName               string                       `json:"room_name"`
	SubmissionID           interface{}                  `json:"submission_id"`
	MaxTickets             int                          `json:"max_tickets"`
	AllowScheduleConflicts int                          `json:"allow_schedule_conflicts"`
	AvailableQuantity      int                          `json:"available_quantity"`
	Description            string                       `json:"description"`
	Sellable               int                          `json:"sellable"`
	LongDescription        string                       `json:"long_description"`
	MoreInfoURI            string                       `json:"more_info_uri"`
	Name                   string                       `json:"name"`
	DateCreated            string                       `json:"date_created"`
	ObjectType             string                       `json:"object_type"`
}

type ConventionEventRelationships struct {
	Room             string `json:"room"`
	Type             string `json:"type"`
	Preferreddaypart string `json:"preferreddaypart"`
	Invitees         string `json:"invitees"`
	Eventgroupevents string `json:"eventgroupevents"`
	UnassignSlots    string `json:"unassign_slots"`
	Convention       string `json:"convention"`
	Tickets          string `json:"tickets"`
	Dayparts         string `json:"dayparts"`
	Possiblerooms    string `json:"possiblerooms"`
	Self             string `json:"self"`
	Hosts            string `json:"hosts"`
	Slots            string `json:"slots"`
	Sessions         string `json:"sessions"`
	Rooms            string `json:"rooms"`
	Startdaypart     string `json:"startdaypart"`
	Spaces           string `json:"spaces"`
	Eventhosts       string `json:"eventhosts"`
	OpenReservations string `json:"open_reservations"`
	Mytickets        string `json:"mytickets"`
	Warnings         string `json:"warnings"`
	Wait             string `json:"wait"`
	Alternatedaypart string `json:"alternatedaypart"`
	Cancel           string `json:"cancel"`
	AssignSlots      string `json:"assign_slots"`
	Unschedule       string `json:"unschedule"`
	Space            string `json:"space"`
	Orderitems       string `json:"orderitems"`
	Eventgroups      string `json:"eventgroups"`
	Waits            string `json:"waits"`
	Reservations     string `json:"reservations"`
}
