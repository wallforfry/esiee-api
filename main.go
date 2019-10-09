package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
	"wallforfry/esiee-api/ade"
	"wallforfry/esiee-api/aurion"
	"wallforfry/esiee-api/matcher"
	"wallforfry/esiee-api/utils"
)

var logger = utils.InitLogger("main-logger")

const ParameterValueSeparatorCharacter = ","

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func updateLocalCache() {
	logger.Info("Updating local cache")
	aurion.Aurion()
	ade.DownloadXml()
	ade.XmlToJson()
	logger.Info("Cache is up-to-date")
}

func main() {
	/*inputString := flag.String("input-string", "", "A sample input string. (Required)")
	flag.Parse()

	logger.Infof("Received inputString: %s", *inputString)

	for _, inputStringPart := range utils.SplitStringParameter(*inputString, ParameterValueSeparatorCharacter) {
		logger.Infof("Parsed input value: %s", inputStringPart)
	}*/

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", ping)

	r.GET("/status", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"uptime": utils.Uptime().String(),
			"files": gin.H{
				"ade.xml":             utils.FileInfos("ade.xml"),
				"calendar.json":       utils.FileInfos("calendar.json"),
				"BDE_MES_GROUPES.csv": utils.FileInfos("BDE_MES_GROUPES.csv"),
			},
		})
	})

	//Old api
	old := r.Group("/api/ade-esiee")
	{
		old.POST("/agenda", func(context *gin.Context) {
			username := context.PostForm("mail")
			events := matcher.GetOldFormatEvents(username)
			if events == nil {
				events = []ade.OldFormat{}
			}
			context.JSON(200, events)
		})
		old.GET("/agenda/:mail", func(context *gin.Context) {
			username := context.Param("mail")
			events := matcher.GetOldFormatEvents(username)
			if events == nil {
				events = []ade.OldFormat{}
			}
			context.JSON(200, events)
		})
	}

	r.POST("/agenda", func(context *gin.Context) {
		username := context.PostForm("mail")
		events := matcher.GetOldFormatEvents(username)
		if events == nil {
			events = []ade.OldFormat{}
		}
		context.JSON(200, events)
	})
	r.GET("/agenda/:mail", func(context *gin.Context) {
		username := context.Param("mail")
		events := matcher.GetOldFormatEvents(username)
		if events == nil {
			events = []ade.OldFormat{}
		}
		context.JSON(200, events)
	})

	// New api

	v2 := r.Group("/v2")
	{
		v2.GET("/agenda/:mail", func(context *gin.Context) {
			username := context.Param("mail")
			context.JSON(200, matcher.GetEvents(username))
		})

		v2.GET("/groups/:mail", func(context *gin.Context) {
			username := context.Param("mail")
			context.JSON(200, aurion.GetUserGroups(username))
		})
	}

	gocron.Every(viper.GetUint64("global.refreshInterval")).Minutes().Do(updateLocalCache)

	utils.Init()

	if viper.GetBool("global.debug") {
		r.Run()
	} else {
		updateLocalCache()

		go r.Run() // listen and serve on 0.0.0.0:8080

		<-gocron.Start()
	}

}
