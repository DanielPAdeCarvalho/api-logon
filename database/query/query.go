package query

import (
	"context"
	"login-app/model"
	"login-app/utils"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/beevik/ntp"
)

func InsertUser(dynamoClient *dynamodb.Client, user model.User, logs utils.Loggar) {

	//o codigo esta indo no observatorio nacional pegar a data e hora
	datatemp, err := ntp.Time("a.st1.ntp.br")
	utils.Check(err, logs)

	//formatando a data retornada do observatorio para a data no formato desejado (yy-mm-dd_hh:mm)
	tempY := datatemp.Format("06")
	tempM := datatemp.Format("01")
	tempD := datatemp.Format("02")
	tempH := strconv.Itoa(datatemp.Hour())
	tempMin := strconv.Itoa(datatemp.Minute())
	user.DataCriacao = tempY + "-" + tempM + "-" + tempD + "_" + tempH + ":" + tempMin

	//converter a struct em um json
	item, err := attributevalue.MarshalMap(user)
	utils.Check(err, logs)

	_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("LoginColaborador"),
		Item:      item,
	})
	utils.Check(err, logs)
}

func SelectUser(Nome string, Senha string, dynamoClient dynamodb.Client, logs utils.Loggar) model.User {

	query := expression.Name("Nome").Equal(expression.Value(Nome))
	proj := expression.NamesList(expression.Name("Nome"), expression.Name("Senha"))
	expr, err := expression.NewBuilder().WithFilter(query).WithProjection(proj).Build()
	utils.Check(err, logs)

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("LoginColaborador"),
	}

	// Make the DynamoDB Query API call
	result, err := dynamoClient.Scan(context.TODO(), params)
	utils.Check(err, logs)

	var user model.User
	for _, i := range result.Items {
		item := model.User{}
		err = attributevalue.UnmarshalMap(i, &item)
		utils.Check(err, logs)
		user = item
	}

	return user
}

func InsertClient(dynamoClient *dynamodb.Client, cliente model.Client, logs utils.Loggar) {

	// Get current date and time from observatory
	datatemp, err := ntp.Time("a.st1.ntp.br")
	utils.Check(err, logs)

	// Format the date and time as "yy-mm-dd_hh:mm"
	cliente.DataCriacao = datatemp.Format("06-01-02_15:04")

	//converter a struct em um json
	item, err := attributevalue.MarshalMap(cliente)
	utils.Check(err, logs)

	_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("LoginCliente"),
		Item:      item,
	})
	utils.Check(err, logs)
}

func SelectClient(Email string, Senha string, dynamoClient dynamodb.Client, logs utils.Loggar) model.Client {

	query := expression.Name("Email").Equal(expression.Value(Email))
	proj := expression.NamesList(expression.Name("Email"), expression.Name("Senha"))
	expr, err := expression.NewBuilder().WithFilter(query).WithProjection(proj).Build()
	utils.Check(err, logs)

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("LoginCliente"),
	}

	// Make the DynamoDB Query API call
	result, err := dynamoClient.Scan(context.TODO(), params)
	utils.Check(err, logs)

	var cliente model.Client
	for _, i := range result.Items {
		item := model.Client{}
		err = attributevalue.UnmarshalMap(i, &item)
		utils.Check(err, logs)
		cliente = item
	}

	return cliente
}
