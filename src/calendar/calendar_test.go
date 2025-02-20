package calendar

import (
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

	if !strings.Contains(icsData, "SUMMARY:Data Communications & Network") {
		t.Errorf("expected event summary to be module name")
	}
	if !strings.Contains(icsData, "LOCATION:A-09-05 | APU CAMPUS") {
		t.Errorf("expected event location to be correct and in correct format")
	}
	if strings.Contains(icsData, "LOCATION:A-05-05 | APU CAMPUS") {
		t.Errorf("expected G1 class to be filtered")
	}
	if strings.Contains(icsData, "LOCATION:A-09-02 | APU CAMPUS") {
		t.Errorf("expected class from different intake to be filtered")
	}
	if !strings.Contains(icsData, "LOCATION:B-04-02 | APU CAMPUS") {
		t.Errorf("expected all classes on weeks with no grouping")
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
