package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/huh"

	"github.com/dan-frohlich/tabetopevents/internal/gateway/tte"
	"github.com/dan-frohlich/tabetopevents/internal/logging"
	"github.com/dan-frohlich/tabetopevents/internal/tui"
)

func main() {
	log := logging.Log{Level: logging.LogLevelInfo}
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

	var ev []tte.ConventionEvent
	ev, err = getEvents(log, s, con)
	if err != nil {
		log.Fatal("failed to get events", "con", con.ViewURI, "error", err)
	}
	log.Info("found", "event_count", len(ev))
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
