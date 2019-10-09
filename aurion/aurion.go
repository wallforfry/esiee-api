package aurion

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"io"
	"os"
	"regexp"
	"wallforfry/esiee-api/utils"
)

var logger = utils.InitLogger("aurion-parser")

type Table struct {
	XMLName xml.Name    `xml:"table"`
	Rows    []AurionRow `xml:"row"`
}

type AurionRow struct {
	XMLName xml.Name `xml:"row",csv:"-"`
	Login   string   `xml:"login.Individu",csv:"login.Individu"`
	Mail    string   `xml:"Coordonnée.Coordonnée",csv:"Coordonnée.Coordonnée"`
	Group   string   `xml:"Code.Groupe",csv:"Code.Groupe"`
}

func (a AurionRow) String() string {
	return fmt.Sprintf("\t Username : %s, Email : %s, Groups : %s \n", a.Login, a.Mail, a.Group)
}

func (a AurionRow) CSV() []string {
	return []string{a.Login, a.Mail, a.Group}
}

type GroupEntry struct {
	Unite  string
	Groups []string
}

func (g GroupEntry) String() string {
	return fmt.Sprintf("%s %s\n", g.Unite, g.Groups)
}

func callAurionApi(username string, password string) {
	url := "https://webaurion.esiee.fr/ws/services/executeFavori"
	data := "<service><user>" + username + "</user><password>" + password + "</password><database>prod</database><dataxml><![CDATA[<executeFavori><favori><id>18152763</id></favori><database>prod</database></executeFavori>]]></dataxml></service>"

	client := resty.New()

	logger.Info("Connecting to Aurion")

	resp, err := client.R().
		EnableTrace().
		SetHeader("Content-Type", "text/plain").
		SetBody([]byte(data)).
		Post(url)

	utils.CheckError(logger, "Error during Aurion Api call", err)

	if resp.StatusCode() != 200 {
		logger.Errorf("Aurion Api answers with %s code", resp.StatusCode())
	}

	var c Table
	err = xml.Unmarshal(resp.Body(), &c)
	utils.CheckError(logger, "Can't unmarshall aurion xml", err)

	logger.Info("Writing BDE_MES_GROUPES.csv")

	groupsFile, err := os.OpenFile("BDE_MES_GROUPES.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	utils.CheckError(logger, "Can't open BDE_MES_GROUPES.csv", err)
	defer groupsFile.Close()

	writer := csv.NewWriter(groupsFile)
	defer writer.Flush()

	err = writer.Write([]string{"login.Individu", "Coordonnée.Coordonnée", "Code.Groupe"})
	utils.CheckError(logger, "Cannot write to file", err)

	for _, value := range c.Rows {
		err := writer.Write(value.CSV())
		utils.CheckError(logger, "Cannot write to file", err)
	}
}

func Aurion() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	utils.CheckError(logger, "Can't read config file", err)

	username := viper.GetString("aurion.username")
	password := viper.GetString("aurion.password")

	callAurionApi(username, password)
}

func transformGroup(aurionFormat string) GroupEntry {
	adeFormat := aurionFormat
	//*
	yearRe := regexp.MustCompile(`^\d\d_`)
	adeFormat = yearRe.ReplaceAllString(adeFormat, "")
	//*
	promoRe := regexp.MustCompile(`^\w{1,4}_`)
	adeFormat = promoRe.ReplaceAllString(adeFormat, "")
	//*/

	uniteRe := regexp.MustCompile(`^([^_]*)_([^_]*)_(.*)`)
	result := uniteRe.FindStringSubmatch(adeFormat)
	if len(result) == 0 {
		return GroupEntry{Unite: adeFormat, Groups: []string{}}
	}

	unite := fmt.Sprintf("%s-%s", result[1], result[2])
	return GroupEntry{Unite: unite, Groups: result[3:]}
}

func GetUserGroups(username string) []GroupEntry {
	var groups []GroupEntry

	groupsFile, err := os.OpenFile("BDE_MES_GROUPES.csv", os.O_RDONLY, os.ModePerm)
	utils.CheckError(logger, "Can't open BDE_MES_GROUPES.csv", err)
	defer groupsFile.Close()

	reader := csv.NewReader(groupsFile)

	_, err = reader.Read()
	utils.CheckError(logger, "Error reading BDE_MES_GROUPES.csv", err)

	for i := 0; ; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		}
		utils.CheckError(logger, "Error reading BDE_MES_GROUPES.csv", err)

		if record[0] == username || record[1] == username {
			groups = append(groups, transformGroup(record[2]))
			//fmt.Printf("Row %d : %v \n", i, record)
		}
	}
	return groups
}
