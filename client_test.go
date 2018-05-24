package ticketswitch

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetURL(t *testing.T) {
	config := &Config{
		BaseURL: "https://super.awesome.tickets",
	}
	client := NewClient(config)

	req := NewRequest("GET", "events.v1", nil)
	u, err := client.getURL(req)
	if assert.Nil(t, err) {
		assert.Equal(t, u.String(), "https://super.awesome.tickets/f13/events.v1")
	}
}

func TestGetURL_bad(t *testing.T) {
	config := &Config{
		BaseURL: ":!::!::::HAHAHAH NOPE!",
	}
	client := NewClient(config)

	req := NewRequest("GET", "events.v1", nil)
	_, err := client.getURL(req)
	assert.NotNil(t, err)
}

func TestGetURL_with_values(t *testing.T) {
	config := &Config{
		BaseURL: "https://super.awesome.tickets",
	}
	client := NewClient(config)

	req := NewRequest("GET", "events.v1", nil)
	req.Values.Add("foo", "bar")
	req.Values.Add("lol", "beans")
	req.Values.Add("lol", "icoptor")

	u, err := client.getURL(req)
	if assert.Nil(t, err) {
		assert.Equal(t, "https://super.awesome.tickets/f13/events.v1?foo=bar&lol=beans&lol=icoptor", u.String())
	}
}

func TestGetURL_with_crypto_block(t *testing.T) {
	config := &Config{
		BaseURL:     "https://super.awesome.tickets",
		CryptoBlock: "abc123",
		User:        "fred_flintstone",
	}
	client := NewClient(config)
	req := NewRequest("GET", "events.v1", nil)

	u, err := client.getURL(req)

	if assert.Nil(t, err) {
		assert.Equal(t, "https://super.awesome.tickets/f13/events.v1?crypto_block=abc123&user_id=fred_flintstone", u.String())
	}

	config.User = ""
	u, err = client.getURL(req)

	assert.NotNil(t, err)
	assert.Nil(t, u)
}

func TestGetURL_with_sub_user(t *testing.T) {
	config := &Config{
		BaseURL: "https://super.awesome.tickets",
		SubUser: "bambam",
	}

	client := NewClient(config)
	req := NewRequest("GET", "events.v1", nil)

	u, err := client.getURL(req)

	if assert.Nil(t, err) {
		assert.Equal(t, "https://super.awesome.tickets/f13/events.v1?sub_id=bambam", u.String())
	}
}

func TestSetHeaders(t *testing.T) {
	config := &Config{
		User:     "fred_flintstone",
		Password: "yabadabadoo",
		Language: "en-GB",
	}

	client := NewClient(config)
	req := NewRequest("GET", "events.v1", nil)

	err := client.setHeaders(req)

	if assert.Nil(t, err) {
		assert.Equal(t, "en-GB", req.Header.Get("Accept-Language"))
		assert.Equal(t, "Basic ZnJlZF9mbGludHN0b25lOnlhYmFkYWJhZG9v", req.Header.Get("Authorization"))
	}

	req.Header = http.Header{}

	config.Language = ""
	err = client.setHeaders(req)

	if assert.Nil(t, err) {
		assert.Equal(t, "", req.Header.Get("Accept-Language"))
	}

	config.User = ""
	err = client.setHeaders(req)
	assert.NotNil(t, err)

	config.User = "fred_flintstone"
	config.Password = ""
	err = client.setHeaders(req)
	assert.NotNil(t, err)
}

func TestDo_post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/events.v1", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "en-GB", r.Header.Get("Accept-Language"))
			assert.Equal(t, "Basic ZnJlZF9mbGludHN0b25lOnlhYmFkYWJhZG9v", r.Header.Get("Authorization"))
			assert.Equal(t, "37", r.Header.Get("Content-Length"))
			body, err := ioutil.ReadAll(r.Body)
			if assert.Nil(t, err) {
				expected := `{
  "foo": "bar",
  "lol": "beans"
}
`
				assert.Equal(t, expected, string(body))
			}
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "fred_flintstone",
		Password: "yabadabadoo",
		Language: "en-GB",
	}

	client := NewClient(config)
	req := NewRequest("POST", "events.v1", map[string]string{"foo": "bar", "lol": "beans"})

	resp, err := client.Do(req)

	if assert.Nil(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
	}
}

func TestDo_get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/events.v1", r.URL.Path)
			assert.Equal(t, "en-GB", r.Header.Get("Accept-Language"))
			assert.Equal(t, "Basic ZnJlZF9mbGludHN0b25lOnlhYmFkYWJhZG9v", r.Header.Get("Authorization"))

			assert.Equal(t, "", r.Header.Get("Content-Type"))
			assert.Equal(t, "", r.Header.Get("Content-Length"))
			body, err := ioutil.ReadAll(r.Body)
			if assert.Nil(t, err) {
				assert.Equal(t, "", string(body))
			}
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "fred_flintstone",
		Password: "yabadabadoo",
		Language: "en-GB",
	}

	client := NewClient(config)
	req := NewRequest("GET", "events.v1", nil)

	resp, err := client.Do(req)

	if assert.Nil(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
	}
}

func TestDo_with_bad_base_url(t *testing.T) {
	config := &Config{
		// this url is unparseable
		BaseURL:  "::!:!:!:!:!:!:!",
		Language: "en-GB",
	}

	client := NewClient(config)
	req := NewRequest("GET", "events.v1", nil)

	_, err := client.Do(req)

	assert.NotNil(t, err)
}

func TestDo_with_header_issues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("call was made when it should have errored")
		}))
	defer server.Close()

	// missing auth info means we can't make headers
	config := &Config{
		BaseURL:  server.URL,
		Language: "en-GB",
	}

	client := NewClient(config)
	req := NewRequest("GET", "events.v1", nil)

	_, err := client.Do(req)

	assert.NotNil(t, err)
}

func TestDo_post_with_unmarshalable_body(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("call was made when it should have errored")
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bilbo_baggins",
		Password: "inthedarknessbindthem",
		Language: "en-GB",
	}

	client := NewClient(config)
	// func cannot be marshalled
	req := NewRequest("POST", "events.v1", func() { t.Fatal("this should not run") })

	_, err := client.Do(req)

	assert.NotNil(t, err)
}

func TestDo_unrequestable_request(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("call was made when it should have errored")
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bilbo_baggins",
		Password: "inthedarknessbindthem",
		Language: "en-GB",
	}

	client := NewClient(config)
	// unicode in the method is a nono
	req := NewRequest("£££££", "events.v1", nil)

	_, err := client.Do(req)

	assert.NotNil(t, err)
}

func TestDo_error_when_doing(t *testing.T) {
	config := &Config{
		// invalid protocol will trigger an error on request
		BaseURL:  "NOPE://google.com",
		User:     "bilbo_baggins",
		Password: "inthedarknessbindthem",
		Language: "en-GB",
	}

	client := NewClient(config)
	req := NewRequest("POST", "events.v1", nil)

	_, err := client.Do(req)

	assert.NotNil(t, err)
}

func TestTest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/test.v1", r.URL.Path)
			w.Write([]byte(`{
				"user_id": "bill",
				"real_name": "fred"
			}`))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	user, err := client.Test()

	if assert.Nil(t, err) {
		assert.Equal(t, "bill", user.ID)
		assert.Equal(t, "fred", user.Name)
	}
}

func TestTest_error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/test.v1", r.URL.Path)
			w.WriteHeader(401)
			w.Write([]byte(`{
				"error_code": 3,
				"error_desc": "User authentication failure"
			}`))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	user, err := client.Test()

	if assert.NotNil(t, err) {
		assert.Nil(t, user)
		assert.IsType(t, Error{}, err)
	}
}

func TestTest_request_error(t *testing.T) {
	config := &Config{
		// invalid protocol will trigger an error on request
		BaseURL:  "NOPE://google.com",
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	_, err := client.Test()

	assert.NotNil(t, err)
}

func TestTest_request_read_error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// explicitly set the content-length to many more bytes then we are
			// actually sending before closing the connection will trigger an
			// error when the request body is read
			w.Header().Set("Content-Length", "100")
			w.Write([]byte("foo"))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	_, err := client.Test()

	assert.NotNil(t, err)
}

func TestTest_request_bad_user_json(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			// the User struct expects the IsB2B flag to be a boolean and this
			// string should not unmarshal into it
			w.Write([]byte(`{"is_b2b": "HAHAHAHAH"}`))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	_, err := client.Test()

	assert.NotNil(t, err)
}

func TestUniversalParams_Universal(t *testing.T) {
	var params UniversalParams
	var values map[string]string

	values = params.Universal()
	assert.Equal(t, map[string]string{}, values)

	params = UniversalParams{
		Availability: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_avail_details"])

	params = UniversalParams{
		AvailabilityWithPerformances: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_avail_details"])
	assert.Equal(t, "1", values["req_avail_details_with_perfs"])

	params = UniversalParams{
		ExtraInfo: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_extra_info"])

	params = UniversalParams{
		Reviews: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_reviews"])

	params = UniversalParams{
		Media: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_media_triplet_one"])
	assert.Equal(t, "1", values["req_media_triplet_two"])
	assert.Equal(t, "1", values["req_media_triplet_three"])
	assert.Equal(t, "1", values["req_media_triplet_four"])
	assert.Equal(t, "1", values["req_media_triplet_five"])
	assert.Equal(t, "1", values["req_media_square"])
	assert.Equal(t, "1", values["req_media_landscape"])
	assert.Equal(t, "1", values["req_media_marquee"])
	assert.Equal(t, "1", values["req_video_iframe"])

	params = UniversalParams{
		CostRange: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_cost_range"])

	params = UniversalParams{
		CostRangeDetails: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_cost_range_details"])

	params = UniversalParams{
		BestValueOffer: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_cost_range"])
	assert.Equal(t, "1", values["req_best_value_offer"])

	params = UniversalParams{
		MaxSavingOffer: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_cost_range"])
	assert.Equal(t, "1", values["req_max_saving_offer"])

	params = UniversalParams{
		MinCostOffer: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_cost_range"])
	assert.Equal(t, "1", values["req_min_cost_offer"])

	params = UniversalParams{
		TopPriceOffer: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_cost_range"])
	assert.Equal(t, "1", values["req_top_price_offer"])

	params = UniversalParams{
		NoSinglesData: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_cost_range"])
	assert.Equal(t, "1", values["req_no_singles_data"])

	params = UniversalParams{
		SourceInfo: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["req_src_info"])

	params = UniversalParams{
		TrackingID: "abc123",
	}
	values = params.Universal()
	assert.Equal(t, "abc123", values["custom_tracking_id"])

	params = UniversalParams{
		Misc: map[string]string{
			"foo": "bar",
			"lol": "beans",
		},
	}
	values = params.Universal()
	assert.Equal(t, "bar", values["foo"])
	assert.Equal(t, "beans", values["lol"])

}

func TestPaginationParams_Pagination(t *testing.T) {
	var params PaginationParams
	var values map[string]string

	values = params.Pagination()
	assert.Equal(t, map[string]string{}, values)

	params = PaginationParams{
		PageLength: 10,
		PageNumber: 4,
	}
	values = params.Pagination()
	assert.Equal(t, "10", values["page_len"])
	assert.Equal(t, "4", values["page_no"])
}

func TestListEventParams_Params(t *testing.T) {
	var params ListEventsParams
	var values map[string]string

	values = params.Params()
	assert.Equal(t, map[string]string{}, values)

	params = ListEventsParams{
		Keywords: []string{"foo", "bar", "lol"},
	}
	values = params.Params()
	assert.Equal(t, "foo,bar,lol", values["keywords"])

	params = ListEventsParams{
		StartDate: time.Date(2015, 4, 3, 2, 1, 0, 0, time.UTC),
		EndDate:   time.Date(2015, 6, 7, 8, 9, 0, 0, time.UTC),
	}
	values = params.Params()
	assert.Equal(t, "20150403:20150607", values["date_range"])

	params = ListEventsParams{
		CountryCode: "uk",
	}
	values = params.Params()
	assert.Equal(t, "uk", values["country_code"])

	params = ListEventsParams{
		CityCode: "uk-london",
	}
	values = params.Params()
	assert.Equal(t, "uk-london", values["city_code"])

	params = ListEventsParams{
		Circle: NewCircle(12.345, 67.890, 98.765),
	}
	values = params.Params()
	assert.Equal(t, "12.345000:67.890000:98.765000", values["circle"])

	params = ListEventsParams{
		IncludeDead: true,
	}
	values = params.Params()
	assert.Equal(t, "1", values["include_dead"])

	params = ListEventsParams{
		SortOrder: SortCostAscending,
	}
	values = params.Params()
	assert.Equal(t, "cost_ascending", values["sort_order"])

	// check that it's pulling the pagination params
	params = ListEventsParams{}
	params.PageLength = 10
	values = params.Params()
	assert.Equal(t, "10", values["page_len"])

	// check that it's pulling the universal params
	params = ListEventsParams{}
	params.TrackingID = "abc123"
	values = params.Params()
	assert.Equal(t, "abc123", values["custom_tracking_id"])

}

func TestListEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/events.v1", r.URL.Path)
			r.ParseForm()
			assert.Equal(t, "foo,bar,lol", r.Form.Get("keywords"))
			w.Write([]byte(`
				{
				  "currency_details": {
					"gbp": {
					  "currency_code": "gbp"
					}
				  },
				  "results": {
					"event": [
					  {
						"event_id": "6KT"
					  },
					  {
						"event_id": "I3S"
					  }
					],
					"paging_status": {
					  "page_length": 2,
					  "page_number": 1,
					  "pages_remaining": 18,
					  "results_remaining": 35,
					  "total_unpaged_results": 39
					}
				  }
				}
			`))
		}))

	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	params := &ListEventsParams{
		Keywords: []string{"foo", "bar", "lol"},
	}
	results, err := client.ListEvents(params)

	if assert.Nil(t, err) {
		assert.Equal(t, 2, results.PagingStatus.PageLength)
		assert.Equal(t, 1, results.PagingStatus.PageNumber)
		events := results.Events
		if assert.Len(t, events, 2) {
			assert.Equal(t, events[0].ID, "6KT")
			assert.Equal(t, events[1].ID, "I3S")
		}
	}
}

func TestListEvents_request_error(t *testing.T) {
	config := &Config{
		// Invalid protocol will result in a http.Do error
		BaseURL:  "NOPE://google.com",
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	_, err := client.ListEvents(nil)

	assert.NotNil(t, err)
}

func TestListEvents_read_error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Write less bytes then we said we were going to
			w.Header().Set("Content-Length", "100")
			w.Write([]byte("foo"))
		}))

	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	_, err := client.ListEvents(nil)

	assert.NotNil(t, err)
}

func TestGetEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/events_by_id.v1", r.URL.Path)
			r.ParseForm()
			assert.Equal(t, "1AA,2BB,3CC", r.Form.Get("event_id_list"))
			w.Write([]byte(`
				{
				  "events_by_id": {
					  "1AA": {
						  "event": {
							  "event_id": "1AA"
						  }
					  },
					  "2BB": {
						  "event": {
							  "event_id": "2BB"
						  }
					  },
					  "3CC": {
						  "event": {
							  "event_id": "3CC"
						  }
					  }
				  }
				}
			`))
		}))

	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	events, err := client.GetEvents([]string{"1AA", "2BB", "3CC"}, nil)

	if err != nil {
		t.Fatal(err)
	}

	if assert.Contains(t, events, "1AA") {
		assert.Equal(t, events["1AA"].ID, "1AA")
	}
	if assert.Contains(t, events, "2BB") {
		assert.Equal(t, events["2BB"].ID, "2BB")
	}
	if assert.Contains(t, events, "3CC") {
		assert.Equal(t, events["3CC"].ID, "3CC")
	}
}

func TestGetEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/events_by_id.v1", r.URL.Path)
			r.ParseForm()
			assert.Equal(t, "1AA", r.Form.Get("event_id_list"))
			w.Write([]byte(`
				{
				  "events_by_id": {
					  "1AA": {
						  "event": {
							  "event_id": "1AA"
						  }
					  }
				  }
				}
			`))
		}))

	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	event, err := client.GetEvent("1AA", nil)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, event.ID, "1AA")
}

func TestListPerformancesParams_Params(t *testing.T) {
	var params ListPerformancesParams
	var values map[string]string

	values = params.Params()
	assert.Equal(t, map[string]string{}, values)

	params = ListPerformancesParams{
		EventID: "25DR",
	}
	values = params.Params()
	assert.Equal(t, "25DR", values["event_id"])

	params = ListPerformancesParams{
		StartDate: time.Date(2015, 4, 3, 2, 1, 0, 0, time.UTC),
		EndDate:   time.Date(2015, 6, 7, 8, 9, 0, 0, time.UTC),
	}
	values = params.Params()
	assert.Equal(t, "20150403:20150607", values["date_range"])

	params = ListPerformancesParams{}
	params.PageLength = 10
	values = params.Params()
	assert.Equal(t, "10", values["page_len"])

	params = ListPerformancesParams{}
	params.TrackingID = "abc123"
	values = params.Params()
	assert.Equal(t, "abc123", values["custom_tracking_id"])

}

func TestListPerformances(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/performances.v1", r.URL.Path)
			r.ParseForm()
			assert.Equal(t, "20150403:20150607", r.Form.Get("date_range"))
			assert.Equal(t, "ABCD", r.Form.Get("event_id"))
			w.Write([]byte(`{
				"results": {
					"has_perf_names": true,
					"auto_select": true,
					"paging_status": {
						"page_length": 10,
						"page_number": 2
					},
					"performance": [
						{
							"perf_id": "ABCD-1"
						},
						{
							"perf_id": "ABCD-2"
						}
					]
				}
			}`))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	params := &ListPerformancesParams{
		EventID:   "ABCD",
		StartDate: time.Date(2015, 4, 3, 2, 1, 0, 0, time.UTC),
		EndDate:   time.Date(2015, 6, 7, 8, 9, 0, 0, time.UTC),
	}
	results, err := client.ListPerformances(params)

	if assert.Nil(t, err) {
		assert.True(t, results.HasPerfNames)
		assert.Equal(t, 10, results.PagingStatus.PageLength)
		assert.Equal(t, 2, results.PagingStatus.PageNumber)
		perfs := results.Performances
		assert.Len(t, perfs, 2)
		assert.Equal(t, perfs[0].ID, "ABCD-1")
		assert.Equal(t, perfs[1].ID, "ABCD-2")
	}
}

func TestListPerformancesSingleResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/performances.v1", r.URL.Path)
			r.ParseForm()
			assert.Equal(t, "20150403:20150607", r.Form.Get("date_range"))
			assert.Equal(t, "ABCD", r.Form.Get("event_id"))
			w.Write([]byte(`{
				"autoselect_this_performance": true,
				"results": {
					"has_perf_names": false,
					"performance": [
						{
							"event_id": "ABCD",
							"has_pool_seats": true,
							"is_ghost": false,
							"is_limited": false,
							"perf_id": "ABCD-1"
						}
					]
				}
			}`))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	params := &ListPerformancesParams{
		EventID:   "ABCD",
		StartDate: time.Date(2015, 4, 3, 2, 1, 0, 0, time.UTC),
		EndDate:   time.Date(2015, 6, 7, 8, 9, 0, 0, time.UTC),
	}
	results, err := client.ListPerformances(params)

	if assert.Nil(t, err) {
		assert.False(t, results.HasPerfNames)
		perfs := results.Performances
		assert.Len(t, perfs, 1)
		assert.Equal(t, perfs[0].ID, "ABCD-1")
	}
}

func TestListPerformances_error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/performances.v1", r.URL.Path)
			w.WriteHeader(400)
			w.Write([]byte(`{
				"error_code": 8,
				"error_desc": "Bad data provided"
			}`))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	user, err := client.ListPerformances(nil)

	if assert.NotNil(t, err) {
		assert.Nil(t, user)
		assert.IsType(t, Error{}, err)
	}
}

func TestGetAvailability(t *testing.T) {
	availabilityJson, error := ioutil.ReadFile("test_data/availability.json")
	if error != nil {
		t.Fatalf("Cannot find test_data/availability.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/availability.v1", r.URL.Path)
			r.ParseForm()
			w.Write(availabilityJson)
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	params := GetAvailabilityParams{
		NumberOfSeats: 2,
	}
	results, err := client.GetAvailability("7AA-5", &params)

	if assert.Nil(t, err) {
		assert.NotNil(t, results)
		assert.IsType(t, Availability{}, results.Availability)
		assert.Equal(t, results.ValidQuantities, []int{1, 2, 3, 4, 5, 6})
		assert.Equal(t, results.CurrencyCode, "gbp")
		assert.Equal(t, len(results.Availability.TicketTypes), 2)
		assert.Equal(t, len(results.Availability.TicketTypes[0].PriceBands), 2)
		assert.Equal(t, len(results.Availability.TicketTypes[1].PriceBands), 2)
	}
}

func TestGetSources(t *testing.T) {
	sourcesJSON, error := ioutil.ReadFile("test_data/sources.json")
	if error != nil {
		t.Fatalf("Cannot find test_data/sources.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/sources.v1", r.URL.Path)
			r.ParseForm()
			w.Write(sourcesJSON)
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	results, err := client.GetSources(nil)

	if assert.Nil(t, err) {
		sources := results.Sources
		assert.Len(t, sources, 3)
		assert.Equal(t, sources[0].Code, "ext_test0")
		assert.Equal(t, sources[0].Description, "External Test Backend 0")
		assert.Equal(t, sources[1].Code, "ext_test1")
		assert.Equal(t, sources[1].Description, "External Test Backend 1")
		assert.Equal(t, sources[2].Code, "generic_test")
		assert.Equal(t, sources[2].Description, "Generic JSONRPC backend TEST")
	}
}

func TestGetSourcesWithReqSrcInfo(t *testing.T) {
	sourcesJSON, error := ioutil.ReadFile("test_data/sources_req_src_info.json")
	if error != nil {
		t.Fatalf("Cannot find test_data/sources_req_src_info.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/sources.v1", r.URL.Path)
			r.ParseForm()
			w.Write(sourcesJSON)
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	params := &UniversalParams{
		SourceInfo: true,
	}

	client := NewClient(config)
	results, err := client.GetSources(params)

	if assert.Nil(t, err) {
		sources := results.Sources
		assert.Len(t, sources, 3)
		assert.Equal(t, sources[0].Code, "ext_test0")
		assert.Equal(t, sources[0].Description, "External Test Backend 0")
		assert.Equal(t, sources[0].Email, "systems@ingresso.co.uk")
		assert.Equal(t, sources[0].Address, "Suite 75, Victoria Place, Wellwood Street, Belfast BT12 5FX")
		assert.Equal(t, sources[0].Class, "ext_test")
		assert.Equal(t, sources[0].Type, "ext_test_core_direct")
		assert.Equal(t, sources[0].TermsAndConditions, "A: Privacy Policy")

		assert.Equal(t, sources[1].Code, "ext_test1")
		assert.Equal(t, sources[1].Description, "External Test Backend 1")
		assert.Equal(t, sources[1].Email, "systems@ingresso.co.uk")
		assert.Equal(t, sources[1].Address, "Suite 75, Victoria Place, Wellwood Street, Belfast BT12 5FX")
		assert.Equal(t, sources[1].Class, "ext_test")
		assert.Equal(t, sources[1].Type, "ext_test_core_direct")
		assert.Equal(t, sources[1].TermsAndConditions, "A: Privacy Policy")

		assert.Equal(t, sources[2].Code, "generic_test")
		assert.Equal(t, sources[2].Description, "Generic JSONRPC backend TEST")
		assert.Equal(t, sources[2].Class, "generic_jsonrpc")
		assert.Equal(t, sources[2].Type, "iguana")
		assert.Equal(t, sources[2].TermsAndConditions, "A: Privacy Policy")
	}
}

func TestGetSourcesError(t *testing.T) {
	errorResponse := []byte(`{
  "error_code": 8,
  "error_desc": "Bad data supplied"
}`)
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/sources.v1", r.URL.Path)
			r.ParseForm()
			w.WriteHeader(400)
			w.Write(errorResponse)
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	params := &UniversalParams{
		SourceInfo: true,
	}

	client := NewClient(config)
	results, err := client.GetSources(params)

	assert.Nil(t, results)
	assert.Equal(t, err.Error(), "ticketswitch: API error 8: Bad data supplied")
}

func TestMakeReservation(t *testing.T) {
	reservationJSON, error := ioutil.ReadFile("test_data/reservation.json")
	if error != nil {
		t.Fatalf("Cannot find test_data/reservation.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/reserve.v1", r.URL.Path)
			r.ParseForm()
			w.Write(reservationJSON)
		}))
	defer server.Close()

	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	params := &MakeReservationParams{
		PerformanceID:  "6IF-B1S",
		PriceBandCode:  "C/pool",
		TicketTypeCode: "CIRCLE",
		NumberOfSeats:  3,
		Discounts:      []string{"ADULT", "CHILD", "CHILD"},
	}

	client := NewClient(config)
	results, err := client.MakeReservation(params)

	if assert.Nil(t, err) {
		assert.Equal(t, results.AllowedCountries["ad"], "Andorra")
		assert.Equal(t, results.AllowedCountries["uk"], "United Kingdom")
		assert.Equal(t, results.AllowedCountries["uy"], "Uruguay")
		assert.True(t, results.CanEditAddress)
		assert.Equal(t, results.CurrencyDetails["gbp"].Code, "gbp")
		assert.False(t, results.InputContainedUnavailableOrder)
		assert.Equal(t, results.LanguageList[0], "en")
		assert.Equal(t, results.MinutesLeftOnReserve, 15.0)
		assert.False(t, results.NeedsAgentReference)
		assert.False(t, results.NeedsEmailAddress)
		assert.False(t, results.NeedsPaymentCard)
		assert.Equal(t, results.PrefilledAddress["country_code"], "uk")
		assert.Equal(t, results.ReserveTime, time.Date(2018, 5, 24, 15, 13, 7, 0, time.UTC))
		assert.Equal(t, results.TransactionStatus, "reserved")
		assert.Equal(t, results.Trolley.TransactionUUID, "e18c20fc-042e-11e7-975c-002590326962")
		assert.Equal(t, len(results.Trolley.Bundles), results.Trolley.BundleCount)
		assert.Equal(t, results.Trolley.Bundles[0].TotalCost, decimal.NewFromFloat(76.5))
		assert.Equal(t, len(results.Trolley.Bundles[0].Orders), 1)
		assert.Equal(t, len(results.Trolley.Bundles[0].Orders[0].TicketOrdersHolder.TicketOrders), 1)
		assert.Equal(t, len(results.Trolley.Bundles[0].Orders[0].TicketOrdersHolder.TicketOrders[0].Seats), 3)
	}
}

func TestMakeReservationFailure(t *testing.T) {
	reservationJSON, error := ioutil.ReadFile("test_data/reservation_failure.json")
	if error != nil {
		t.Fatalf("Cannot find test_data/reservation_failure.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/reserve.v1", r.URL.Path)
			r.ParseForm()
			w.Write(reservationJSON)
		}))
	defer server.Close()

	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	params := &MakeReservationParams{
		PerformanceID:  "7AB-5",
		PriceBandCode:  "B/pool",
		TicketTypeCode: "CIRCLE",
		NumberOfSeats:  2,
		Seats:          []string{"H9", "H10"},
	}

	client := NewClient(config)
	results, err := client.MakeReservation(params)

	if assert.Nil(t, err) {
		assert.Equal(t, results.AllowedCountries["ad"], "Andorra")
		assert.Equal(t, results.AllowedCountries["uk"], "United Kingdom")
		assert.Equal(t, results.AllowedCountries["zw"], "Zimbabwe")
		assert.True(t, results.CanEditAddress)
		assert.Equal(t, results.CurrencyDetails["gbp"].Code, "gbp")
		assert.False(t, results.InputContainedUnavailableOrder)
		assert.Equal(t, results.LanguageList[0], "en")
		assert.Equal(t, results.MinutesLeftOnReserve, 15.0)
		assert.False(t, results.NeedsAgentReference)
		assert.False(t, results.NeedsEmailAddress)
		assert.False(t, results.NeedsPaymentCard)
		assert.Equal(t, results.ReserveTime, time.Date(2018, 5, 24, 15, 45, 49, 0, time.UTC))
		assert.Equal(t, results.TransactionStatus, "reserved")
		assert.Equal(t, results.Trolley.TransactionUUID, "U-8841ADC8-5F69-11E8-A0DD-AC1F6B466128-EC1A0BEE-LDNX")
		assert.Equal(t, len(results.Trolley.Bundles), results.Trolley.BundleCount)
		assert.Equal(t, results.Trolley.Bundles[0].TotalCost, decimal.NewFromFloat(20))
		assert.Equal(t, len(results.Trolley.Bundles[0].Orders), 0)
		assert.Equal(t, len(results.UnreservedOrders), 1)
		assert.Equal(t, results.UnreservedOrders[0].ItemNumber, 1)
		assert.Equal(t, results.UnreservedOrders[0].RequestedSeatIDs, []string{"H9", "H10"})
	}
}
