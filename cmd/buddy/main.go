package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"

	"github.com/dan-frohlich/tabetopevents/internal/gateway/tte"
	"github.com/dan-frohlich/tabetopevents/internal/logging"
	"github.com/dan-frohlich/tabetopevents/internal/tui"
)

func main() {
	log := logging.Log{Level: logging.LogLevelInfo}
	if term.IsTerminal(0) {
		log.Debug("in a term")
	} else {
		log.Error("not in a term")
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}
	log.Debug("terminal dimaensions", "width", width, "height", height)

	if len(os.Args) > 1 {
		for _, arg := range os.Args {
			switch arg {
			case "-v", "--verbose":
				log.Level = logging.LogLevelDebug
			}
		}
	}
	s, err := extablishSession(log)
	if err != nil {
		log.Fatal("failed to establish tabletop.events session", "error", err)
		return
	}

	con := SelectConvention(log, s)
	fmt.Println(tui.H3.Border(tui.DataBorder, true).Render(
		fmt.Sprintf("%s : %s - %s\n\t%s\n\thttp://tabletop.events%s",
			con.Name, con.StartDate, con.EndDate, con.WebsiteURI, con.ViewURI)))
	// fmt.Println("selected:", con.Name)
	// os.Exit(1)

	var evz []tte.ConventionEvent
	evz, err = getEvents(log, s, con)
	if err != nil {
		log.Fatal("failed to get events", "con", con.ViewURI, "error", err)
	}
	log.Info("found", "event_count", len(evz))

	counts, eventTypeURIByTypeName := getEventTypes(evz, err, s, log)

	log.Info("found", "event_type_count", len(counts))

	displayEventTypeSummary(counts, eventTypeURIByTypeName)

	eventTypeNameByURI := map[string]string{}
	for k, v := range eventTypeURIByTypeName {
		eventTypeNameByURI[v] = k
	}

	filteredEvents := filterEventTypes(log, s, con, evz, eventTypeURIByTypeName)
	log.Info("filtered events", "filtered", len(filteredEvents), "total", len(evz))
	displayEvents(log, width, filteredEvents, eventTypeNameByURI)
}

func displayEvents(log logging.Logger, width int, events []tte.ConventionEvent, eventTypeNameByURI map[string]string) {
	keys := []string{"name", "type", "start", "duration", "description", "publisher", "url"} //, "host"}
	var maxFieldWidth = width - 12 - 4
	for _, ev := range events {
		var out string
		// out += fmt.Sprintf("%7d - %-20s - %s\n", counts[tn], trimLen(tn, 20), eventURIByTypeName[tn])
		m := map[string]string{
			"name":        ev.Name,
			"type":        eventTypeNameByURI[ev.Relationships.Type],
			"start":       ev.StartdaypartName,
			"duration":    fmt.Sprintf("%s", time.Duration(ev.Duration)*time.Minute),
			"description": ev.Description,
			"publisher":   ev.CustomFields.Publisher,
			"url":         "https://tabletop.events" + ev.ViewURI,
			// "game": ev.
			// "host":        ev.Relationships.Eventhosts,
		}
		for _, key := range keys {
			value := m[key]
			if len(value) < maxFieldWidth {
				out += fmt.Sprintf("%-12s: %s\n", key, value)
			} else {
				field := key
				vz := splitString(value, maxFieldWidth)
				for i, v := range vz {
					if i > 0 {
						field = "           ."
					}
					out += fmt.Sprintf("%-12s: %s\n", field, v)
				}
			}
		}
		fmt.Println(lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Italic(true).Render(out[:len(out)-1]))
	}

}

func splitString(s string, partLength int) []string {
	if partLength <= 0 {
		return []string{s}
	}

	var result []string
	for i := 0; i < len(s); i += partLength {
		end := i + partLength
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}

func filterEventTypes(log logging.Logger, s tte.Session, con tte.Convention, events []tte.ConventionEvent, eventTypeURIByTypeName map[string]string) (filteredEvents []tte.ConventionEvent) {
	eventTypeNameByURI := make(map[string]string)
	var eventTypeOpts []huh.Option[string]
	for k, v := range eventTypeURIByTypeName {
		eventTypeNameByURI[v] = k
		eventTypeOpts = append(eventTypeOpts, huh.NewOption(k, v))
	}

	var eventTypes []string
	var title string
	var description string
	huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Event Type(s)").
				Options(eventTypeOpts...).
				Value(&eventTypes),
			huh.NewInput().
				Title("Title Match").
				Value(&title),
			huh.NewInput().
				Title("Description Match").
				Value(&description),
		).Title("Filter Events"),
	).WithTheme(huh.ThemeBase16()).
		WithShowHelp(true).
		WithShowErrors(true).
		Run()
	var pred []tte.EventPredicate
	if len(title) > 0 {
		pred = append(pred, func(ce tte.ConventionEvent) bool {
			return strings.Contains(strings.ToLower(ce.Name), strings.ToLower(title))
		})
	}
	if len(description) > 0 {
		pred = append(pred, func(ce tte.ConventionEvent) bool {
			lmatch := strings.ToLower(description)
			ld := strings.ToLower(ce.Description)
			lld := strings.ToLower(ce.LongDescription)
			return strings.Contains(ld, lmatch) || strings.Contains(lld, lmatch)
		})
	}
	if len(eventTypes) > 0 {
		pred = append(pred, func(ce tte.ConventionEvent) bool {
			for _, et := range eventTypes {
				if ce.Relationships.Type == et {
					return true
				}
			}
			return false
		})
	}

	return tte.FilterableConventionEvents(events).Filter(pred)
}

func displayEventTypeSummary(counts map[string]int, eventURIByTypeName map[string]string) {
	var typeNames = make([]string, 0, len(counts))
	for k := range counts {
		typeNames = append(typeNames, k)
	}
	sort.Strings(typeNames)
	var out string
	out = "  COUNT - EVENT TYPE           - URI\n"
	for _, tn := range typeNames {
		out += fmt.Sprintf("%7d - %-20s - %s\n", counts[tn], trimLen(tn, 20), eventURIByTypeName[tn])
	}
	fmt.Println(lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Italic(true).Render(out[:len(out)-1]))
}

func getEventTypes(evz []tte.ConventionEvent, err error, s tte.Session, log logging.Log) (map[string]int, map[string]string) {
	counts := make(map[string]int)
	// eventTypeByID := make(map[string]tte.ConventionEventType)
	// eventTypeByName := make(map[string]tte.ConventionEventType)
	eventTypeByURI := make(map[string]tte.ConventionEventType)
	eventURIByTypeName := make(map[string]string)

	var (
		cet tte.ConventionEventType
		ok  bool
	)
	for _, ev := range evz {
		if cet, ok = eventTypeByURI[ev.Relationships.Type]; !ok {
			cet, err = s.GetConventionEventType(ev.Relationships.Type)
			if err != nil {
				log.Error("failed to get event type from event", "event_type_uri", ev.Relationships.Type, "event_number", ev.EventNumber, "error", err)
			}
			eventTypeByURI[ev.Relationships.Type] = cet
			eventURIByTypeName[cet.Name] = ev.Relationships.Type
		}
		counts[cet.Name] += 1
	}
	return counts, eventURIByTypeName
}

func trimLen(s string, maxLex int) string {
	if len(s) < maxLex {
		return s
	}
	return s[:maxLex]
}

func getEvents(log logging.Logger, s tte.Session, con tte.Convention) (events []tte.ConventionEvent, err error) {
	var ignoreCachedEventInfo bool
	ignoreCachedEventInfo = true
	cache, err := s.GetCachedConventionEvents(con)
	events = cache.ConventionEvents
	if err == nil {
		ignoreCachedEventInfo = false
		huh.NewConfirm().
			Title(fmt.Sprintf("cached convention event data was found [%s old], shall we use it?", cache.Age)).
			Affirmative("No.").
			Negative("Yes!").
			Value(&ignoreCachedEventInfo).
			WithTheme(huh.ThemeBase16()).
			Run()
	} else {
		log.Error("GetCachedConventionEvents", "error", err)
	}
	log.Info("use cache?", "ignoreCachedEventInfo", ignoreCachedEventInfo)
	if ignoreCachedEventInfo {
		events, err = s.GetConventionEvents(con)

	}
	return events, err
}

func extablishSession(log logging.Log) (tte.Session, error) {
	var useCachedApiKey bool = true
	c, err := tte.RestoreClient(log)
	useCachedApiKey = err == nil
	if err != nil || !useCachedApiKey {
		var apiKey string
		huh.NewInput().
			Title("establishing tabletop.events session").
			Prompt("input tabletop.events api key:").
			Value(&apiKey).
			WithTheme(huh.ThemeBase16()).
			Run()
		c = tte.NewClient(log, apiKey)
	}

	var s tte.Session
	s, err = c.RestoreSession(log)
	if err != nil { //|| !useCachedSessionId {
		log.Info("creating a new session")

		var username string
		var password string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("tabletop.evetns login"),
			),
			huh.NewGroup(
				huh.NewInput().
					// Title("tabletop.evetns login").
					Prompt("username:").
					Value(&username),
				huh.NewInput().EchoMode(huh.EchoModePassword).
					Prompt("password:").
					Value(&password),
			),
			// WithTheme(huh.ThemeBase16()).
			// Title("tabletop.evetns login").Description("the login page").WithShowErrors(true),
		).
			WithLayout(huh.LayoutStack).
			WithShowErrors(true).
			WithTheme(huh.ThemeBase16())
		// form.Update(form.Init())
		form.NextGroup()
		form.Run()

		if len(username) == 0 || len(password) == 0 {
			log.Fatal("username and password must be provided")
			os.Exit(1)
		}
		s, err = c.NewSession(log, username, password)
	}
	return s, err
}

func SelectConvention(log logging.Logger, s tte.Session) tte.Convention {

	var ignoreCachedConventionInfo bool
	ignoreCachedConventionInfo = true
	cache, err := s.GetCachedActiveConventions()
	if err == nil {
		ignoreCachedConventionInfo = false
		huh.NewConfirm().
			Title(fmt.Sprintf("cached convention data was found [%s old], shall we use it?", cache.Age)).
			Affirmative("No.").
			Negative("Yes!").
			Value(&ignoreCachedConventionInfo).
			WithTheme(huh.ThemeBase16()).
			Run()
	}
	cz := cache.Conventions
	if ignoreCachedConventionInfo {
		cz, err = s.GetActiveConventions()
	}
	var conmap = make(map[string]tte.Convention)

	for _, con := range cz {
		conmap[con.Name] = con
	}
	var conNames = make([]string, 0, len(cz))
	for k := range conmap {
		conNames = append(conNames, k)
	}
	sort.Strings(conNames)
	// fmt.Println(strings.Join(conNames, "\n"))

	var conName string
	field := huh.NewSelect[string]().
		Height(21).
		Title("Pick a convention.").
		Options(huh.NewOptions(conNames...)...).
		Value(&conName).
		WithTheme(huh.ThemeBase16())
	huh.NewForm(huh.NewGroup(field)).WithShowHelp(true).Run()

	return conmap[conName]
}
