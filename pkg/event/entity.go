package event

import (
	"fmt"
	"time"
)

type Event struct {
	EventId     string   `json:"event_id" bson:"event_id"`
	Name        string   `json:"name" bson:"name"`
	StartHour   string   `json:"start_hour" bson:"start_hour"`
	EndHour     string   `json:"end_hour" bson:"end_hour"`
	Date        string   `json:"date" bson:"date"`
	Duration    int      `json:"duration" bson:"duration"`
	Color       string   `json:"color" bson:"color"` // r,g,b
	CreatedAt   string   `json:"created_at" bson:"created_at"`
	UpdatedAt   string   `json:"updated_at" bson:"updated_at"`
	Info        string   `json:"info" bson:"info"`
	Trainees    []string `json:"trainees" bson:"trainees"`
	Unite       string   `json:"unite" bson:"unite"`
	UniteName   string   `json:"unite_name" bson:"unite_name"`
	Instructors []string `json:"instructors" bson:"instructors"`
	Classrooms  []string `json:"classrooms" bson:"classrooms"`
	Majors      []string `json:"majors" bson:"majors"`
}

func (e Event) String() string {
	return fmt.Sprintf("Name : %s\n\t StartHour : %s\n\t EndHour : %s\n\t Date : %s\n\t Duration : %d\n\t Color : rgb(%s)\n\t CreatedAt : %s\n\t UpdatedAt : %s\n\t Info : %s\n\t Trainees : %v\n\t Unite : %s\n\t UniteName : %s\n\t Instructors : %v\n\t Classrooms : %v\n\t Majors : %v\n", e.Name, e.StartHour, e.EndHour, e.Date, e.Duration, e.Color, e.CreatedAt, e.UpdatedAt, e.Info, e.Trainees, e.Unite, e.UniteName, e.Instructors, e.Classrooms, e.Majors)
}

func (e Event) IsAt(at time.Time) bool {
	start, _ := time.Parse("02/01/2006 15:04", fmt.Sprintf("%s %s", e.Date, e.StartHour))
	start = start.Add(-2 * time.Hour)

	end, _ := time.Parse("02/01/2006 15:04", fmt.Sprintf("%s %s", e.Date, e.EndHour))
	end = end.Add(-2 * time.Hour)

	return at.After(start) && at.Before(end)
}
