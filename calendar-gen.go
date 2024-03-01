package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func fridayBeforeIfWeekend(t time.Time) time.Time {
	if t.Weekday() == time.Sunday {
		return t.AddDate(0, 0, -2)
	} else if t.Weekday() == time.Saturday {
		return t.AddDate(0, 0, -1)
	} else {
		return t
	}
}

type event struct {
	subject     string
	startDate   string
	allDayEvent bool
}

func (e *event) row() []string {
	return []string{e.subject, e.startDate, strconv.FormatBool(e.allDayEvent)}
}

type events []event

func (e *events) csv(w io.Writer) error {
	csvW := csv.NewWriter(w)

	if err := csvW.Write([]string{"Subject", "Start Date", "All Day Event"}); err != nil {
		return err
	}
	for _, event := range *e {
		if err := csvW.Write(event.row()); err != nil {
			return err
		}
	}
	csvW.Flush()
	return nil
}

func aalyriaPayday() (events events) {
	addEvent := func(date time.Time) {
		events = append(events, event{
			subject:     "Payday!",
			startDate:   date.Format("01-02-2006"),
			allDayEvent: true,
		})
	}

	now := time.Now()
	year, month, _ := now.Date()
	for year < 2026 {
		midDay := fridayBeforeIfWeekend(time.Date(year, month, 15, 0, 0, 0, 0, now.Location()))
		addEvent(midDay)

		lastDay := fridayBeforeIfWeekend(time.Date(year, month, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1))
		addEvent(lastDay)

		month++
		year += int(month) / 13
		month = ((month - 1) % 12) + 1
	}
	return
}

func lastDayOfMonth(subject string) (events events) {
	now := time.Now()
	year, month, _ := now.Date()
	for year < 2026 {
		lastDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1)

		events = append(events, event{
			subject:     subject,
			startDate:   lastDayOfMonth.Format("01-02-2006"),
			allDayEvent: true,
		})

		month++
		year += int(month) / 13
		month = ((month - 1) % 12) + 1
	}
	return
}

func main() {
	// events := aalyriaPayday()
	events := lastDayOfMonth("$ - Casella Trash Bill Due")

	if err := events.csv(os.Stdout); err != nil {
		log.Printf("csv: %v", err)
	}
}
