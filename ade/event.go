package ade

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"
	"wallforfry/esiee-api/aurion"
	"wallforfry/esiee-api/pkg/event"
)

type Event struct {
	XMLName    xml.Name   `xml:"event"`
	Id         string     `xml:"id,attr"`
	Name       string     `xml:"name,attr"`
	StartHour  string     `xml:"startHour,attr"`
	EndHour    string     `xml:"endHour,attr"`
	Date       string     `xml:"date,attr"`
	Duration   int        `xml:"duration,attr"` //hour quarters count duration
	Color      string     `xml:"color,attr"`    // color r,g,b
	Creation   string     `xml:"creation,attr"`
	LastUpdate string     `xml:"lastUpdate,attr"`
	Info       string     `xml:"info,attr"` // info
	Resources  []Resource `xml:"resources>resource"`
}

func (e Event) String() string {
	return fmt.Sprintf("Name : %s\n\t StartHour : %s\n\t EndHour : %s\n\t Date : %s\n\t Duration : %d\n\t Color : rgb(%s)\n\t Creation : %s\n\t LastUpdate : %s\n\t Info : %s\n\t Ressources : %v\n", e.Name, e.StartHour, e.EndHour, e.Date, e.Duration, e.Color, e.Creation, e.LastUpdate, e.Info, e.Resources)
}

func FromNewFormat(e event.Event) OldFormat {
	//2018-11-21T14:00:00.000Z
	start, _ := time.Parse("02/01/2006 15:04", fmt.Sprintf("%s %s", e.Date, e.StartHour))
	start = start.Add(-2 * time.Hour)
	startString := start.Format("2006-01-02T15:04:05.000Z")
	end, _ := time.Parse("02/01/2006 15:04", fmt.Sprintf("%s %s", e.Date, e.EndHour))
	end = end.Add(-2 * time.Hour)
	endString := end.Format("2006-01-02T15:04:05.000Z")
	rooms := strings.Join(e.Classrooms, ", ")
	unite := e.UniteName
	description := fmt.Sprintf("%s\n%s", strings.Join(e.Trainees, "\n"), strings.Join(e.Instructors, "\n"))
	return OldFormat{Name: e.Name, Start: startString, End: endString, Rooms: rooms, Prof: strings.Join(e.Instructors, ", "), Unite: unite, Description: description}
}

func (e Event) ToEventAde() event.Event {
	code := e.Name

	reg := regexp.MustCompile(`^(.*):`)
	codes := reg.FindStringSubmatch(code)
	if len(codes) != 0 {
		code = codes[1]
	}

	event := event.Event{
		EventId:     e.Id,
		Name:        e.Name,
		StartHour:   e.StartHour,
		EndHour:     e.EndHour,
		Date:        e.Date,
		Duration:    e.Duration,
		Color:       e.Color,
		CreatedAt:   e.Creation,
		UpdatedAt:   e.LastUpdate,
		Info:        e.Info,
		UniteName:   aurion.GetUnite(code).Label,
		Trainees:    []string{},
		Instructors: []string{},
		Classrooms:  []string{},
		Majors:      []string{},
	}

	for _, resource := range e.Resources {
		switch resource.Category {
		case "trainee":
			event.Trainees = append(event.Trainees, resource.Name)
			break
		case "instructor":
			event.Instructors = append(event.Instructors, resource.Name)
			break
		case "classroom":
			event.Classrooms = append(event.Classrooms, resource.Name)
			break
		case "category6":
			event.Unite = resource.Name
			break
		case "equipment":
			event.Majors = append(event.Majors, resource.Name)
			break
		}

	}

	return event
}
