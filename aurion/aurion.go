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
	"wallforfry/esiee-api/database"
	"wallforfry/esiee-api/models/aurion"
	"wallforfry/esiee-api/pkg/unite"
	"wallforfry/esiee-api/utils"
)

var logger = utils.InitLogger("aurion-parser")

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

	var c aurion.UniteTable
	err = xml.Unmarshal(resp.Body(), &c)
	utils.CheckError(logger, "Can't unmarshall aurion xml", err)

	logger.Info("Saving unites to database")

	uniteRepo := unite.NewMongoRepository(database.Database)

	for _, value := range c.Rows {
		value.Code = transformGroup(value.Code).Unite
		value.Code = strings.ReplaceAll(value.Code, "_", "-")

		err = uniteRepo.Update(&unite.Unite{Code: value.Code, Label: value.Label})
		logger.Debug("Store " + value.Code + " : " + value.Label)
		utils.CheckError(logger, "Cannot store unite \""+value.Code+"\" to database", err)
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

	var c aurion.GroupTable
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

func transformGroup(aurionFormat string) aurion.GroupEntry {
	adeFormat := aurionFormat

	yearRe := regexp.MustCompile(`^\d\d_`)
	adeFormat = yearRe.ReplaceAllString(adeFormat, "")

	promoRe := regexp.MustCompile(`^\w{1,4}_`)
	adeFormat = promoRe.ReplaceAllString(adeFormat, "")

	uniteRe := regexp.MustCompile(`^([^_]*)_([^_]*)_(.*)`)
	result := uniteRe.FindStringSubmatch(adeFormat)
	if len(result) == 0 {
		return aurion.GroupEntry{Unite: adeFormat, Groups: []string{}}
	}

	unite := fmt.Sprintf("%s-%s", result[1], result[2])
	return aurion.GroupEntry{Unite: unite, Groups: result[3:]}
}

func GetUserGroups(username string) []aurion.GroupEntry {
	var groups []aurion.GroupEntry

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

func GetUnite(code string) unite.Unite {
	uniteRepo := unite.NewMongoRepository(database.Database)
	convertedCode := strings.ReplaceAll(code, "_", "-")
	result, err := uniteRepo.FindByCode(convertedCode)
	if err != nil {
		return unite.Unite{Code: convertedCode, Label: "Unknown"}
	}
	return *result
}
