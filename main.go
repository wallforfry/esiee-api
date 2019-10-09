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
				"BDE_UNITES.csv":      utils.FileInfos("BDE_UNITES.csv"),
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

		v2.GET("/events/:name", func(context *gin.Context) {
			name := context.Param("name")
			var events []ade.EventAde
			for _, event := range ade.GetEvents() {
				if event.Unite == name {
					events = append(events, event)
				}
			}
			context.JSON(200, events)
		})

		v2.GET("/unite/:name", func(context *gin.Context) {
			name := context.Param("name")
			context.JSON(200, aurion.GetUnite(name))
		})
	}

	gocron.Every(viper.GetUint64("global.refreshInterval")).Minutes().Do(updateLocalCache)

	utils.Init()

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	utils.CheckError(logger, "Can't read config file", err)

	debug := viper.GetBool("global.debug")

	logger.Infof("Running in debug : %t", debug)

	if viper.GetBool("global.refreshCache") {
		updateLocalCache()
		go r.Run() // listen and serve on 0.0.0.0:8080
		<-gocron.Start()
	} else {
		r.Run()
	}

}
