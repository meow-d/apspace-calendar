package calendar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

type TimetableEntry struct {
	Intake      string `json:"INTAKE"`
	ModuleID    string `json:"MODID"`
	ModuleName  string `json:"MODULE_NAME"`
	Day         string `json:"DAY"`
	Location    string `json:"LOCATION"`
	Room        string `json:"ROOM"`
	DateISO     string `json:"DATESTAMP_ISO"`
	TimeFromISO string `json:"TIME_FROM_ISO"`
	TimeToISO   string `json:"TIME_TO_ISO"`
	Grouping    string `json:"GROUPING"`
}

func FetchAndConvert(intake, group, titleFormat string) (string, error) {
	// fetch timetable
	resp, err := http.Get("https://s3-ap-southeast-1.amazonaws.com/open-ws/weektimetable")
	if err != nil {
		return "", fmt.Errorf("failed to fetch timetable: %w", err)
	}
	defer resp.Body.Close()

	// decode json
	var entries []TimetableEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}

	// group entries by week
	weeklyGroups := make(map[string]map[string]bool)
	for _, entry := range entries {
		if entry.Intake != intake {
			continue
		}

		week, err := getWeek(entry.DateISO)
		if err != nil {
			continue
		}

		if _, exists := weeklyGroups[week]; !exists {
			weeklyGroups[week] = make(map[string]bool)
		}
		weeklyGroups[week][entry.Grouping] = true
	}

	// determine if each week should ignore grouping
	assumeNoGrouping := make(map[string]bool)
	for week, groups := range weeklyGroups {
		if len(groups) == 1 && groups["G1"] {
			assumeNoGrouping[week] = true
		} else {
			assumeNoGrouping[week] = false
		}
	}

	// create calendar
	cal := ics.NewCalendar()
	cal.SetName("Apspace")

	for _, entry := range entries {
		if entry.Intake != intake {
			continue
		}

		week, err := getWeek(entry.DateISO)
		if err != nil {
			continue
		}
		if group != "" && entry.Grouping != "" && entry.Grouping != group && !assumeNoGrouping[week] {
			continue
		}

		var title string
		var module string
		var class string

		parts := strings.Split(entry.ModuleID, "-")

		if titleFormat == "" {
			titleFormat = "module_name"
		}
		switch titleFormat {
		case "module_name", "module_name_class":
			module = entry.ModuleName
		case "module_code", "module_code_class":
			module = parts[len(parts)-3]
		}

		switch titleFormat {
		case "module_name", "module_code":
			title = module
		case "module_name_class", "module_code_class":
			class = strings.Join(parts[len(parts)-2:], "-")
			title = module + " " + class
		case "module_id":
			title = entry.ModuleID
		}

		loc := entry.Room + " | " + entry.Location
		if entry.Room == "" {
			loc = entry.Location
		}

		start, _ := time.Parse(time.RFC3339, entry.TimeFromISO)
		end, _ := time.Parse(time.RFC3339, entry.TimeToISO)

		event := cal.AddEvent(fmt.Sprintf("%s@%s@calendar", start.Unix(), entry.ModuleID))
		event.SetDtStampTime(time.Now())
		event.SetSummary(title)
		event.SetLocation(loc)
		event.SetStartAt(start)
		event.SetEndAt(end)
	}

	return cal.Serialize(), nil
}

func getWeek(dateStr string) (string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week), nil
}
