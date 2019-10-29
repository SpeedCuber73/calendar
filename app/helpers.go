package app

import "time"

func hasFreeTime(existingEvents []Event, start, end time.Time) bool {
	for _, event := range existingEvents {
		if (event.StartAt.Before(start) || event.StartAt.Equal(start)) && event.EndAt.After(start) {
			return false
		}
		if event.StartAt.After(start) && event.StartAt.Before(end) {
			return false
		}
	}
	return true
}
