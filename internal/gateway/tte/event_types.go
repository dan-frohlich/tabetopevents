package tte

import "encoding/json"

func (s Session) GetConventionEventType(uri string) (cet ConventionEventType, err error) {
	var resp ConventionEventTypeResponse

	params := map[string]string{
		"session_id":             s.ID,
		"_include_relationships": "1",
	}

	var b []byte
	b, err = s.client.httpGet(uri, params, nil)
	if err != nil {
		return cet, err
	}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return cet, err
	}
	s.client.db.Store("event_type", uri, "json", b)

	if resp.Err != nil {
		return cet, resp.Err
	}
	cet = resp.Result
	return cet, nil
}

type ConventionEventTypeResponse struct {
	Result ConventionEventType `json:"result"`
	Err    *ApiError           `json:"error"`
}

type ConventionEventType struct {
	// AllowScheduleConflicts  int64         `json:"allow_schedule_conflicts"`
	// ConventionID            string        `json:"convention_id"`
	CustomFields []CustomField `json:"custom_fields"`
	// DateCreated             string        `json:"date_created"`
	// DateUpdated             string        `json:"date_updated"`
	// DefaultCostPerSlot      int64         `json:"default_cost_per_slot"`
	Description string `json:"description"`
	// EndBuffer               int64         `json:"end_buffer"`
	// GlobalTicketPrice       int64         `json:"global_ticket_price"`
	ID string `json:"id"`
	// LimitTicketAvailability int64         `json:"limit_ticket_availability"`
	// LimitVolunteers         int64         `json:"limit_volunteers"`
	// MaxTickets              int64         `json:"max_tickets"`
	Name string `json:"name"`
	// ObjectName           string        `json:"object_name"`
	// ObjectType           string        `json:"object_type"`
	// OverrideTicketPrices int64         `json:"override_ticket_prices"`
	// SalesStopAfter       int64         `json:"sales_stop_after"`
	// SpecialRequests      []interface{} `json:"special_requests"`
}

type CustomField struct {
	Conditional        int64           `json:"conditional"`
	Edit               int64           `json:"edit"`
	Label              string          `json:"label"`
	Name               string          `json:"name"`
	Required           int64           `json:"required"`
	SequenceNumber     int64           `json:"sequence_number"`
	Type               CustomFieldType `json:"type"`
	View               int64           `json:"view"`
	Options            *string         `json:"options,omitempty"`
	Optionsdescription *string         `json:"optionsdescription,omitempty"`
	ConditionalName    *string         `json:"conditional_name,omitempty"`
	ConditionalValue   *string         `json:"conditional_value,omitempty"`
}

type CustomFieldType string

const (
	Select   CustomFieldType = "select"
	Text     CustomFieldType = "text"
	Textarea CustomFieldType = "textarea"
)
