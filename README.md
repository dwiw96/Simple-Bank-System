# Overview
This is a Go and postgresql REST API project to reengineering digital bank system in backend side.
This project is for my personal learning about backend and the project itself based on website ("https://dev.to/techschoolguru/series/7172") and ("https://github.com/dwiw96/Sea-Wave-Measurenment-Using-9dof").

## Supported Features
* User can create account and login to that account
* Every account can create wallet that belong the that account.
* Authentication and authorization using PASETO by earer token.
* Only authentication account that can access all cause including create wallet, transfer balance, edit information, etc..
* When performing transfer balance system will automatically create entry and transfer record.

## Installation Guide
* Clone this repository (https://github.com/dwiw96/Simple-Bank-System.git).
* The main branch is the most stable branch at any given time, ensure you're working from it.
* This project run database via docker, so ensure that you installed docker in your machine.

### To Run
For more information, see (https://github.com/dwiw96/Simple-Bank-System/blob/main/Makefile)
* To run the postgresql server inside docker, run this task:
```
make docckerStart
```
* Create database table use migrate library
```
make migrateUp
```
* To Login into database using sql
```
make dockerExec
```
* Run local server
Server is run at http://localhost:8080
```
make server
```

## Request and Response Example
### Create an account
* Request
```
curl -X POST "http://localhost:8080/account" \
     -H "Content-Type: application/json" \
     --data-binary @- << EOF
        {  "username": "marker",
            "fullname": "Mark Spekman",
            "date_of_birth": "1986-05-12",
            "address": {
                "province": "DKI Jakarta",
                "city": "Jakarta Barat",
                "zip": 10203,
                "street": "Jl Tj Karang 3-4 A Ged Dana Pensiun Bank Mandiri Room 302 Lt 3"
            },
            "email": "marker@mail.com",
            "password": "secret1"
        }
EOF
```
* Response
```
{"AccountNumber":1019194980,"Username":"marker","FullName":"Mark Spekman","DateOfBirth":"1986-05-12","Address":{"province":"DKI Jakarta","city":"Jakarta Barat","zip":10203,"street":"Jl Tj Karang 3-4 A Ged Dana Pensiun Bank Mandiri Room 302 Lt 3"},"Email":"marker@mail.com","PasswordChangeAt":"0001-01-01T07:07:12+07:07","CreatedAt":"2023-11-23T13:01:24.761992+07:00"}
```

## API Endpoint
For full complete api endpoint you can use openapi.yml. <br>
file: (https://github.com/dwiw96/Simple-Bank-System/blob/main/openapi.yml)

## Status Code
| Status Code | Description |
| :--- | :--- |
| 200 | OK |
| 400 | BAD REQUEST |
| 401 | UNAUTHORIZATE |
| 422 | UNPROCESSABLE ENTITY |
| 500 | INTERNAL SERVER ERROR |

## Technologies Used
* [Go Programming Language]
* [Postgresql] sql database for saving all data in this project
* [Docker] (https://www.docker.com/) Docker is a software platform that allows you to build, test, and deploy applications quickly. - amazon
* [PASETO] (https://paseto.io/) Paseto (Platform-Agnostic SEcurity TOkens) is a specification and reference implementation for secure stateless tokens.
