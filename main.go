package main

import (
	"log"
	"login-app/driver"
	"login-app/handlers"
	"login-app/utils"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var (
	dynamoClient *dynamodb.Client
	logs         utils.Loggar
)

func inLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func setupRouter() *gin.Engine {

	appRouter := gin.New()

	// TSA
	appRouter.OPTIONS("/logon", func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.JSON(200, nil)
	})

	appRouter.GET("/", func(ctx *gin.Context) {
		logs.InfoLogger.Println("Servidor Ok")
		handlers.ResponseOK(ctx, logs)
	})

	appRouter.POST("/logon", func(ctx *gin.Context) {
		handlers.GetUser(ctx, dynamoClient, logs)
	})

	appRouter.POST("/logonclient", func(ctx *gin.Context) {
		var credentials struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&credentials); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		handlers.GetClient(ctx, dynamoClient, logs, credentials.Email, credentials.Password)
	})

	appRouter.POST("/signin", func(ctx *gin.Context) {
		handlers.PostUser(ctx, dynamoClient, logs)
	})

	appRouter.POST("/signclient", func(ctx *gin.Context) {
		handlers.PostClient(ctx, dynamoClient, logs)
	})

	return appRouter
}

// Para compilar o binario do sistema usamos:
//
//	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o api-login .
//
// para criar o zip do projeto comando:
//
// zip lambda.zip api-login
//
// main.go
func main() {
	InfoLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)

	logs.InfoLogger = InfoLogger
	logs.ErrorLogger = ErrorLogger
	var err error
	// chamada de função para a criação da sessao de login com o banco
	dynamoClient, err = driver.ConfigAws()
	//chamada da função para revificar o erro retornado
	utils.Check(err, logs)

	if inLambda() {

		log.Fatal(gateway.ListenAndServe(":8080", setupRouter()))
	} else {

		log.Fatal(http.ListenAndServe(":8080", setupRouter()))
	}
}
