package integrationtests

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sender", func() {
	Describe("Event", func() {
		It("should receive event", func() {
			startDate := time.Now()
			title := fmt.Sprintf("Test Send Event %d", startDate.Unix())
			payload := Event{
				Title:       title,
				Description: "Test description",
				UserID:      user.ID,
				StartDate:   startDate,
				EndDate:     startDate.Add(time.Hour),
				NotifiedFor: 1,
			}
			statusCode, event, err := ProcessEvent(EventsURL, http.MethodPost, payload)

			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))
			Expect(user.ID).To(Equal(user.ID))
			Expect(event.Title).To(Equal(title))
			Expect(event.Status).To(Equal("New"))

			// waiting for event processing
			time.Sleep(time.Duration(10 * time.Second))

			statusCode, result, err := ProcessEvents(fmt.Sprintf("%s/1/%d", EventsURL, startDate.Unix()), http.MethodGet, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(statusCode).To(Equal(200))

			event.Status = "Sent"
			Expect(result.Events).To(ContainElement(event))
		})
	})
})
