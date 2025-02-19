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
	Intake       string `json:"INTAKE"`
	ModuleID     string `json:"MODID"`
	ModuleName   string `json:"MODULE_NAME"`
	Day          string `json:"DAY"`
	Location     string `json:"LOCATION"`
	Room         string `json:"ROOM"`
	DateISO      string `json:"DATESTAMP_ISO"`
	TimeFromISO  string `json:"TIME_FROM_ISO"`
	TimeToISO    string `json:"TIME_TO_ISO"`
}

func FetchAndConvert(intake, titleFormat string) (string, error) {
	resp, err := http.Get("https://s3-ap-southeast-1.amazonaws.com/open-ws/weektimetable")
	if err != nil {
		return "", fmt.Errorf("failed to fetch timetable: %w", err)
	}
	defer resp.Body.Close()

	var entries []TimetableEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}

	cal := ics.NewCalendar()
  cal.SetName("Apspace")

	for _, entry := range entries {
		if entry.Intake != intake {
			continue
		}

		var title string
		switch titleFormat {
		case "module_code":
			parts := strings.Split(entry.ModuleID, "-")
			if len(parts) > 3 {
				title = parts[len(parts)-3]
			} else {
				title = entry.ModuleID
			}
		case "module_id":
			title = entry.ModuleID
		default:
			title = entry.ModuleName
		}

		loc := entry.Room + " | " + entry.Location
		if loc == "" {
			loc = entry.Location
		}

		start, _ := time.Parse(time.RFC3339, entry.TimeFromISO)
		end, _ := time.Parse(time.RFC3339, entry.TimeToISO)

		event := cal.AddEvent(fmt.Sprintf("%s@calendar", entry.ModuleID))
		event.SetSummary(title)
		event.SetLocation(loc)
		event.SetStartAt(start)
		event.SetEndAt(end)
	}

	return cal.Serialize(), nil
}
