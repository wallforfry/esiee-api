package aurion

import (
	"encoding/xml"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"wallforfry/esiee-api/database"
	"wallforfry/esiee-api/models/aurion"
	"wallforfry/esiee-api/pkg/group"
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
		value.Code = group.CreateGroupFromAurionEntry("", "", value.Code).Unite
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

	logger.Info("Storing user groups to database")

	groupRepo := group.NewMongoRepository(database.Database)

	var groups []group.Group

	if len(c.Rows) > 0 {
		_ = groupRepo.Dump()
		for _, value := range c.Rows {
			g := group.CreateGroupFromAurionEntry(value.Login, value.Mail, value.Group)
			groups = append(groups, g)
			utils.CheckError(logger, "Cannot create group : "+value.Group+" for user : "+value.Group, err)
		}

		err = groupRepo.StoreMany(groups)
		utils.CheckError(logger, "Cannot store groups to database", err)
	}
}

func Aurion() {
	username := viper.GetString("aurion.username")
	password := viper.GetString("aurion.password")

	retrieveGroups(username, password)
	retrieveUnites(username, password)
}

func GetUserGroups(username string) []group.Group {

	groupRepo := group.NewMongoRepository(database.Database)

	groups, err := groupRepo.FindByUsername(username)
	if err != nil {
		logger.Fatal("Can't retrieve groups for : " + username)
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
