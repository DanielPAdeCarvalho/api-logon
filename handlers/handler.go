package handlers

import (
	"login-app/database/query"
	"login-app/encrypt"
	"login-app/model"
	"login-app/utils"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func MarshalUser(c *gin.Context, dynamoClient *dynamodb.Client, logs utils.Loggar) (model.User, model.User) {
	var userModel model.User
	err := c.BindJSON(&userModel)
	utils.Check(err, logs)
	user := query.SelectUser(userModel.Nome, userModel.Senha, *dynamoClient, logs)
	return userModel, user
}

// this funcion going into AWS, than user and passaword not is hash
func GetUser(c *gin.Context, dynamoClient *dynamodb.Client, logs utils.Loggar) {
	userModel, user := MarshalUser(c, dynamoClient, logs)
	if user.Nome == "" {
		c.IndentedJSON(http.StatusNotFound, "Nome de usuário "+userModel.Nome+" não encontrado")
		return
	}
	if encrypt.CheckHash(userModel.Senha, user.Senha, logs) {
		c.IndentedJSON(http.StatusAccepted, "Authorized")
	} else {
		c.IndentedJSON(http.StatusUnauthorized, "Senha do usuário "+userModel.Nome+" não confere")
		return
	}
}

// GetClient receives a GET request with a client's email and password, checks if the email exists in the database and if the password is correct.
// It takes a gin.Context, a dynamodb.Client, a configuration.GoAppTools, a string and another string as parameters.
// Returns nothing.
func GetClient(c *gin.Context, dynamoClient *dynamodb.Client, logs utils.Loggar, email string, password string) {
	client := query.SelectClient(email, password, *dynamoClient, logs)
	if client.Email == "" {
		c.IndentedJSON(http.StatusNotFound, "Email "+email+" não encontrado")
		return
	}
	if encrypt.CheckHash(password, client.Senha, logs) {
		c.IndentedJSON(http.StatusAccepted, "Authorized")
	} else {
		c.IndentedJSON(http.StatusUnauthorized, "Senha do usuário "+client.Nome+" não confere")
		return
	}
}

func ResponseOK(c *gin.Context, logs utils.Loggar) {
	c.IndentedJSON(http.StatusOK, "Servidor up")
}

func PostUser(c *gin.Context, dynamoClient *dynamodb.Client, logs utils.Loggar) {
	var newUser model.User

	//configue o model punh with the retorn of context gin
	err := c.BindJSON(&newUser)
	//faz a chacagem de errode forma unificada
	utils.Check(err, logs)
	//call the package to encrypt password
	newUser.Senha = encrypt.EncrytpHash(newUser.Senha, logs)
	//calling the quiry package to mount the request for a DB
	query.InsertUser(dynamoClient, newUser, logs)
	name := ("Colaborador " + newUser.Nome + " criado com sucesso!")
	c.IndentedJSON(http.StatusCreated, (name))
}

// PostClient receives a POST request with a new client's information, encrypts the password and inserts it into the database.
// It takes a gin.Context, a dynamodb.Client and a loggar as parameters.
// Returns nothing.
func PostClient(c *gin.Context, dynamoClient *dynamodb.Client, logs utils.Loggar) {
	var newClient model.Client

	//configue o model punh with the retorn of context gin
	err := c.BindJSON(&newClient)
	//faz a chacagem de errode forma unificada
	utils.Check(err, logs)
	//call the package to encrypt password
	newClient.Senha = encrypt.EncrytpHash(newClient.Senha, logs)
	//calling the quiry package to mount the request for a DB
	query.InsertClient(dynamoClient, newClient, logs)
	name := ("Cliente " + newClient.Nome + " resgistrado com sucesso!")
	c.IndentedJSON(http.StatusCreated, name)
}
