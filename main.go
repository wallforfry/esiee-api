package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"wallforfry/esiee-api/ade"
	"wallforfry/esiee-api/aurion"
	"wallforfry/esiee-api/database"
	_ "wallforfry/esiee-api/docs"
	"wallforfry/esiee-api/utils"
)

var logger = utils.InitLogger("main-logger")

const ParameterValueSeparatorCharacter = ","

func updateLocalCache() {
	logger.Info("Updating local cache")
	aurion.Aurion()
	ade.DownloadXml()
	ade.XmlToJson()
	logger.Info("Cache is up-to-date")
}

// @title ESIEE API
// @version 0.6.2
// @description API pour ade et aurion

// @host ade.wallforfry.fr
// @BasePath /
func main() {

	/*inputString := flag.String("input-string", "", "A sample input string. (Required)")
	  flag.Parse()

	  logger.Infof("Received inputString: %s", *inputString)

	  for _, inputStringPart := range utils.SplitStringParameter(*inputString, ParameterValueSeparatorCharacter) {
	  	logger.Infof("Parsed input value: %s", inputStringPart)
	  }*/

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	utils.CheckError(logger, "Can't read config file", err)

	debug := viper.GetBool("global.debug")

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	database.CreateMongoDatabase(
		viper.GetString("mongodb.host"),
		viper.GetInt("mongodb.port"),
		viper.GetString("mongodb.database"),
		viper.GetString("mongodb.username"),
		viper.GetString("mongodb.password"),
	)

	//uniteRepo := unite.NewMongoRepository(Database)
	//result, err := uniteRepo.Store(&unite.Unite{Code: "IGI-101", Label: "Info 101"})
	//if err == nil {
	//    fmt.Println(result.String())
	//}
	//fmt.Println(err)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", ping)

	r.GET("/status", status)

	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//Old api
	old := r.Group("/api/ade-esiee")
	{
		old.POST("/agenda", postAgendaOld)
		old.GET("/agenda/:mail", getAgendaOld)
	}

	r.POST("/agenda", postAgendaOldShort)
	r.GET("/agenda/:mail", getAgendaOldShort)

	r.GET("/api/ics/:mail", getICS)

	r.GET("/rooms", rooms)
	r.GET("/rooms/:hour", rooms)
	r.GET("/api/rooms", rooms)
	r.GET("/api/rooms/:hour", rooms)

	// New api

	v2 := r.Group("/v2")
	{
		v2.GET("/agenda/:mail", getAgenda)

		v2.GET("/ics/:mail", getICS)

		v2.GET("/groups/:mail", getGroups)

		v2.GET("/events/:name", getEventFilterByUnite)

		v2.GET("/unite/:name", getUniteInfo)
	}

	utils.Init()

	logger.Infof("Running in debug : %t", debug)

	if viper.GetBool("global.refreshCache") {
		gocron.Every(viper.GetUint64("global.refreshInterval")).Minutes().Do(updateLocalCache)
		logger.Info("Starting Gocron")
		updateLocalCache()
		gocron.Start()
	}

	r.Run()
}
