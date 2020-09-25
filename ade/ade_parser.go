package ade

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"wallforfry/esiee-api/models/ade"
	"wallforfry/esiee-api/utils"
)

func XmlToJson() {
	adeFile, err := os.OpenFile("ade.xml", os.O_RDONLY, os.ModePerm)
	utils.CheckError(logger, "Can't open ade.xml", err)
	defer adeFile.Close()

	logger.Info("Reading ade.xml file")

	byteValue, _ := ioutil.ReadAll(adeFile)

	var events ade.Events

	logger.Info("Unmarshalling xml")

	err = xml.Unmarshal(byteValue, &events)
	utils.CheckError(logger, "Can't unmarshall ade.xml", err)

	//    fmt.Println(events.Events[1290].ToEventAde())

	calendarFile, err := os.OpenFile("calendar.json", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	utils.CheckError(logger, "Can't open calendar.json", err)
	defer calendarFile.Close()

	var calendar []ade.EventAde

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
