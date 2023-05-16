package ticketswitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
	ctx := SetSessionTrackingID(context.Background(), "trackingid")

	err := client.setHeaders(ctx, req)

	if assert.Nil(t, err) {
		assert.Equal(t, "en-GB", req.Header.Get("Accept-Language"))
		assert.Equal(t, "trackingid", req.Header.Get("x-request-id"))
		assert.Equal(t, "Basic ZnJlZF9mbGludHN0b25lOnlhYmFkYWJhZG9v", req.Header.Get("Authorization"))
	}

	req.Header = http.Header{}

	config.Language = ""
	err = client.setHeaders(ctx, req)

	if assert.Nil(t, err) {
		assert.Equal(t, "", req.Header.Get("Accept-Language"))
	}

	config.User = ""
	err = client.setHeaders(ctx, req)
	assert.NotNil(t, err)

	config.User = "fred_flintstone"
	config.Password = ""
	err = client.setHeaders(ctx, req)
	assert.NotNil(t, err)
}

func TestDo_post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/events.v1", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "en-GB", r.Header.Get("Accept-Language"))
			assert.Equal(t, "Basic ZnJlZF9mbGludHN0b25lOnlhYmFkYWJhZG9v", r.Header.Get("Authorization"))
			assert.Equal(t, "postid", r.Header.Get("x-request-id"))
			assert.Equal(t, "37", r.Header.Get("Content-Length"))
			body, err := io.ReadAll(r.Body)
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
	ctx := SetSessionTrackingID(context.Background(), "postid")
	resp, err := client.Do(ctx, req)

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

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
			assert.Equal(t, "foobar123", r.Header.Get("x-request-id"))

			assert.Equal(t, "", r.Header.Get("Content-Type"))
			assert.Equal(t, "", r.Header.Get("Content-Length"))
			body, err := io.ReadAll(r.Body)
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

	ctx := SetSessionTrackingID(context.Background(), "foobar123")
	resp, err := client.Do(ctx, req)

	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

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

	resp, err := client.Do(context.Background(), req)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

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

	resp, err := client.Do(context.Background(), req)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

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
	// func cannot be marshaled
	req := NewRequest("POST", "events.v1", func() { t.Fatal("this should not run") })

	resp, err := client.Do(context.Background(), req)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

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

	resp, err := client.Do(context.Background(), req)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

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

	resp, err := client.Do(context.Background(), req)
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

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
	user, err := client.Test(context.Background())

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
	user, err := client.Test(context.Background())

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
	_, err := client.Test(context.Background())

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
	_, err := client.Test(context.Background())

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
	_, err := client.Test(context.Background())

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

	params = UniversalParams{
		AddCustomer: true,
	}
	values = params.Universal()
	assert.Equal(t, "1", values["add_customer"])
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
	results, err := client.ListEvents(context.Background(), params)

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
	_, err := client.ListEvents(context.Background(), nil)

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
	_, err := client.ListEvents(context.Background(), nil)

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
	events, err := client.GetEvents(context.Background(), []string{"1AA", "2BB", "3CC"}, nil)

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
	event, err := client.GetEvent(context.Background(), "1AA", nil)

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

	params = ListPerformancesParams{UniversalParams: UniversalParams{AddCustomer: true}}
	values = params.Params()
	assert.Equal(t, "1", values["add_customer"])
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
	results, err := client.ListPerformances(context.Background(), params)

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
	results, err := client.ListPerformances(context.Background(), params)

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
	user, err := client.ListPerformances(context.Background(), nil)

	if assert.NotNil(t, err) {
		assert.Nil(t, user)
		assert.IsType(t, Error{}, err)
	}
}

func TestListPerformanceTimesForMultipleDates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/times.v1", r.URL.Path)
			r.ParseForm()
			assert.Equal(t, "ABCD", r.Form.Get("event_id"))
			w.Write([]byte(`{
                "results": {
                    "time": [
                        {
                            "iso8601_date_and_time": "2020-01-29T14:30:00Z",
                            "time_desc": "2.30 PM"
                        },
                        {
                            "iso8601_date_and_time": "2020-01-29T19:00:00Z",
                            "time_desc": "7.00 PM"
                        },
                        {
                            "iso8601_date_and_time": "2020-01-29T19:30:00Z",
                            "time_desc": "7.30 PM"
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
		EventID: "ABCD",
	}
	results, err := client.ListPerformanceTimes(context.Background(), params)
	assert.Nil(t, err)
	assert.Equal(t, len(results.Times), 3)
	assert.Equal(t, results.Times[0].Datetime.Format(time.RFC3339), "2020-01-29T14:30:00Z")
	assert.Equal(t, results.Times[0].TimeDesc, "2.30 PM")
	assert.Equal(t, results.Times[1].Datetime.Format(time.RFC3339), "2020-01-29T19:00:00Z")
	assert.Equal(t, results.Times[1].TimeDesc, "7.00 PM")
	assert.Equal(t, results.Times[2].Datetime.Format(time.RFC3339), "2020-01-29T19:30:00Z")
	assert.Equal(t, results.Times[2].TimeDesc, "7.30 PM")
}

func TestListPerformanceTimesForSingleDate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/times.v1", r.URL.Path)
			r.ParseForm()
			assert.Equal(t, "ABCD", r.Form.Get("event_id"))
			w.Write([]byte(`{
                "results": {
                    "time": [
                        {
                            "iso8601_date_and_time": "2020-01-29T14:30:00Z",
                            "time_desc": "2.30 PM"
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
		EventID: "ABCD",
	}
	results, err := client.ListPerformanceTimes(context.Background(), params)
	assert.Nil(t, err)
	assert.Equal(t, len(results.Times), 1)
	assert.Equal(t, results.Times[0].Datetime.Format(time.RFC3339), "2020-01-29T14:30:00Z")
	assert.Equal(t, results.Times[0].TimeDesc, "2.30 PM")
}

func TestListPerformanceTimesWithNoDates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/times.v1", r.URL.Path)
			r.ParseForm()
			assert.Equal(t, "ABCD", r.Form.Get("event_id"))
			w.Write([]byte(`{
                "results": {
                    "time": []
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
		EventID: "ABCD",
	}
	results, err := client.ListPerformanceTimes(context.Background(), params)
	assert.Nil(t, err)
	assert.Equal(t, len(results.Times), 0)
}

func TestGetAvailability(t *testing.T) {
	availabilityJSON, err := os.ReadFile("testdata/availability.json")
	if err != nil {
		t.Fatalf("Cannot find testdata/availability.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/availability.v1", r.URL.Path)
			r.ParseForm()
			w.Write(availabilityJSON)
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
		Discounts:     true,
	}
	results, err := client.GetAvailability(context.Background(), "7AA-5", &params)

	if assert.Nil(t, err) {
		assert.NotNil(t, results)
		assert.IsType(t, Availability{}, results.Availability)
		assert.Equal(t, results.ValidQuantities, []int{1, 2, 3, 4, 5, 6})
		assert.Equal(t, results.CurrencyCode, "gbp")
		assert.Equal(t, len(results.Availability.TicketTypes), 2)
		assert.Equal(t, len(results.Availability.TicketTypes[0].PriceBands), 2)
		assert.Equal(t, len(results.Availability.TicketTypes[1].PriceBands), 2)
		assert.Equal(t, len(results.Availability.TicketTypes[0].PriceBands[0].PossibleDiscounts.Discounts), 2)
		assert.Equal(t, results.Availability.TicketTypes[0].PriceBands[0].Desc, "TEST PB1")
	}
}

func TestGetDiscounts(t *testing.T) {
	discountsJSON, err := os.ReadFile("testdata/discounts.json")
	if err != nil {
		t.Fatalf("Cannot find testdata/discounts.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/discounts.v1", r.URL.Path)
			r.ParseForm()
			w.Write(discountsJSON)
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	results, err := client.GetDiscounts(context.Background(), "6IF-C5O", "CIRCLE", "A/pool", nil)

	if assert.Nil(t, err) {
		discounts := results.DiscountsHolder.Discounts
		assert.Len(t, discounts, 4)
		assert.Equal(t, discounts[0].Code, "ADULT")
		assert.Equal(t, discounts[0].Description, "Adult")
		assert.Equal(t, discounts[0].Seatprice, decimal.NewFromFloat(35))
		assert.Equal(t, discounts[0].Surcharge, decimal.NewFromFloat(4))
		assert.Equal(t, discounts[0].NonOfferSeatprice, decimal.NewFromFloat(35))
		assert.Equal(t, discounts[0].NonOfferSurcharge, decimal.NewFromFloat(4))
		assert.Equal(t, discounts[1].Code, "CHILD")
		assert.Equal(t, discounts[1].Description, "Child rate")
		assert.Equal(t, discounts[2].Code, "STUDENT")
		assert.Equal(t, discounts[2].Description, "Student rate")
		assert.Equal(t, discounts[3].Code, "OAP")
		assert.Equal(t, discounts[3].Description, "Senior citizen rate")
		for _, discount := range discounts {
			assert.Equal(t, discount.AllowsLeavingSingleSeats, "always")
			assert.Equal(t, discount.NumberAvailable, 6)
			assert.Equal(t, discount.IsOffer, false)
		}
	}
}

func TestGetSources(t *testing.T) {
	sourcesJSON, err := os.ReadFile("testdata/sources.json")
	if err != nil {
		t.Fatalf("Cannot find testdata/sources.json")
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
	results, err := client.GetSources(context.Background(), nil)

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
	sourcesJSON, err := os.ReadFile("testdata/sources_req_src_info.json")
	if err != nil {
		t.Fatalf("Cannot find testdata/sources_req_src_info.json")
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
	results, err := client.GetSources(context.Background(), params)

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
	results, err := client.GetSources(context.Background(), params)

	assert.Nil(t, results)
	assert.Equal(t, err.Error(), "ticketswitch: API error 8: Bad data supplied")
}

func TestGetSendMethods(t *testing.T) {
	sourcesJSON, err := os.ReadFile("testdata/send_methods.json")
	if err != nil {
		t.Fatalf("Cannot find testdata/send_methods.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/send_methods.v1", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "6IF-C30", r.URL.Query().Get("perf_id"))
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
	results, err := client.GetSendMethods(context.Background(), "6IF-C30", nil)

	if assert.Nil(t, err) {
		assert.Equal(t, results.SourceCode, "ext_test0")
		sendMethods := results.SendMethodsHolder.SendMethods
		assert.Len(t, sendMethods, 2)
		assert.Equal(t, sendMethods[0].Type, "collect")
		assert.Equal(t, sendMethods[0].Code, "COBO")
		assert.Equal(t, sendMethods[0].Cost, decimal.NewFromFloat(1.5))
		assert.Equal(t, sendMethods[1].Type, "post")
		assert.Equal(t, sendMethods[1].Code, "POST")
		assert.Equal(t, sendMethods[1].Cost, decimal.NewFromFloat(3.5))
		assert.Len(t, sendMethods[1].PermittedCountries.Countries, 2)
		assert.Equal(t, sendMethods[1].PermittedCountries.Countries[0].Code, "ie")
		assert.Equal(t, sendMethods[1].PermittedCountries.Countries[0].Desc, "Ireland")
		assert.Equal(t, sendMethods[1].PermittedCountries.Countries[1].Code, "uk")
		assert.Equal(t, sendMethods[1].PermittedCountries.Countries[1].Desc, "United Kingdom")
	}
}

func TestMakeReservation(t *testing.T) {
	reservationJSON, err := os.ReadFile("testdata/reservation.json")
	if err != nil {
		t.Fatalf("Cannot find testdata/reservation.json")
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
	results, err := client.MakeReservation(context.Background(), params)

	if assert.Nil(t, err) {
		assert.Equal(t, results.AllowedCountries["ad"], "Andorra")
		assert.Equal(t, results.AllowedCountries["uk"], "United Kingdom")
		assert.Equal(t, results.AllowedCountries["uy"], "Uruguay")
		assert.True(t, results.CanEditAddress)
		assert.Equal(t, results.CurrencyDetails["gbp"].Code, "gbp")
		assert.False(t, results.InputContainedUnavailableOrder)
		assert.Equal(t, results.Languages[0], "en")
		assert.Equal(t, results.MinutesLeftOnReserve, 15.0)
		assert.False(t, results.NeedsAgentReference)
		assert.False(t, results.NeedsEmailAddress)
		assert.False(t, results.NeedsPaymentCard)
		assert.Equal(t, results.PrefilledAddress["country_code"], "uk")
		assert.Equal(t, results.ReserveTime, time.Date(2018, 5, 24, 15, 13, 7, 0, time.UTC))
		assert.Equal(t, results.Status, "reserved")
		assert.Equal(t, results.Trolley.TransactionUUID, "e18c20fc-042e-11e7-975c-002590326962")
		assert.Equal(t, len(results.Trolley.Bundles), results.Trolley.BundleCount)
		assert.Equal(t, results.Trolley.Bundles[0].TotalCost, decimal.NewFromFloat(76.5))
		assert.Equal(t, len(results.Trolley.Bundles[0].Orders), 1)
		assert.Equal(t, len(results.Trolley.Bundles[0].Orders[0].TicketOrdersHolder.TicketOrders), 1)
		assert.Equal(t, len(results.Trolley.Bundles[0].Orders[0].TicketOrdersHolder.TicketOrders[0].Seats), 3)
	}
}

func TestMakeReservationFailure(t *testing.T) {
	reservationJSON, err := os.ReadFile("testdata/reservation_failure.json")
	if err != nil {
		t.Fatalf("Cannot find testdata/reservation_failure.json")
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
	results, err := client.MakeReservation(context.Background(), params)

	if assert.Nil(t, err) {
		assert.Equal(t, results.AllowedCountries["ad"], "Andorra")
		assert.Equal(t, results.AllowedCountries["uk"], "United Kingdom")
		assert.Equal(t, results.AllowedCountries["zw"], "Zimbabwe")
		assert.True(t, results.CanEditAddress)
		assert.Equal(t, results.CurrencyDetails["gbp"].Code, "gbp")
		assert.False(t, results.InputContainedUnavailableOrder)
		assert.Equal(t, results.Languages[0], "en")
		assert.Equal(t, results.MinutesLeftOnReserve, 15.0)
		assert.False(t, results.NeedsAgentReference)
		assert.False(t, results.NeedsEmailAddress)
		assert.False(t, results.NeedsPaymentCard)
		assert.Equal(t, results.ReserveTime, time.Date(2018, 5, 24, 15, 45, 49, 0, time.UTC))
		assert.Equal(t, results.Status, "reserved")
		assert.Equal(t, results.Trolley.TransactionUUID, "U-8841ADC8-5F69-11E8-A0DD-AC1F6B466128-EC1A0BEE-LDNX")
		assert.Equal(t, len(results.Trolley.Bundles), results.Trolley.BundleCount)
		assert.True(t, results.Trolley.Bundles[0].TotalCost.Equal(decimal.NewFromFloat(20)))
		assert.Equal(t, len(results.Trolley.Bundles[0].Orders), 0)
		assert.Equal(t, len(results.UnreservedOrders), 1)
		assert.Equal(t, results.UnreservedOrders[0].ItemNumber, 1)
		assert.Equal(t, results.UnreservedOrders[0].RequestedSeatIDs, []string{"H9", "H10"})
	}
}

func TestTransactionParams_Params(t *testing.T) {
	params := TransactionParams{
		UniversalParams: UniversalParams{
			TrackingID: "abc123",
		},
		TransactionUUID: "acb71e1e-20aa-4c74-a607-7d3580b08130",
	}
	values := params.Params()
	assert.Equal(t, "acb71e1e-20aa-4c74-a607-7d3580b08130", values["transaction_uuid"])
	assert.Equal(t, "abc123", values["custom_tracking_id"])
}

func TestReleaseReservation_success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/release.v1", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			var inputs map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&inputs); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1", inputs["transaction_uuid"])
			w.Write([]byte(`{"released_ok": true}`))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	params := &TransactionParams{
		TransactionUUID: "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1",
	}
	success, err := client.ReleaseReservation(context.Background(), params)
	if assert.Nil(t, err) {
		assert.True(t, success)
	}
}

func TestReleaseReservation_failed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"released_ok": false}`))
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	params := &TransactionParams{
		TransactionUUID: "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1",
	}
	success, err := client.ReleaseReservation(context.Background(), params)
	if assert.Nil(t, err) {
		assert.False(t, success)
	}
}

type FakePaymentMethod map[string]string

func (fake FakePaymentMethod) PaymentParams() map[string]string {
	return fake
}

func TestMakePurchaseParams_without_payment_method(t *testing.T) {
	params := MakePurchaseParams{
		UniversalParams: UniversalParams{
			TrackingID: "abc123",
		},
		TransactionUUID: "797aeb33-5199-4457-bfa1-cd505b6943b4",
		Customer: Customer{
			FirstName: "Barney",
			LastName:  "Rubble",
		},
		SendConfirmationEmail: true,
		AgentReference:        "AgentReference",
	}

	values := params.Params()
	assert.Equal(t, "797aeb33-5199-4457-bfa1-cd505b6943b4", values["transaction_uuid"])
	assert.Equal(t, "abc123", values["custom_tracking_id"])
	assert.Equal(t, "Barney", values["first_name"])
	assert.Equal(t, "Rubble", values["last_name"])
	assert.Equal(t, "AgentReference", values["agent_reference"])
	assert.Equal(t, "1", values["send_confirmation_email"])
}

func TestMakePurchaseParams_with_payment_method(t *testing.T) {
	params := MakePurchaseParams{
		PaymentMethod: FakePaymentMethod{
			"foo": "bar",
			"lol": "beans",
		},
	}

	values := params.Params()
	assert.Equal(t, "bar", values["foo"])
	assert.Equal(t, "beans", values["lol"])
}

func TestMakePurchase_success(t *testing.T) {
	data, err := os.ReadFile("testdata/purchase-credit-success.json")
	if err != nil {
		t.Fatalf("testdata/purchase-credit-success.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/purchase.v1", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			var inputs map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			if err2 := decoder.Decode(&inputs); err2 != nil {
				t.Fatal(err2)
			}
			assert.Equal(t, "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1", inputs["transaction_uuid"])
			w.Write(data)
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	params := &MakePurchaseParams{
		TransactionUUID: "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1",
	}
	result, err := client.MakePurchase(context.Background(), params)
	if assert.Nil(t, err) {
		assert.Equal(t, "purchased", result.Status)
		assert.Nil(t, result.Callout)
		assert.Nil(t, result.PendingCallout)
		assert.Equal(t, map[string]Currency{
			"gbp": {
				Code:       "gbp",
				Places:     2,
				Factor:     100,
				PreSymbol:  "£",
				PostSymbol: "",
				Number:     826,
			},
		}, result.Currency)
		assert.Equal(t, Customer{
			FirstName:                  "Test",
			LastName:                   "Tester",
			AddressLineOne:             "Metro Building",
			AddressLineTwo:             "1 Butterwick",
			CountryCode:                "uk",
			EmailAddress:               "testing@gmail.com",
			WorkPhone:                  "0203 137 7420",
			HomePhone:                  "0203 137 7420",
			Postcode:                   "W6 8DL",
			Town:                       "London",
			SupplierCanUseCustomerData: false,
			UserCanUseCustomerData:     true,
			WorldCanUseCustomerData:    false,
		}, result.Customer)
		expectedReserve := time.Date(2017, 4, 12, 8, 38, 20, 0, time.UTC)
		assert.Equal(t, expectedReserve, result.ReserveDatetime)
		expectedPurchase := time.Date(2017, 4, 12, 8, 38, 35, 0, time.UTC)
		assert.Equal(t, expectedPurchase, result.PurchaseDatetime)
		assert.Equal(t, []string{
			"en-gb",
			"en",
			"en-us",
			"nl",
		}, result.Languages)
		assert.Equal(t, "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1", result.Trolley.TransactionUUID)
	}
}

func TestGetStatusWithCustomer(t *testing.T) {
	data, err := os.ReadFile("testdata/status_with_customer.json")
	if err != nil {
		t.Fatalf("testdata/status_with_customer.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/status.v1", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1", r.URL.Query().Get("transaction_uuid"))
			w.Write(data)
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)

	customerParam := UniversalParams{
		AddCustomer: true,
	}
	params := &TransactionParams{
		TransactionUUID: "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1", UniversalParams: customerParam,
	}
	result, err := client.GetStatus(context.Background(), params)
	customer := result.Customer

	if assert.Nil(t, err) {
		assert.Equal(t, "1234567", customer.AgentReference)
		assert.Equal(t, "Fred", customer.FirstName)
		assert.Equal(t, "Flinstone", customer.LastName)
	}
}

func TestGetStatus(t *testing.T) {
	data, err := os.ReadFile("testdata/status.json")
	if err != nil {
		t.Fatalf("testdata/status.json")
	}
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/status.v1", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1", r.URL.Query().Get("transaction_uuid"))
			w.Write(data)
		}))
	defer server.Close()
	config := &Config{
		BaseURL:  server.URL,
		User:     "bill",
		Password: "hahaha",
	}

	client := NewClient(config)
	params := &TransactionParams{
		TransactionUUID: "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1",
	}
	result, err := client.GetStatus(context.Background(), params)
	if assert.Nil(t, err) {
		assert.Equal(t, "purchased", result.Status)
		assert.Equal(t, map[string]Currency{
			"gbp": {
				Code:       "gbp",
				Places:     2,
				Factor:     100,
				PreSymbol:  "£",
				PostSymbol: "",
				Number:     826,
			},
		}, result.CurrencyDetails)
		expectedReserve := time.Date(2018, 5, 27, 13, 3, 14, 0, time.UTC)
		assert.Equal(t, expectedReserve, result.ReserveDatetime)
		expectedPurchase := time.Date(2018, 5, 27, 13, 3, 15, 0, time.UTC)
		assert.Equal(t, expectedPurchase, result.PurchaseDatetime)
		assert.Equal(t, []string{"en"}, result.Languages)
		assert.Equal(t, "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1", result.Trolley.TransactionUUID)
		assert.True(t, result.Trolley.PurchaseResult.Success)
		assert.False(t, result.Trolley.PurchaseResult.IsPartial)
		assert.Equal(t, len(result.Trolley.Bundles), 1)
		assert.Equal(t, len(result.Trolley.Bundles[0].Orders), 1)
		assert.Equal(t, result.Trolley.Bundles[0].Orders[0].Event.UpsellList.EventIds, []string{"7AA", "6IF"})
	}
}

func TestCancel(t *testing.T) {
	data, err := os.ReadFile("testdata/cancel.json")
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}
	transUUID := "4df498e9-2daa-4393-a6bb-cc3dfefa7cc1"
	happyServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/f13/cancel.v1", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, transUUID, r.URL.Query().Get("transaction_uuid"))
			w.Write(data)
		}))
	defer happyServer.Close()
	jsonErrorServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("not real json"))
		}))
	defer jsonErrorServer.Close()
	table := []struct {
		name        string
		config      *Config
		shouldError bool
	}{
		{
			name: "happy path",
			config: &Config{
				BaseURL:  happyServer.URL,
				User:     "bill",
				Password: "hahaha",
			},
		},
		{
			name: "request error",
			config: &Config{
				BaseURL:  "not a real url",
				User:     "bill",
				Password: "hahaha",
			},
			shouldError: true,
		},
		{
			name: "decoder error",
			config: &Config{
				BaseURL:  jsonErrorServer.URL,
				User:     "bill",
				Password: "hahaha",
			},
			shouldError: true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient(test.config)
			params := &CancellationParams{
				TransactionUUID: transUUID,
				CancelItemsList: []int{1},
			}
			result, err := client.Cancel(context.Background(), params)
			if test.shouldError {
				if !assert.NotNil(t, err) {
					t.Fatal("expected an error!")
				}
			} else {
				if !assert.Nil(t, err) {
					t.Fatal(err)
				}
				assert.Equal(t, map[string]Currency{
					"gbp": {
						Code:       "gbp",
						Places:     2,
						Factor:     100,
						PreSymbol:  "£",
						PostSymbol: "",
						Number:     826,
					},
				}, result.CurrencyDetails)
				assert.Equal(t, transUUID, result.Trolley.TransactionUUID)
				assert.True(t, result.Trolley.PurchaseResult.Success)
				assert.False(t, result.Trolley.PurchaseResult.IsPartial)
				assert.Equal(t, len(result.Trolley.Bundles), 1)
				assert.Equal(t, len(result.Trolley.Bundles[0].Orders), 1)
				// nolint:misspell
				assert.Equal(t, result.Trolley.Bundles[0].Orders[0].CancellationStatus, "cancelled")
			}
		})
	}
}

func TestCancelItemsList_String(t *testing.T) {
	assert.Equal(t, "1,2,3,4,5", fmt.Sprint(CancelItemsList{1, 2, 3, 4, 5}))
	assert.Equal(t, "11,12,31,14,79,111", fmt.Sprint(CancelItemsList{11, 12, 31, 14, 79, 111}))
	assert.Equal(t, "1", fmt.Sprint(CancelItemsList{1}))
	assert.Equal(t, "", fmt.Sprint(CancelItemsList{}))
}
