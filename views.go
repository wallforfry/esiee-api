package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"wallforfry/esiee-api/ade"
	"wallforfry/esiee-api/aurion"
	"wallforfry/esiee-api/matcher"
	"wallforfry/esiee-api/utils"
)

// ping godoc
// @Summary Ask for ping get pong
// @Description Do ping to check api
// @Tags Core
// @Accept json
// @Produce json
// @Success 200 {string} string "{"message": "pong"}
// @Router /ping [get]
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// status godoc
// @Summary Get API status
// @Description Got API informations about local files and uptime
// @Tags Core
// @Accept json
// @Produce json
// @Success 200 {object} string "API informations"
// @Router /status [get]
func status(context *gin.Context) {
	context.JSON(200, gin.H{
		"uptime": utils.Uptime().String(),
		"files": gin.H{
			"ade.xml":             utils.FileInfos("ade.xml"),
			"calendar.json":       utils.FileInfos("calendar.json"),
			"BDE_MES_GROUPES.csv": utils.FileInfos("BDE_MES_GROUPES.csv"),
			"BDE_UNITES.csv":      utils.FileInfos("BDE_UNITES.csv"),
		},
	})
}

// rooms godoc
// @Summary Get free rooms
// @Description Get all the free rooms at now or now + X hours
// @Tags Rooms
// @Accept json
// @Produce json
// @Param hour path int false "Hour shift"
// @Success 200 {array} string "Array of free rooms"
// @Router /api/rooms/{hour} [get]
func rooms(context *gin.Context) {
	v := context.Param("hour")
	hour := 0

	if v != "" {
		var err error
		hour, err = strconv.Atoi(v)
		utils.CheckError(logger, "hour argument isn't int", err)
	}
	context.JSON(200, ade.GetEventsAt(hour))
}

// postAgendaOld godoc
// @Summary Get user agenda
// @Description Get user agenda by username or e-mail
// @Tags Old
// @Param mail formData string true "Username or e-mail"
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {array} ade.OldFormat "List of events"
// @Router /api/ade-esiee/agenda [post]
func postAgendaOld(context *gin.Context) {
	username := context.PostForm("mail")
	events := matcher.GetOldFormatEvents(username)
	if events == nil {
		events = []ade.OldFormat{}
	}
	context.JSON(200, events)
}

// postAgendaOldShort godoc
// @Summary Get user agenda
// @Description Get user agenda by username or e-mail
// @Tags Old
// @Param mail formData string true "Username or e-mail"
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {array} ade.OldFormat "List of events"
// @Router /agenda [post]
func postAgendaOldShort(context *gin.Context) {
	username := context.PostForm("mail")
	events := matcher.GetOldFormatEvents(username)
	if events == nil {
		events = []ade.OldFormat{}
	}
	context.JSON(200, events)
}

// getAgendaOld godoc
// @Summary Get user agenda
// @Description Get user agenda by username or e-mail
// @Tags Old
// @Param mail path string true "Username or e-mail"
// @Accept json
// @Produce json
// @Success 200 {array} ade.OldFormat "List of events"
// @Router /api/ade-esiee/agenda/{mail} [get]
func getAgendaOld(context *gin.Context) {
	username := context.Param("mail")
	events := matcher.GetOldFormatEvents(username)
	if events == nil {
		events = []ade.OldFormat{}
	}
	context.JSON(200, events)
}

// getAgendaOldShort godoc
// @Summary Get user agenda
// @Description Get user agenda by username or e-mail
// @Tags Old
// @Param mail path string true "Username or e-mail"
// @Accept json
// @Produce json
// @Success 200 {array} ade.OldFormat "List of events"
// @Router /agenda/{mail} [get]
func getAgendaOldShort(context *gin.Context) {
	username := context.Param("mail")
	events := matcher.GetOldFormatEvents(username)
	if events == nil {
		events = []ade.OldFormat{}
	}
	context.JSON(200, events)
}

// getAgenda godoc
// @Summary Get user agenda
// @Description Get user agenda by username or e-mail
// @Tags V2,Agenda
// @Param mail path string true "Username or e-mail"
// @Accept json
// @Produce json
// @Success 200 {array} ade.EventAde "List of events"
// @Router /v2/agenda/{mail} [get]
func getAgenda(context *gin.Context) {
	username := context.Param("mail")
	context.JSON(200, matcher.GetEvents(username))
}

// getGroups godoc
// @Summary Get user groups
// @Description Get user groups by username or e-mail
// @Tags V2,Aurion
// @Param mail path string true "Username or e-mail"
// @Accept json
// @Produce json
// @Success 200 {array} aurion.GroupEntry "List of groups"
// @Router /v2/groups/{mail} [get]
func getGroups(context *gin.Context) {
	username := context.Param("mail")
	context.JSON(200, aurion.GetUserGroups(username))
}

// getEventFilterByUnite godoc
// @Summary Get events of specific unite
// @Description Get all events of specific unite with its code
// @Tags V2,Agenda
// @Param name path string true "Unite Code"
// @Accept json
// @Produce json
// @Success 200 {array} ade.EventAde "List of events"
// @Router /v2/events/{name} [get]
func getEventFilterByUnite(context *gin.Context) {
	name := context.Param("name")
	var events []ade.EventAde
	for _, event := range ade.GetEvents() {
		if event.Unite == name {
			events = append(events, event)
		}
	}
	context.JSON(200, events)
}

// getUniteInfo godoc
// @Summary Get unite information
// @Description Get unite code and label
// @Tags V2,Aurion
// @Param name path string true "Unite Code"
// @Accept json
// @Produce json
// @Success 200 {array} aurion.Unite "Unite informations"
// @Router /v2/unite/{name} [get]
func getUniteInfo(context *gin.Context) {
	name := context.Param("name")
	context.JSON(200, aurion.GetUnite(name))
}