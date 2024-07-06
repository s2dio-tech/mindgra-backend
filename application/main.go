package main

import (
	"github.com/gin-gonic/gin"

	// "github.com/jinzhu/gorm"
	"github.com/s2dio-tech/mindgra-backend/datasource"

	docs "github.com/s2dio-tech/mindgra-backend/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	common "github.com/s2dio-tech/mindgra-backend/common"
	auth "github.com/s2dio-tech/mindgra-backend/common/auth"
	_httpCommon "github.com/s2dio-tech/mindgra-backend/common/http"

	_mailService "github.com/s2dio-tech/mindgra-backend/internal/email/service"
	_mailUsecase "github.com/s2dio-tech/mindgra-backend/internal/email/usecase"

	_authHttp "github.com/s2dio-tech/mindgra-backend/internal/auth/delivery/http"
	_tokenRepo "github.com/s2dio-tech/mindgra-backend/internal/auth/repository"
	_authUsecase "github.com/s2dio-tech/mindgra-backend/internal/auth/usecase"

	_userHttp "github.com/s2dio-tech/mindgra-backend/internal/users/delivery/http"
	_userRepo "github.com/s2dio-tech/mindgra-backend/internal/users/repository"
	_userUsecase "github.com/s2dio-tech/mindgra-backend/internal/users/usecase"

	_wordHttp "github.com/s2dio-tech/mindgra-backend/internal/words/delivery/http"
	_wordRepo "github.com/s2dio-tech/mindgra-backend/internal/words/repository"
	_wordUsecase "github.com/s2dio-tech/mindgra-backend/internal/words/usecase"
)

func main() {

	// init configuration
	common.InitConfig()

	// init database
	db := datasource.InitNeo4J()
	defer db.Disconnect()

	// repositories
	userRepo := _userRepo.InitUserRepository(&db)
	tokenRepo := _tokenRepo.InitTokenRepository(&db)
	wordRepo := _wordRepo.InitWordRepository(&db)
	graphRepo := _wordRepo.InitGraphRepository(&db)
	linkRepo := _wordRepo.InitLinkRepository(&db)

	mailUsecase := _mailUsecase.Init(&_mailService.MailJet{
		PublicKey:  *common.AppConfig.MailjetPublicKey,
		PrivateKey: *common.AppConfig.MailjetPrivateKey,
	})
	// mailUsecase := _mailUsecase.Init(&_mailService.SMTPSender{
	// 	Host:     *common.AppConfig.SMTPHost,
	// 	Port:     *common.AppConfig.SMTPPort,
	// 	Username: *common.AppConfig.SMTPUsername,
	// 	Password: *common.AppConfig.SMTPPassword,
	// })
	authUsecase := _authUsecase.InitAuthUsecase(tokenRepo, userRepo, mailUsecase)
	userUsecase := _userUsecase.InitUserUsecase(userRepo, mailUsecase)
	wordUsecase := _wordUsecase.InitWordUsecase(wordRepo, graphRepo)
	linkUsecase := _wordUsecase.InitLinkUsecase(linkRepo, wordRepo)
	graphUsecase := _wordUsecase.InitGraphUsecase(graphRepo)

	///////////////////////////
	// init rest api server
	///////////////////////////
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(_httpCommon.CORSMiddleware())
	docs.SwaggerInfo.BasePath = "/"
	v1 := r.Group("")

	// http handlers
	authHandler := _authHttp.InitHandlers(authUsecase, userUsecase)
	userHandler := _userHttp.InitHandlers(userUsecase)
	wordHandler := _wordHttp.InitWordHandlers(wordUsecase)
	linkHandler := _wordHttp.InitLinkHandlers(linkUsecase)
	graphHandler := _wordHttp.InitGraphHandlers(graphUsecase)

	authGroup := v1.Group("")
	authGroup.Use(_httpCommon.CORSMiddleware())
	authGroup.Use(auth.AuthMiddleware())
	{
		//words
		authGroup.GET("/words/search", wordHandler.SearchWord)
		authGroup.GET("/words/findPath", wordHandler.FindPath)
		authGroup.GET("/words/:id", wordHandler.GetWordDetail)
		authGroup.POST("/words", wordHandler.CreateWord)
		authGroup.POST("/words/links", wordHandler.Link2Words)
		authGroup.PUT("/words/:id", wordHandler.UpdateWord)
		authGroup.DELETE("/words/:id", wordHandler.DeleteWord)
		//links
		authGroup.GET("/links/:path1", linkHandler.GetDetail)
		authGroup.GET("/links/:path1/:path2", linkHandler.GetDetail)
		authGroup.POST("/links", linkHandler.CreateLink)
		authGroup.PUT("/links/:id", linkHandler.UpdateLink)
		authGroup.DELETE("/links", linkHandler.DeleteLink)
		//graphs
		authGroup.GET("/graphs/:id/data", wordHandler.GetGraphData)
		authGroup.GET("/graphs", graphHandler.List)
		authGroup.POST("/graphs", graphHandler.CreateGraph)
		authGroup.PUT("/graphs/:id", graphHandler.UpdateGraph)
		authGroup.DELETE("/graphs/:id", graphHandler.DeleteGraph)
	}

	v1.GET("/graphs/:id", graphHandler.Detail)

	v1.POST("/auth/login", authHandler.Login)
	v1.POST("/auth/grant", authHandler.Grant)
	v1.POST("/auth/password-forgot", authHandler.ForgotPassword)
	v1.POST("/auth/password-verify-otp", authHandler.VerifyResetPasswordOTP)
	v1.POST("/auth/password-reset", authHandler.ResetPassword)
	v1.POST("/users", userHandler.Register)

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// listen and serve on 0.0.0.0:8080
	r.Run()
}
