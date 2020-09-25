package ade

import (
	"fmt"
	"strings"
	"time"
)

type EventAde struct {
	Name        string   `json:"name"`
	StartHour   string   `json:"start_hour"`
	EndHour     string   `json:"end_hour"`
	Date        string   `json:"date"`
	Duration    int      `json:"duration"`
	Color       string   `json:"color"` // r,g,b
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Info        string   `json:"info"`
	Trainees    []string `json:"trainees"`
	Unite       string   `json:"unite"`
	UniteName   string   `json:"unite_name"`
	Instructors []string `json:"instructors"`
	Classrooms  []string `json:"classrooms"`
	Majors      []string `json:"majors"`
}

func (e EventAde) String() string {
	return fmt.Sprintf("Name : %s\n\t StartHour : %s\n\t EndHour : %s\n\t Date : %s\n\t Duration : %d\n\t Color : rgb(%s)\n\t CreatedAt : %s\n\t UpdatedAt : %s\n\t Info : %s\n\t Trainees : %v\n\t Unite : %s\n\t UniteName : %s\n\t Instructors : %v\n\t Classrooms : %v\n\t Majors : %v\n", e.Name, e.StartHour, e.EndHour, e.Date, e.Duration, e.Color, e.CreatedAt, e.UpdatedAt, e.Info, e.Trainees, e.Unite, e.UniteName, e.Instructors, e.Classrooms, e.Majors)
}

func (e EventAde) ToOldFormat() OldFormat {
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

func (e EventAde) IsAt(at time.Time) bool {
	start, _ := time.Parse("02/01/2006 15:04", fmt.Sprintf("%s %s", e.Date, e.StartHour))
	start = start.Add(-2 * time.Hour)

	end, _ := time.Parse("02/01/2006 15:04", fmt.Sprintf("%s %s", e.Date, e.EndHour))
	end = end.Add(-2 * time.Hour)

	return at.After(start) && at.Before(end)
}
