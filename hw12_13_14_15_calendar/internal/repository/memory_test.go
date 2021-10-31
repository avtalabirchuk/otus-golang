package repository

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func wrap(vs ...interface{}) []interface{} {
	return vs
}

func getDate(day int, month time.Month) time.Time {
	return time.Date(2009, month, day, 0, 0, 0, 0, time.UTC)
}

func TestMemoryRepo(t *testing.T) {
	t.Run("get day events", func(t *testing.T) {
		r := NewMemoryRepo()

		r.storage[1] = Event{ID: 1, UserID: 1, StartDate: getDate(1, time.November), EndDate: getDate(2, time.November)}
		r.storage[2] = Event{ID: 2, UserID: 1, StartDate: getDate(1, time.November), EndDate: getDate(15, time.November)}
		r.storage[3] = Event{ID: 3, UserID: 1, StartDate: getDate(3, time.November), EndDate: getDate(4, time.November)}
		r.storage[4] = Event{ID: 4, UserID: 1, StartDate: getDate(15, time.November), EndDate: getDate(16, time.November)}

		result, err := r.GetDayEvents(getDate(3, time.November))
		require.Nil(t, err)

		require.Equal(t, result, []Event{r.storage[2], r.storage[3]})
	})

	t.Run("get week events", func(t *testing.T) {
		r := NewMemoryRepo()

		date01_10 := getDate(1, time.October)

		date01_11 := getDate(1, time.November)
		date02_11 := getDate(2, time.November)
		date03_11 := getDate(3, time.November)
		date05_11 := getDate(5, time.November)
		date15_11 := getDate(15, time.November)
		date25_11 := getDate(25, time.November)

		// before requested week
		r.storage[1] = Event{ID: 1, UserID: 1, StartDate: date01_10, EndDate: date01_11}
		// starts before requested week
		r.storage[2] = Event{ID: 2, UserID: 1, StartDate: date01_10, EndDate: date02_11}
		// starts and ends within requested week
		r.storage[3] = Event{ID: 3, UserID: 1, StartDate: date03_11, EndDate: date05_11}
		// starts within requested week and ends after it
		r.storage[4] = Event{ID: 4, UserID: 1, StartDate: date05_11, EndDate: date15_11}
		// starts and ends after requested week
		r.storage[5] = Event{ID: 5, UserID: 1, StartDate: date15_11, EndDate: date25_11}

		result, err := r.GetWeekEvents(date02_11)
		require.Nil(t, err)

		fmt.Printf("%+v", result)
		require.ElementsMatch(t, []Event{r.storage[2], r.storage[3], r.storage[4]}, result)
	})

	t.Run("get month events", func(t *testing.T) {
		r := NewMemoryRepo()

		date01_10 := getDate(1, time.October)
		date02_10 := getDate(2, time.October)

		date01_11 := getDate(1, time.November)
		date02_11 := getDate(2, time.November)
		date03_11 := getDate(3, time.November)
		date05_11 := getDate(5, time.November)

		date2_12 := getDate(2, time.December)
		date15_12 := getDate(15, time.December)

		// before requested month
		r.storage[1] = Event{ID: 1, UserID: 1, StartDate: date01_10, EndDate: date02_10}
		// starts before requested month
		r.storage[2] = Event{ID: 2, UserID: 1, StartDate: date01_10, EndDate: date01_11}
		// starts and ends within requested month
		r.storage[3] = Event{ID: 3, UserID: 1, StartDate: date02_11, EndDate: date03_11}
		// starts within requested month and ends after it
		r.storage[4] = Event{ID: 4, UserID: 1, StartDate: date05_11, EndDate: date2_12}
		// starts and ends after requested month
		r.storage[5] = Event{ID: 5, UserID: 1, StartDate: date2_12, EndDate: date15_12}

		result, err := r.GetMonthEvents(date01_11)
		require.Nil(t, err)

		fmt.Printf("%+v", result)

		require.ElementsMatch(t, []Event{r.storage[2], r.storage[3], r.storage[4]}, result)
	})

	t.Run("create event", func(t *testing.T) {
		r := NewMemoryRepo()

		startDate := getDate(2, time.October)
		endDate := startDate.AddDate(0, 0, 1)
		event, err := r.CreateEvent(Event{UserID: 1, StartDate: startDate, EndDate: endDate})
		require.Nil(t, err)
		require.NotNil(t, event)
		require.Equal(t, event, r.storage[event.ID])
	})

	t.Run("update event", func(t *testing.T) {
		r := NewMemoryRepo()
		var id int64 = 1
		startDate := getDate(10, time.October)
		r.storage[id] = Event{UserID: 1, StartDate: startDate}

		endDate := startDate.AddDate(0, 0, 1)
		event, err := r.UpdateEvent(id, Event{UserID: 1, EndDate: endDate})

		require.Nil(t, err)
		require.Equal(t, event, r.storage[id])
	})

	t.Run("remove event", func(t *testing.T) {
		r := NewMemoryRepo()
		var id int64 = 1
		startDate := getDate(10, time.October)
		r.storage[id] = Event{UserID: 1, StartDate: startDate}

		err := r.DeleteEvent(id)
		require.Nil(t, err)

		_, ok := r.storage[id]
		require.False(t, ok)
	})

	t.Run("remove unavailable event", func(t *testing.T) {
		r := NewMemoryRepo()
		var id int64 = 1
		startDate := getDate(10, time.October)
		r.storage[id] = Event{UserID: 1, StartDate: startDate}

		err := r.DeleteEvent(111)
		require.EqualError(t, err, fmt.Sprintf("%s", ErrEventNotFound))
	})
}

func TestRepoMultithreading(t *testing.T) {
	t.Run("test multithreading", func(t *testing.T) {
		r := NewMemoryRepo()
		wg := &sync.WaitGroup{}
		wg.Add(2)

		startDate := getDate(10, time.October)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				r.CreateEvent(Event{UserID: 1, StartDate: startDate})
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				r.DeleteEvent(int64(i))
			}
		}()

		wg.Wait()
	})
}
