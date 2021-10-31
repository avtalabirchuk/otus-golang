package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/require"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/service/server"
)

type DataEvent struct {
	UserID    int64
	Title     string
	StartDate string
	EndDate   string
}

type JSONEvent struct {
	ID        string `json:"ID"`
	UserID    string `json:"UserID"`
	Title     string `json:"Title"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
}

type EventsResponse struct {
	Events []JSONEvent `json:"events"`
}

var (
	service Service
	repo    *repository.MemoryRepo
)

var (
	grpcAddress = "localhost:50051"
	httpAddress = "localhost:50052"
)

func createEvent(ID int64, startDate time.Time, endDate time.Time) repository.Event {
	return repository.Event{
		ID:        ID,
		UserID:    1,
		Title:     "Test Title",
		StartDate: startDate,
		EndDate:   endDate,
	}
}

func getEventsFromStorage(ids []int64) (result []JSONEvent) {
	for _, id := range ids {
		evt := (*repo.GetStorage())[id]
		result = append(result, JSONEvent{
			ID:        strconv.Itoa(int(evt.ID)),
			UserID:    strconv.Itoa(int(evt.UserID)),
			Title:     evt.Title,
			StartDate: evt.StartDate.Format(time.RFC3339),
			EndDate:   evt.EndDate.Format(time.RFC3339),
		})
	}
	return
}

func getDate(day int, month time.Month) time.Time {
	return time.Date(2009, month, day, 0, 0, 0, 0, time.UTC)
}

func getUrl() string {
	return fmt.Sprintf("http://%s/events", httpAddress)
}

func getQueryUrl(queryType server.QueryRangeType, ts int64) string {
	return fmt.Sprintf("%s/%d/%d", getUrl(), queryType, ts)
}

func setup() {
	// turn off logs for easier checking the test errors
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	repo = repository.NewMemoryRepo()
	service = New(repo)
	runServers()
}

func resetStorage() {
	repo.ClearStorage()
	fillRepo(createEvent(1, getDate(1, time.August), getDate(2, time.August)))
	fillRepo(createEvent(2, getDate(1, time.August), getDate(15, time.August)))
	fillRepo(createEvent(3, getDate(1, time.August), getDate(1, time.September)))
	fillRepo(createEvent(4, getDate(1, time.September), getDate(2, time.September)))
	fillRepo(createEvent(5, getDate(1, time.September), getDate(5, time.September)))
	fillRepo(createEvent(6, getDate(1, time.September), getDate(15, time.September)))
	fillRepo(createEvent(7, getDate(1, time.September), getDate(1, time.October)))
	fillRepo(createEvent(8, getDate(1, time.September), getDate(15, time.October)))
	fillRepo(createEvent(9, getDate(1, time.October), getDate(2, time.October)))
	fillRepo(createEvent(10, getDate(25, time.September), getDate(1, time.December)))
	fillRepo(createEvent(11, getDate(5, time.October), getDate(6, time.October)))
}

func fillRepo(event repository.Event) {
	(*repo.GetStorage())[event.ID] = event
}

func teardown() {
}

func runServers() {
	go func() {
		if err := service.RunGRPC(grpcAddress); err != nil {
			log.Fatalf("run GRPC failed with %v; want success", err)
			return
		}
	}()
	go func() {
		if err := service.RunHTTP(grpcAddress, httpAddress); err != nil {
			log.Fatalf("run GRPC failed with %v; want success", err)
			return
		}
	}()
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
	os.Exit(0)
}

func makeQuery(t *testing.T, method string, url string, body string) []byte {
	req, err := http.NewRequest(method, url, strings.NewReader(body))

	require.Nil(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)

	result, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	return result
}

func checkEventsQuery(t *testing.T, queryType server.QueryRangeType, startPeriod int64, expectedIds []int64) {
	body := makeQuery(t, http.MethodGet, getQueryUrl(queryType, startPeriod), "{}")

	var actual EventsResponse
	err := json.Unmarshal(body, &actual)
	require.Nil(t, err)

	expected := getEventsFromStorage(expectedIds)
	require.Nil(t, err)

	require.ElementsMatch(t, expected, actual.Events)
}

func TestGetDayEvents(t *testing.T) {
	resetStorage()
	checkEventsQuery(t, server.QueryRangeType_DAY, getDate(5, time.September).Unix(), []int64{5, 6, 7, 8})
}

func TestGetWeekEvents(t *testing.T) {
	resetStorage()
	checkEventsQuery(t, server.QueryRangeType_WEEK, getDate(20, time.September).Unix(), []int64{7, 8, 10})
}

func TestGetMonthEvents(t *testing.T) {
	resetStorage()
	checkEventsQuery(t, server.QueryRangeType_MONTH, getDate(1, time.October).Unix(), []int64{7, 8, 9, 10, 11})
}

func TestDeleteEvent(t *testing.T) {
	var eventId int64 = 5
	resetStorage()

	_, ok := (*repo.GetStorage())[eventId]
	require.True(t, ok)

	body := makeQuery(t, http.MethodDelete, fmt.Sprintf("%s/%d", getUrl(), eventId), "{}")
	require.Equal(t, string(body), "{}")

	_, ok = (*repo.GetStorage())[eventId]
	require.False(t, ok)
}

func TestCreateEvent(t *testing.T) {
	resetStorage()

	event := DataEvent{
		UserID:    1,
		Title:     "Test Title",
		StartDate: getDate(1, time.August).Format(time.RFC3339),
		EndDate:   getDate(2, time.August).Format(time.RFC3339),
	}

	converted, err := json.Marshal(event)
	require.Nil(t, err)

	body := makeQuery(t, http.MethodPost, getUrl(), string(converted))

	var created JSONEvent
	err = json.Unmarshal(body, &created)

	intId, err := strconv.Atoi(created.ID)
	require.Nil(t, err)

	_, ok := (*repo.GetStorage())[int64(intId)]
	require.True(t, ok)
}

func TestUpdateEvent(t *testing.T) {
	resetStorage()

	var eventID int64 = 2

	type DataObject struct{ Title string }
	event := struct {
		Event DataObject `json:"event"`
	}{DataObject{"Updated Title"}}

	require.NotEqual(t, (*repo.GetStorage())[eventID].Title, "Updated Title")

	converted, err := json.Marshal(event)
	require.Nil(t, err)

	body := makeQuery(t, http.MethodPut, fmt.Sprintf("%s/%d", getUrl(), eventID), string(converted))

	var updated JSONEvent
	err = json.Unmarshal(body, &updated)
	require.Nil(t, err)

	require.Equal(t, updated.Title, "Updated Title")
	require.Equal(t, (*repo.GetStorage())[eventID].Title, "Updated Title")
}
