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
	"strconv"
	"strings"
	"wallforfry/esiee-api/utils"
)

var logger = utils.InitLogger("aurion-parser")

type GroupTable struct {
	XMLName xml.Name   `xml:"table"`
	Rows    []GroupRow `xml:"row"`
}

type GroupRow struct {
	XMLName xml.Name `xml:"row",csv:"-"`
	Login   string   `xml:"login.Individu",csv:"login.Individu"`
	Mail    string   `xml:"Coordonnée.Coordonnée",csv:"Coordonnée.Coordonnée"`
	Group   string   `xml:"Code.Groupe",csv:"Code.Groupe"`
}

func (a GroupRow) String() string {
	return fmt.Sprintf("\t Username : %s, Email : %s, Groups : %s \n", a.Login, a.Mail, a.Group)
}

func (a GroupRow) CSV() []string {
	return []string{a.Login, a.Mail, a.Group}
}

type UniteTable struct {
	XMLName xml.Name   `xml:"table"`
	Rows    []UniteRow `xml:"row"`
}

type UniteRow struct {
	XMLName xml.Name `xml:"row",csv:"-"`
	Code    string   `xml:"Code.Unité",csv:"Code.Unité"`
	Label   string   `xml:"Libellé.Unité",csv:"Libellé.Unité"`
}

func (a UniteRow) String() string {
	return fmt.Sprintf("Code : %s, Label : %s\n", a.Code, a.Label)
}

func (a UniteRow) CSV() []string {
	return []string{a.Code, a.Label}
}

type Unite struct {
	Code  string
	Label string
}

type GroupEntry struct {
	Unite  string
	Groups []string
}

func (g GroupEntry) String() string {
	return fmt.Sprintf("%s %s\n", g.Unite, g.Groups)
}

func retrieveUnites(username string, password string) {
	requestId := strconv.Itoa(viper.GetInt("aurion.unitesRequestId"))

	url := "https://webaurion.esiee.fr/ws/services/executeFavori"
	data := "<service><user>" + username + "</user><password>" + password + "</password><database>prod</database><dataxml><![CDATA[<executeFavori><favori><id>" + requestId + "</id></favori><database>prod</database></executeFavori>]]></dataxml></service>"

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

	var c UniteTable
	err = xml.Unmarshal(resp.Body(), &c)
	utils.CheckError(logger, "Can't unmarshall aurion xml", err)

	logger.Info("Writing BDE_UNITES.csv")

	groupsFile, err := os.OpenFile("BDE_UNITES.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	utils.CheckError(logger, "Can't open BDE_UNITES.csv", err)
	defer groupsFile.Close()

	writer := csv.NewWriter(groupsFile)
	defer writer.Flush()

	err = writer.Write([]string{"Code.Unité", "Libellé.Unité"})
	utils.CheckError(logger, "Cannot write to file", err)

	for _, value := range c.Rows {
		value.Code = transformGroup(value.Code).Unite
		err := writer.Write(value.CSV())
		utils.CheckError(logger, "Cannot write to file", err)
	}
}

func retrieveGroups(username string, password string) {
	requestId := strconv.Itoa(viper.GetInt("aurion.groupsRequestId"))

	url := "https://webaurion.esiee.fr/ws/services/executeFavori"
	data := "<service><user>" + username + "</user><password>" + password + "</password><database>prod</database><dataxml><![CDATA[<executeFavori><favori><id>" + requestId + "</id></favori><database>prod</database></executeFavori>]]></dataxml></service>"

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

	var c GroupTable
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
	username := viper.GetString("aurion.username")
	password := viper.GetString("aurion.password")

	retrieveGroups(username, password)
	retrieveUnites(username, password)
}

func transformGroup(aurionFormat string) GroupEntry {
	adeFormat := aurionFormat

	yearRe := regexp.MustCompile(`^\d\d_`)
	adeFormat = yearRe.ReplaceAllString(adeFormat, "")

	promoRe := regexp.MustCompile(`^\w{1,4}_`)
	adeFormat = promoRe.ReplaceAllString(adeFormat, "")

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

func GetUnites() {
	username := viper.GetString("aurion.username")
	password := viper.GetString("aurion.password")
	retrieveUnites(username, password)
}

func GetUnite(code string) Unite {
	groupsFile, err := os.OpenFile("BDE_UNITES.csv", os.O_RDONLY, os.ModePerm)
	utils.CheckError(logger, "Can't open BDE_UNITES.csv", err)
	defer groupsFile.Close()

	reader := csv.NewReader(groupsFile)

	_, err = reader.Read()
	utils.CheckError(logger, "Error reading BDE_UNITES.csv", err)

	for i := 0; ; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		}
		utils.CheckError(logger, "Error reading BDE_UNITES.csv", err)

		convertedCode := strings.ReplaceAll(record[0], "_", "-")
		if convertedCode == code {
			return Unite{Code: convertedCode, Label: record[1]}
		}
	}
	return Unite{Code: code, Label: code}
}
