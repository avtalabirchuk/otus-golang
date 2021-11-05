package integrationtests

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	now   = time.Now()
	day   = time.Duration(time.Hour * 24)
	week  = 7 * time.Duration(day)
	month = 30 * time.Duration(day)
)

func createEvent(title string, startDate time.Time, duration time.Duration) (int, Event, error) {
	payload := Event{
		Title:       title,
		UserID:      user.ID,
		StartDate:   startDate,
		EndDate:     startDate.Add(duration),
		NotifiedFor: 1,
	}
	return ProcessEvent(EventsURL, http.MethodPost, payload)
}

var _ = Describe("Calendar", func() {
	Describe("Event", func() {
		It("should create event", func() {
			title := "Test Event"
			statusCode, event, err := createEvent(title, now.Add(month), time.Hour)

			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))
			Expect(event.ID).NotTo(Equal(0))
			Expect(event.UserID).To(Equal(user.ID))
			Expect(event.Title).To(Equal(title))
		})

		It("should not create event with empty title", func() {
			startDate := now.Add(month)
			payload := Event{
				Title:       "",
				UserID:      user.ID,
				StartDate:   startDate,
				EndDate:     startDate.Add(time.Hour),
				NotifiedFor: 1,
			}
			statusCode, result, err := ProcessError(EventsURL, http.MethodPost, payload)

			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(400))
			Expect(result.Error).To(ContainSubstring("Field validation for 'Title' failed"))
		})

		It("should update event", func() {
			_, baseEvent, _ := createEvent("Test Event", now.Add(month), time.Hour)

			updatedTitle := "Updated title"
			payload := UpdateEventPayload{
				Event: Event{
					Title: updatedTitle,
				},
			}

			statusCode, event, err := ProcessEvent(fmt.Sprintf("%s/%d", EventsURL, baseEvent.ID), http.MethodPut, payload)
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))
			Expect(event.ID).To(Equal(baseEvent.ID))
			Expect(event.UserID).To(Equal(user.ID))
			Expect(event.Title).To(Equal(updatedTitle))
		})

		It("should remove event", func() {
			_, baseEvent, _ := createEvent("Test Event", now.Add(month), time.Hour)

			period := baseEvent.StartDate.Unix()
			statusCode, result, err := ProcessEvents(fmt.Sprintf("%s/1/%d", EventsURL, period), http.MethodGet, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))
			Expect(result.Events).To(ContainElement(baseEvent))

			statusCode, _, err = ProcessEvent(fmt.Sprintf("%s/%d", EventsURL, baseEvent.ID), http.MethodDelete, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))

			statusCode, result, err = ProcessEvents(fmt.Sprintf("%s/1/%d", EventsURL, period), http.MethodGet, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))
			Expect(result.Events).NotTo(ContainElement(baseEvent))
		})

		Context("Get events", func() {
			It("should get events per day", func() {
				startDate := now.Add(month)

				_, evt1, _ := createEvent("Test Event", startDate.Add(-2*day), day)
				_, evt2, _ := createEvent("Test Event", startDate.Add(-2*day), 15*day)
				_, evt3, _ := createEvent("Test Event", startDate, 15*day)
				_, evt4, _ := createEvent("Test Event", startDate.Add(15*day), month)

				statusCode, result, err := ProcessEvents(fmt.Sprintf("%s/1/%d", EventsURL, startDate.Unix()), http.MethodGet, nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(200))

				Expect(result.Events).NotTo(ContainElement(evt1))
				Expect(result.Events).To(ContainElement(evt2))
				Expect(result.Events).To(ContainElement(evt3))
				Expect(result.Events).NotTo(ContainElement(evt4))
			})

			It("should get events per week", func() {
				startDate := now.Add(month)

				_, evt1, _ := createEvent("Test Event", startDate.Add(-2*day), day)
				_, evt2, _ := createEvent("Test Event", startDate.Add(-2*day), 15*day)
				_, evt3, _ := createEvent("Test Event", startDate.Add(1*day), day)
				_, evt4, _ := createEvent("Test Event", startDate.Add(1*day), 15*day)
				_, evt5, _ := createEvent("Test Event", startDate.Add(15*day), day)

				statusCode, result, err := ProcessEvents(fmt.Sprintf("%s/2/%d", EventsURL, startDate.Unix()), http.MethodGet, nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(200))

				Expect(result.Events).NotTo(ContainElement(evt1))
				Expect(result.Events).To(ContainElement(evt2))
				Expect(result.Events).To(ContainElement(evt3))
				Expect(result.Events).To(ContainElement(evt4))
				Expect(result.Events).NotTo(ContainElement(evt5))
			})

			It("should get events per month", func() {
				startDate := now.Add(month)

				_, evt1, _ := createEvent("Test Event", startDate.Add(-2*day), day)
				_, evt2, _ := createEvent("Test Event", startDate.Add(-2*day), 15*day)
				_, evt3, _ := createEvent("Test Event", startDate.Add(1*day), day)
				_, evt4, _ := createEvent("Test Event", startDate.Add(day+week), day)
				_, evt5, _ := createEvent("Test Event", startDate.Add(month+month), day)

				statusCode, result, err := ProcessEvents(fmt.Sprintf("%s/3/%d", EventsURL, startDate.Unix()), http.MethodGet, nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(statusCode).To(Equal(200))

				Expect(result.Events).NotTo(ContainElement(evt1))
				Expect(result.Events).To(ContainElement(evt2))
				Expect(result.Events).To(ContainElement(evt3))
				Expect(result.Events).To(ContainElement(evt4))
				Expect(result.Events).NotTo(ContainElement(evt5))
			})
		})
	})
})
