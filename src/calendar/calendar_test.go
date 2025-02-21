package calendar

import (
	// "fmt"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchAndConvert(t *testing.T) {
	// mock data
	mockResponse := `[
	  {
	    "INTAKE": "AFCF2411ICT",
	    "MODID": "AICT009-3-C-DCN-L-2",
	    "MODULE_NAME": "Data Communications & Network",
	    "DAY": "MON",
	    "LOCATION": "APU CAMPUS",
	    "ROOM": "A-05-05",
	    "DATESTAMP_ISO": "2025-02-17",
	    "TIME_FROM_ISO": "2025-02-17T14:00:00+08:00",
	    "TIME_TO_ISO": "2025-02-17T16:00:00+08:00",
	    "GROUPING": "G1"
	  },
	  {
	    "INTAKE": "AFCF2411ICT",
	    "MODID": "AICT009-3-C-DCN-L-2",
	    "MODULE_NAME": "Data Communications & Network",
	    "DAY": "MON",
	    "LOCATION": "APU CAMPUS",
	    "ROOM": "A-09-05",
	    "DATESTAMP_ISO": "2025-02-17",
	    "TIME_FROM_ISO": "2025-02-17T14:00:00+08:00",
	    "TIME_TO_ISO": "2025-02-17T16:00:00+08:00",
	    "GROUPING": "G2"
	  },
	  {
	    "INTAKE": "APU2390IDK",
	    "MODID": "AICT009-3-C-DCN-L-2",
	    "MODULE_NAME": "Data Communications & Network",
	    "DAY": "MON",
	    "LOCATION": "APU CAMPUS",
	    "ROOM": "A-09-02",
	    "DATESTAMP_ISO": "2025-02-17",
	    "TIME_FROM_ISO": "2025-02-17T14:00:00+08:00",
	    "TIME_TO_ISO": "2025-02-17T16:00:00+08:00",
	    "GROUPING": "G1"
	  },
	  {
	    "INTAKE": "AFCF2411ICT",
	    "MODID": "AICT009-3-C-DCN-L-2",
	    "MODULE_NAME": "Data Communications & Network",
	    "DAY": "MON",
	    "LOCATION": "APU CAMPUS",
	    "ROOM": "B-04-02",
	    "DATESTAMP_ISO": "2025-02-24",
	    "TIME_FROM_ISO": "2025-02-24T14:00:00+08:00",
	    "TIME_TO_ISO": "2025-02-24T16:00:00+08:00",
	    "GROUPING": "G1"
	  }
	]
  `

	// mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	originalURL := "https://s3-ap-southeast-1.amazonaws.com/open-ws/weektimetable"
	http.DefaultClient = &http.Client{Transport: &mockTransport{url: originalURL, mockURL: ts.URL}}

	// test
	icsData, err := FetchAndConvert("AFCF2411ICT", "G2", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	contains(t, icsData, "BEGIN:VCALENDAR", "expected ics data to start with BEGIN:VCALENDAR")
	contains(t, icsData, "END:VCALENDAR", "expected ics data to end with END:VCALENDAR")
	contains(t, icsData, "SUMMARY:Data Communications & Network", "expected event summary to be module name")
	contains(t, icsData, "LOCATION:A-09-05 | APU CAMPUS", "expected event location to be correct and in correct format")
	notContains(t, icsData, "LOCATION:A-05-05 | APU CAMPUS", "expected G1 class to be filtered")
	notContains(t, icsData, "LOCATION:A-09-02 | APU CAMPUS", "expected class from different intake to be filtered")
	contains(t, icsData, "LOCATION:B-04-02 | APU CAMPUS", "expected all classes on weeks with no grouping")

	icsData, err = FetchAndConvert("AFCF2411ICT", "G2", "module_name_class")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	contains(t, icsData, "SUMMARY:Data Communications & Network L-2", "expected event summary to be module name and class")

	icsData, err = FetchAndConvert("AFCF2411ICT", "G2", "module_code_class")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	contains(t, icsData, "SUMMARY:DCN L-2", "expected event summary to be module code and class")
}

func contains(t *testing.T, testString, expected, message string) {
	if !strings.Contains(testString, expected) {
		t.Errorf(message, expected)
	}
}

func notContains(t *testing.T, testString, expected, message string) {
	if strings.Contains(testString, expected) {
		t.Errorf(message, expected)
	}
}

type mockTransport struct {
	url     string
	mockURL string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.HasPrefix(req.URL.String(), m.url) {
		req.URL.Scheme = "http"
		req.URL.Host = strings.TrimPrefix(m.mockURL, "http://")
	}
	return http.DefaultTransport.RoundTrip(req)
}
