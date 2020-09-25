package ade

import (
	"encoding/xml"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"regexp"
	"strconv"
	"time"
	"wallforfry/esiee-api/database"
	"wallforfry/esiee-api/pkg/event"
	"wallforfry/esiee-api/utils"
)

var logger = utils.InitLogger("ade-logger")

func DownloadXml() {
	projectId := viper.GetInt("ade.projectId")

	client := resty.New()

	baseUrl := "https://planif.esiee.fr/jsp/webapi"

	logger.Info("Connecting to ADE")
	url := baseUrl + "?function=connect&login=lecteur1&password="

	resp, err := client.R().
		EnableTrace().
		Get(url)
	utils.CheckError(logger, "Can't connect to ADE", err)

	if resp.StatusCode() == 200 {
		pat := regexp.MustCompile(`id=\"(.*)\"`)
		session_id := pat.FindStringSubmatch(resp.String())[1]

		logger.Info("Setting ADE project id")
		url = baseUrl + "?sessionId=" + session_id + "&function=setProject&projectId=" + strconv.Itoa(projectId)

		resp, err = client.R().
			EnableTrace().
			Get(url)
		utils.CheckError(logger, "Can't set ADE project id", err)

		logger.Info("Retrieving ADE xml file")
		url = baseUrl + "?sessionId=" + session_id + "&function=getEvents&tree=true&detail=8"

		resp, err = client.R().
			EnableTrace().
			Get(url)
		utils.CheckError(logger, "Can't retrieve ADE xml", err)

		var events Events
		logger.Info("Unmarshalling xml")
		err = xml.Unmarshal(resp.Body(), &events)
		utils.CheckError(logger, "Can't unmarshall ade.xml", err)

		var calendar []event.Event

		logger.Info("Convert event format")

		if len(events.Events) > 0 {
			for _, e := range events.Events {
				calendar = append(calendar, e.ToEventAde())
			}

			eventRepo := event.NewMongoRepository(database.Database)

			logger.Info("Store events in database")
			err = eventRepo.Dump()
			utils.CheckError(logger, "Can't dump events database", err)
			err = eventRepo.StoreMany(calendar)
			utils.CheckError(logger, "Can't store events to database", err)
		} else {
			logger.Error("Events are not stored")
		}

		logger.Info("Disconnecting from ADE")
		url = baseUrl + "?function=disconnect"

		resp, err = client.R().
			EnableTrace().
			Get(url)
		utils.CheckError(logger, "Can't disconnect from ADE", err)
	} else {
		logger.Error("Error while connecting to planif.esiee.fr..")
	}

}

func GetEvents() []event.Event {
	eventRepo := event.NewMongoRepository(database.Database)

	events, err := eventRepo.FindAll()
	utils.CheckError(logger, "Can't get events", err)
	return events
}

func GetEventsAt(offset int) []string {
	allRooms := viper.GetStringSlice("global.rooms")
	events := GetEvents()

	var used []string

	loc, err := time.LoadLocation("Europe/Paris")
	utils.CheckError(logger, "Can't load location", err)
	//now := time.Date(2019,8, 22,19,0,0,1, loc)
	now := time.Now()
	at := now.Add(1 * time.Nanosecond).Add(time.Hour * time.Duration(offset)).In(loc)

	for _, event := range events {
		if event.IsAt(at) {
			for _, room := range event.Classrooms {
				used = utils.AppendIfMissing(used, room)
			}
		}
	}

	return utils.Difference(allRooms, used)
}
