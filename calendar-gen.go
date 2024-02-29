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

type events struct {
	events []event
}

func (e *events) add(date time.Time) {
	e.events = append(e.events, event{
		subject:     "Payday!",
		startDate:   date.Format("01-02-2006"),
		allDayEvent: true,
	})
}

func (e *events) csv(w io.Writer) error {
	csvW := csv.NewWriter(w)

	if err := csvW.Write([]string{"Subject", "Start Date", "All Day Event"}); err != nil {
		return err
	}
	for _, event := range e.events {
		if err := csvW.Write(event.row()); err != nil {
			return err
		}
	}
	csvW.Flush()
	return nil
}

func main() {
	events := events{}

	now := time.Now()
	year, month, _ := now.Date()
	for year < 2026 {
		midDay := fridayBeforeIfWeekend(time.Date(year, month, 15, 0, 0, 0, 0, now.Location()))
		events.add(midDay)

		lastDay := fridayBeforeIfWeekend(time.Date(year, month, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, -1))
		events.add(lastDay)

		month++
		year += int(month) / 13
		month = ((month - 1) % 12) + 1
	}

	if err := events.csv(os.Stdout); err != nil {
		log.Printf("csv: %v", err)
	}
}
