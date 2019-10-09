package ade

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
	"wallforfry/esiee-api/aurion"
	"wallforfry/esiee-api/utils"
)

type Resource struct {
	XMLName  xml.Name `xml:"resource"`
	Category string   `xml:"category,attr"`
	Name     string   `xml:"name,attr"`
}

func (r Resource) String() string {
	return fmt.Sprintf("\t  %s : %s\n", r.Category, r.Name)
}

type Events struct {
	XMLName xml.Name `xml:"events"`
	Events  []Event  `xml:"event"`
}

type Event struct {
	XMLName    xml.Name   `xml:"event"`
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

func (e Event) ToEventAde() EventAde {
	code := e.Name

	reg := regexp.MustCompile(`^(.*):`)
	codes := reg.FindStringSubmatch(code)
	if len(codes) != 0 {
		code = codes[1]
	}

	event := EventAde{
		Name:      e.Name,
		StartHour: e.StartHour,
		EndHour:   e.EndHour,
		Date:      e.Date,
		Duration:  e.Duration,
		Color:     e.Color,
		CreatedAt: e.Creation,
		UpdatedAt: e.LastUpdate,
		Info:      e.Info,
		UniteName: aurion.GetUnite(code).Label,
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

type EventAde struct {
	Name        string   `json:"name"`
	StartHour   string   `json:"start_hour"`
	EndHour     string   `json:"end_hour"`
	Date        string   `json:"date"`
	Duration    int      `json:"duration"` //hour quarters count duration
	Color       string   `json:"color"`    // color r,g,b
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Info        string   `json:"info"`        // info
	Trainees    []string `json:"trainees"`    //trainee
	Unite       string   `json:"unite"`       // category6
	UniteName   string   `json:"unite_name"`  //fullname not from ade
	Instructors []string `json:"instructors"` //instructor
	Classrooms  []string `json:"classrooms"`  //classroom
	Majors      []string `json:"majors"`      // equipment
}

func (e EventAde) String() string {
	return fmt.Sprintf("Name : %s\n\t StartHour : %s\n\t EndHour : %s\n\t Date : %s\n\t Duration : %d\n\t Color : rgb(%s)\n\t CreatedAt : %s\n\t UpdatedAt : %s\n\t Info : %s\n\t Trainees : %v\n\t Unite : %s\n\t UniteName : %s\n\t Instructors : %v\n\t Classrooms : %v\n\t Majors : %v\n", e.Name, e.StartHour, e.EndHour, e.Date, e.Duration, e.Color, e.CreatedAt, e.UpdatedAt, e.Info, e.Trainees, e.Unite, e.UniteName, e.Instructors, e.Classrooms, e.Majors)
}

func (e EventAde) ToOldFormat() OldFormat {
	//2018-11-21T14:00:00.000Z
	start, _ := time.Parse("02/01/2006 15:04", fmt.Sprintf("%s %s", e.Date, e.StartHour))
	startString := start.Format("2006-01-02T15:04:05.000Z")
	end, _ := time.Parse("02/01/2006 15:04", fmt.Sprintf("%s %s", e.Date, e.EndHour))
	endString := end.Format("2006-01-02T15:04:05.000Z")
	rooms := strings.Join(e.Classrooms, ", ")
	unite := e.UniteName
	description := fmt.Sprintf("%s\n%s", strings.Join(e.Trainees, "\n"), strings.Join(e.Instructors, "\n"))
	return OldFormat{Name: e.Name, Start: startString, End: endString, Rooms: rooms, Prof: strings.Join(e.Instructors, ", "), Unite: unite, Description: description}
}

type OldFormat struct {
	Name        string `json:"name"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Rooms       string `json:"rooms"`
	Prof        string `json:"prof"`
	Unite       string `json:"unite"`
	Description string `json:"description"`
}

func XmlToJson() {
	adeFile, err := os.OpenFile("ade.xml", os.O_RDONLY, os.ModePerm)
	utils.CheckError(logger, "Can't open ade.xml", err)
	defer adeFile.Close()

	logger.Info("Reading ade.xml file")

	byteValue, _ := ioutil.ReadAll(adeFile)

	var events Events

	logger.Info("Unmarshalling xml")

	err = xml.Unmarshal(byteValue, &events)
	utils.CheckError(logger, "Can't unmarshall ade.xml", err)

	//    fmt.Println(events.Events[1290].ToEventAde())

	calendarFile, err := os.OpenFile("calendar.json", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	utils.CheckError(logger, "Can't open calendar.json", err)
	defer calendarFile.Close()

	var calendar []EventAde

	logger.Info("Convert xml to json")

	for _, event := range events.Events {
		calendar = append(calendar, event.ToEventAde())
	}

	result, err := json.Marshal(calendar)
	utils.CheckError(logger, "Can't marshall events to json", err)

	logger.Info("Writing calendar.json")

	_, err = calendarFile.Write(result)
	utils.CheckError(logger, "Can't write calendar.json", err)
}
