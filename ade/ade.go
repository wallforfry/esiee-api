package ade

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"wallforfry/esiee-api/utils"
)

var logger = utils.InitLogger("ade-logger")

func DownloadXml() {

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	utils.CheckError(logger, "Can't read config file", err)

	projectId := viper.GetInt("ade.project_id")

	client := resty.New()

	base_url := "https://planif.esiee.fr/jsp/webapi"

	logger.Info("Connecting to ADE")
	url := base_url + "?function=connect&login=lecteur1&password="

	resp, err := client.R().
		EnableTrace().
		Get(url)
	utils.CheckError(logger, "Can't connect to ADE", err)

	pat := regexp.MustCompile(`id=\"(.*)\"`)
	session_id := pat.FindStringSubmatch(resp.String())[1]

	logger.Info("Setting ADE project id")
	url = base_url + "?sessionId=" + session_id + "&function=setProject&projectId=" + strconv.Itoa(projectId)

	resp, err = client.R().
		EnableTrace().
		Get(url)
	utils.CheckError(logger, "Can't set ADE project id", err)

	logger.Info("Retrieving ADE xml file")
	url = base_url + "?sessionId=" + session_id + "&function=getEvents&tree=true&detail=8"

	resp, err = client.R().
		EnableTrace().
		Get(url)
	utils.CheckError(logger, "Can't retrieve ADE xml", err)

	adeFile, err := os.OpenFile("ade.xml", os.O_RDWR|os.O_CREATE, os.ModePerm)
	utils.CheckError(logger, "Can't open ade.xml", err)
	defer adeFile.Close()

	logger.Info("Writing xml file")
	n, err := adeFile.Write(resp.Body())
	utils.CheckError(logger, "Writing xml file", err)
	logger.Infof("Write %d bytes", n)

	logger.Info("Disconnecting from ADE")
	url = base_url + "?function=disconnect"

	resp, err = client.R().
		EnableTrace().
		Get(url)
	utils.CheckError(logger, "Can't disconnect from ADE", err)

	//////////////

	/*
	   client := resty.New()

	   project_id := viper.GetInt("ade.project_id")

	   url := "https://planif.esiee.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=147,738,739,743,744,2841,5757,746,747,748,2781,2782,3286,682,683,684,685,659,665,674,680,681,727,733,785,998,1295,2555,2743,5215,5688,731,734,735,736,740,741,742,780,782,1852,2584,4350,5321,786,787,788,789,790,2270,2275,2277,2278,2282,704,745,773,775,776,4937,728,2117,772,719,2112,183,185,196,4051,4679,2072,2074,2272,2276,2089,154,713,163,167,700,701,705,707,708,712,714,715,716,724,725,726,737,749,758,759,1057,1858,1908,2090,2108,2281,428,717,720,721,722,2265,2274,2279&projectId="+strconv.Itoa(project_id)+"&calType=ical&nbWeeks=12"

	   resp, err := client.R().
	       EnableTrace().
	       Get(url)
	   checkError("Download ICS file from ADE", err)

	   adeFile, err := os.OpenFile("ade.ics", os.O_RDWR|os.O_CREATE, os.ModePerm)
	   checkError("Can't open ade.ics", err)
	   defer adeFile.Close()

	   _, err = adeFile.Write(resp.Body())
	   checkError("Can't write ade.ics", err)
	   logger.Info("ICS file is up to date")
	*/

	/*
	   icsFile, err := os.OpenFile("ade.ics", os.O_RDONLY, os.ModePerm)
	   checkError("Can't open ade.ics", err)
	   defer icsFile.Close()

	   start, end := time.Now(), time.Now().Add(12*30*24*time.Hour)

	   c := gocal.NewParser(icsFile)
	   c.Start, c.End = &start, &end
	   c.Parse()

	   calendarFile, err := os.OpenFile("calendar.json", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	   checkError("Can't open calendar.json", err)
	   defer calendarFile.Close()

	   value, err := json.Marshal(c.Events)
	   checkError("", err)

	   _, err = calendarFile.Write(value)
	   checkError("Can't write calendar.json", err)
	   logger.Info("Calendar.json is up to date")
	   //for _, e := range c.Events {
	   //    fmt.Printf("%s on %s\n", e.Summary, e.Start)
	   //}



	   url := "http://test.wallforfry.fr/BDE_UNITES.csv"

	   resp, err := client.R().
	       EnableTrace().
	       Get(url)
	   checkError("Can't download BDE_UNITES.csv", err)

	   unitesFile, err := os.OpenFile("BDE_UNITES.csv", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	   checkError("Can't open BDE_UNITES.csv", err)
	   defer unitesFile.Close()

	   resp.StatusCode()
	   //_, err = unitesFile.Write(resp.Body())
	   checkError("Can't write BDE_UNITES.csv", err)

	   logger.Info("Unites names are up to date")

	   // regex : ^DESCRIPTION:\\n(\d*)\\n(.*)\\n(.*)\\nAurion\\n(.*)?\(
	   // gp 1 = id chelou
	   // gp 2 = groupe ou promo
	   // gp 3 = unite
	   // gp 4 = prof
	*/
}

func GetEvents() []EventAde {
	calendarFile, err := os.OpenFile("calendar.json", os.O_RDONLY, os.ModePerm)
	utils.CheckError(logger, "Can't open calendar.json", err)
	defer calendarFile.Close()

	var calendar []EventAde

	byteValue, _ := ioutil.ReadAll(calendarFile)

	err = json.Unmarshal(byteValue, &calendar)
	utils.CheckError(logger, "Can't unmarshall calendar.json", err)

	//logger.Info("Writing calendar.json")
	return calendar
}
