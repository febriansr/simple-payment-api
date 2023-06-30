# Simple Payment API

Simple Payment API is a RESTful API application designed to provide various functionalities related to user authentication, payment processing, and session management. The API enables users to perform actions such as login, logout, and make a payment.

## Table of contents
* [Technologies](#technologies)
* [Setup](#setup)
* [Features](#features)
    * [Login](#login)
    * [Logout](#logout)
    * [Payment](#payment)

## Technologies
This project is built using the following technologies:

- Go Lang
- GIN Framework
- Redis

## Setup
To run this project locally, follow these steps:

1. Clone the repository or download the source code.
```
git clone https://github.com/febriansr/simple-payment-api
```
2. Go to the repository.
```
cd simple-payment-api
```
3. Start the Redis server. Refer to the Redis documentation for instructions on starting the server based on your operating system [here](#https://redis.io/docs/getting-started/).
```
sudo service redis-server start
```
4. Set up the environment by editing the .env
```
SERVER_PORT=[ServerPort]
SERVER_HOST=[ServerHost]
JSON_FILE_NAME_CUSTOMER=./data/customer.json
JSON_FILE_NAME_MERCHANT=./data/merchant.json
JSON_FILE_NAME_HISTORY=./data/history.json
ACCESS_TOKEN_LIFETIME=[AccessTokenLifetimeinMinutes]
APPLICATION_NAME=[ApplicationName]
JWT_SIGNATURE_KEY=[SignatureKey]
REDDIS_ADDRESS=[RedisHost]:[RedisPort]
REDDIS_PASSWORD=[RedisPassword]
```
5. Run the project.
```
go run main.go
```
6. To Run the tests you can simply use these commands.
```
go test -v ./... -coverprofile=cover.out  && go tool cover -html=cover.out
```
Make sure you have installed Go Lang and Redis on your machine before running the project.

## Features

### Login 
To use the application, you need to login by sending a POST request to the following endpoint:
```
http://[ServerHost]:[ServerPort]/v1/login/
```

Include the following JSON request format in the request body:
```
{
    "username": [username],
    "password": [password]
}
```
The password field in the JSON request body should be a plaintext version of the hashed password saved in the customer JSON file. If the request is successful, you will receive the following response:
```
{
    "code": 200,
    "message": "Success",
    "data": {
    "token" : [token] 
    }
}
```
If the server is unable to process your request, you will receive an error response with the appropriate error code and message.

### Payment
To make a payment, send a POST request to the following endpoint:
```
http://[ServerHost]:[ServerPort]/v1/payment/
```
Include the access token in the Authorization header of the request and provide the necessary payment details in the request body using the following format:
```
{
    "merchant_code": [merchant code],
    "amount": [amount]
}
```
The amount inputted should be less than or equal to the customer's balance and greater than 0. The token in Authorization should be valid and not expired. The transaction can only be made by registered users to registered merchants. A registered user cannot make a payment for another registered user without changing the token.
If the payment request is successful, you will receive a success response. If there is an error, you will receive an appropriate error response.

### Logout
To logout from the application, send a POST request to the following endpoint:
```
http://[ServerHost]:[ServerPort]/v1/logout/
```
Include the access token in the Authorization header of the request. If the logout request is successful, you will receive a success response. If there is an error, you will receive an appropriate error response.
If you already logged out, you have to login again to access the application.