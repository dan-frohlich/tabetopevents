package main

import (
	"fmt"
	"os"
	"os/exec"
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
	a := &app{log: log, db: tte.NewDB(log)}

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
	err = a.extablishSession()
	if err != nil {
		log.Fatal("failed to establish tabletop.events session", "error", err)
		return
	}

	con := a.SelectConvention()
	println(tui.H3.Border(tui.DataBorder, true).Render(
		fmt.Sprintf("%s : %s - %s\n\t%s\n\thttp://tabletop.events%s",
			con.Name, con.StartDate, con.EndDate, con.WebsiteURI, con.ViewURI)))
	a.println("selected:", con.Name)
	// os.Exit(1)

	var evz []tte.ConventionEvent
	evz, err = a.getEvents()
	if err != nil {
		log.Fatal("failed to get events", "con", con.ViewURI, "error", err)
	}
	log.Info("found", "event_count", len(evz))

	counts, eventTypeURIByTypeName := a.getEventTypes(evz)

	log.Info("found", "event_type_count", len(counts))

	a.displayEventTypeSummary(counts, eventTypeURIByTypeName)

	eventTypeNameByURI := map[string]string{}
	for k, v := range eventTypeURIByTypeName {
		eventTypeNameByURI[v] = k
	}

	a.readLikesFromCache()
	var allLiked = a.likes

	var stop bool
	for !stop {

		filteredEvents := a.filterEventTypes(evz, eventTypeURIByTypeName)
		log.Info("filtered events", "filtered", len(filteredEvents), "total", len(evz))

		sort.Slice(filteredEvents, func(i, j int) bool {
			if filteredEvents[i].StartdaypartName != filteredEvents[j].StartdaypartName {
				return filteredEvents[i].StartdaypartName.Compare(filteredEvents[j].StartdaypartName)
			}
			return filteredEvents[i].Name < filteredEvents[j].Name
		})

		a.displayEvents(log, width, filteredEvents, eventTypeNameByURI)

		var like = false

		if len(filteredEvents) > 0 {

			huh.NewConfirm().
				Title("like some of these?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&like).
				WithTheme(huh.ThemeBase16()).
				Run()

			if like {
				//TODO diplay a ui to multiselect these
				var options []huh.Option[string]
				for _, fe := range filteredEvents {
					if a.isLiked(fe) {
						continue
					}
					value := fe.ViewURI
					eventType := eventTypeNameByURI[fe.Relationships.Type]
					key := fmt.Sprintf("%4d - [%-16s] %s [%s] (%s)", fe.EventNumber, eventType, tui.H4.Render(fe.Name), string(fe.StartdaypartName), fmt.Sprintf("%s", time.Duration(fe.Duration)*time.Minute))
					options = append(options, huh.NewOption(key, value))
				}
				if len(options) > 0 {
					var liked []string
					huh.NewMultiSelect[string]().
						Title("Like Events").
						Options(options...).
						Value(&liked).
						WithTheme(huh.ThemeBase16()).
						Run()
					allLiked = append(allLiked, liked...)
				}
			}

			var unlike = false

			huh.NewConfirm().
				Title("un-like some of these?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&unlike).
				WithTheme(huh.ThemeBase16()).
				Run()

			if unlike {
				//TODO diplay a ui to multiselect these
				var options []huh.Option[string]
				for _, fe := range filteredEvents {
					if !a.isLiked(fe) {
						continue
					}
					value := fe.ViewURI
					eventType := eventTypeNameByURI[fe.Relationships.Type]
					key := fmt.Sprintf("%4d - [%-16s] %s [%s] (%s)", fe.EventNumber, eventType, tui.H4.Render(fe.Name), string(fe.StartdaypartName), fmt.Sprintf("%s", time.Duration(fe.Duration)*time.Minute))
					options = append(options, huh.NewOption(key, value))
				}
				if len(options) > 0 {
					var unliked []string
					huh.NewMultiSelect[string]().
						Title("Un-Like Events").
						Options(options...).
						Value(&unliked).
						WithTheme(huh.ThemeBase16()).
						Run()
					//TODO remove unlined from allLiked
					var result []string
					var set = make(map[string]struct{})
					for _, ul := range unliked {
						set[ul] = struct{}{}
					}
					for _, l := range allLiked {
						if _, ok := set[l]; !ok {
							result = append(result, l)
						}
					}
					allLiked = result

				}
			}
		}
		sort.Strings(allLiked)
		a.likes = allLiked
		a.db.Store("liked", con.ViewURI, "txt", []byte(strings.Join(allLiked, "\n")))
		huh.NewConfirm().
			Title("again?").
			Affirmative("No.").
			Negative("Yes!").
			Value(&stop).
			WithTheme(huh.ThemeBase16()).
			Run()
	}
	var open bool
	huh.NewConfirm().
		Title("open liked events?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&open).
		WithTheme(huh.ThemeBase16()).
		Run()

	sort.Strings(allLiked)
	for _, like := range allLiked {
		if open {
			var url string
			url = fmt.Sprintf("https://tabletop.events%s", like)

			log.Info("opening", "url", url)
			cmd := exec.Command("open", url)
			_, _ = cmd.Output()
		} else {
			log.Info("liked", "uri", like)
		}
	}
	a.db.Store("liked", con.ViewURI, "txt", []byte(strings.Join(allLiked, "\n")))
}

func (a *app) readLikesFromCache() {
	var b []byte
	b, _ = a.db.Read("liked", a.con.ViewURI, "txt")
	a.likes = strings.Split(string(b), "\n")
}

type app struct {
	con   tte.Convention
	db    tte.DB
	likes []string
	log   logging.Logger
	s     tte.Session
}

func (a *app) isLiked(ce tte.ConventionEvent) bool {
	for _, l := range a.likes {
		if l == ce.ViewURI {
			return true
		}
	}
	return false
}
func (a *app) displayEvents(log logging.Logger, width int, events []tte.ConventionEvent, eventTypeNameByURI map[string]string) {
	keys := []string{"name", "number", "type", "start", "duration", "description", "publisher", "host group", "game master", "url"} //, "host"}
	const padding = 8
	var maxFieldWidth = 80 // width - 12 - padding - 18
	likeMap := make(map[string]struct{})
	for _, l := range a.likes {
		likeMap[l] = struct{}{}
	}
	for _, ev := range events {
		var out string
		// out += fmt.Sprintf("%7d - %-20s - %s\n", counts[tn], trimLen(tn, 20), eventURIByTypeName[tn])

		name := ev.Name
		if _, ok := likeMap[ev.ViewURI]; ok {
			liked := logging.WarnStyle.Style.Bold(true).Render("(*)")
			name = fmt.Sprintf("%s%s", liked, ev.Name)
		}
		m := map[string]string{
			"name":        name,
			"number":      fmt.Sprintf("%d", ev.EventNumber),
			"type":        eventTypeNameByURI[ev.Relationships.Type],
			"start":       string(ev.StartdaypartName),
			"duration":    fmt.Sprintf("%s", time.Duration(ev.Duration)*time.Minute),
			"description": strip(ev.Description, "\n"),
			"publisher":   ev.CustomFields.Publisher,
			"host group":  ev.CustomFields.HostingGroup,
			"game master": ev.CustomFields.GM,
			"url":         "https://tabletop.events" + ev.ViewURI,
			// "game": ev.
			// "host":        ev.Relationships.Eventhosts,
		}
		for _, key := range keys {
			value := m[key]
			if len(value) < maxFieldWidth {
				out += fmt.Sprintf("%12s: %s\n", key, value)
			} else {
				field := key
				vz := splitString(value, maxFieldWidth)
				for i, v := range vz {
					if i > 0 {
						field = "           "
						out += fmt.Sprintf("%12s  %s\n", field, v)
					} else {
						out += fmt.Sprintf("%12s: %s\n", field, v)
					}
				}
			}
		}
		println(lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			Italic(true).
			MaxWidth(width - padding).
			Render(out[:len(out)-1]))
	}
	println("\n")
}

func (a *app) println(args ...any) {
	fmt.Println(args...)
	// filePath := "./log.txt"
	// if log != nil {
	// 	log.Debug("writing log", "path", filePath)
	// }
	// data := strings.Join(argsToStrings(args), " ")
	// _ = os.WriteFile(filePath, []byte(data), os.FileMode(0644))
}

func argsToStrings(args []any) (sz []string) {
	sz = make([]string, 0, len(args))
	for _, arg := range args {
		sz = append(sz, argToString(arg))
	}
	return sz
}

func argToString(arg any) string {
	switch v := arg.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	case int, int16, int32, int64, int8, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func strip(s string, cutset string) string {
	var out []rune
	for _, r := range s {
		if strings.Contains(cutset, string(r)) {
			continue
		}
		out = append(out, r)
	}
	return string(out)
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

func (a *app) filterEventTypes(events []tte.ConventionEvent, eventTypeURIByTypeName map[string]string) (filteredEvents []tte.ConventionEvent) {

	eventTypeNameByURI := make(map[string]string)
	var eventTypeOpts []huh.Option[string]
	for k, v := range eventTypeURIByTypeName {
		eventTypeNameByURI[v] = k
		eventTypeOpts = append(eventTypeOpts, huh.NewOption(k, v))
	}

	sort.Slice(eventTypeOpts, func(i, j int) bool {
		return eventTypeOpts[i].Key < eventTypeOpts[j].Key
	})

	var (
		eventTypes  []string
		description string
		host        string
		id          string
		title       string
		areLiked    string
	)

	huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("Event Type(s)").
				Options(eventTypeOpts...).Value(&eventTypes),
			huh.NewInput().Title("ID").Value(&id),
			huh.NewInput().Title("Title Match").Value(&title),
			huh.NewSelect[string]().Title("Liked Events").
				Options(huh.NewOption[string]("either", "either"), huh.NewOption[string]("liked", "liked"), huh.NewOption[string]("not liked", "not liked")).
				Value(&areLiked),
			huh.NewInput().Title("Description Match").Value(&description),
			huh.NewInput().Title("Hosting Group").Value(&host),
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
	if len(id) > 0 {
		pred = append(pred, func(ce tte.ConventionEvent) bool {
			return strings.ToLower(ce.ID) == strings.ToLower(id)
		})
	}
	if len(areLiked) > 0 && areLiked != "either" {
		pred = append(pred, func(ce tte.ConventionEvent) bool {
			if areLiked == "liked" {
				for _, l := range a.likes {
					if l == ce.ViewURI {
						return true
					}
				}
				return false
			}
			for _, l := range a.likes {
				if l == ce.ViewURI {
					return false
				}
			}
			return true
		})
	}
	if len(host) > 0 {
		pred = append(pred, func(ce tte.ConventionEvent) bool {
			return strings.Contains(strings.ToLower(ce.CustomFields.HostingGroup), strings.ToLower(host))
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

func (a *app) displayEventTypeSummary(counts map[string]int, eventURIByTypeName map[string]string) {
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
	println(lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Italic(true).Render(out[:len(out)-1]))
}

func (a *app) getEventTypes(evz []tte.ConventionEvent) (map[string]int, map[string]string) {

	counts := make(map[string]int)
	eventTypeByURI := make(map[string]tte.ConventionEventType)
	eventURIByTypeName := make(map[string]string)

	var (
		err error
		cet tte.ConventionEventType
		log logging.Logger = a.log
		ok  bool
		s   tte.Session = a.s
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

func (a *app) getEvents() (events []tte.ConventionEvent, err error) {
	var (
		log logging.Logger = a.log
		s   tte.Session    = a.s
		con tte.Convention = a.con
	)
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

func (a *app) extablishSession() error {
	var log logging.Logger = a.log
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
	s, err = c.RestoreSession()
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
		s, err = c.NewSession(username, password)
	}
	a.s = s
	return err
}

func (a *app) SelectConvention() tte.Convention {

	var (
		s tte.Session = a.s
	)
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
	a.println(strings.Join(conNames, "\n"))

	var conName string
	field := huh.NewSelect[string]().
		Height(21).
		Title("Pick a convention.").
		Options(huh.NewOptions(conNames...)...).
		Value(&conName).
		WithTheme(huh.ThemeBase16())
	huh.NewForm(huh.NewGroup(field)).WithShowHelp(true).Run()

	con := conmap[conName]
	a.con = con
	return con
}
