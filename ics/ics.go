package ics

import (
	"fmt"
	ics "github.com/arran4/golang-ical"
	"strings"
	"time"
	"wallforfry/esiee-api/ade"
)

func EventsToICS(events []ade.EventAde) *ics.Calendar {
	//Load timezone
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		panic(err.Error())
	}
	time.Local = loc

	calendar := ics.NewCalendar()
	calendar.SetMethod(ics.MethodPublish)

	for _, event := range events {
		calEvent := calendar.AddEvent(event.Name)

		createDate, _ := time.Parse("01/02/2006 15:04", event.CreatedAt)
		calEvent.SetCreatedTime(createDate)

		updateDate, _ := time.Parse("01/02/2006 15:04", event.CreatedAt)
		calEvent.SetModifiedAt(updateDate)

		startAt, _ := time.Parse("02/01/2006 15:04", event.Date+" "+event.StartHour)
		calEvent.SetStartAt(startAt)

		endAt, _ := time.Parse("02/01/2006 15:04", event.Date+" "+event.EndHour)
		calEvent.SetEndAt(endAt)

		calEvent.SetSummary(event.Name)

		calEvent.SetLocation(strings.Join(event.Classrooms, ", "))

		calEvent.SetDescription(fmt.Sprintf("%s\n%s\n%s\n%s", event.UniteName, strings.Join(event.Instructors, ", "), strings.Join(event.Classrooms, ", "), strings.Join(event.Trainees, ", ")))
	}
	return calendar
}
