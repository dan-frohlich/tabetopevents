package tte

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
