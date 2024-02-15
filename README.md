<h1 align="center">Welcome Simple Bank API GO Clean Architecture</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-1.0.0-blue.svg?cacheSeconds=2592000" />
  <a href="https://goreportcard.com/badge/github.com/benebobaa/amikom-bri-api" target="_blank">
    <img alt="Build" src="https://goreportcard.com/badge/github.com/benebobaa/amikom-bri-api" />
  </a>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>

</p>

- This project is a user-friendly solution for basic banking tasks, built using the Go programming language and adhering to the Clean Architecture principles. It encompasses essential features such as account creation, listing, balance checking, account transfers, and transaction history tracking. Additionally, it incorporates advanced functionalities like user registration, email verification, authentication with JWT (JSON Web Token) using access and refresh tokens, password recovery with email-based reset links, email notifications for transfers, and the ability to create and export transaction and expense history to PDF.
- Tech Stack Used: Go, PostgreSQL, Onesignal ,GORM, JWT, Validator, Viper and Fiber.
- Infrastructure: Docker, AWS EC2, AWS Sec Manager, AWS RDS, Github Actions, and Cloudflare.
## Architecture
-  I have crafted this architecture inspired by Uncle Bob's Clean Architecture, the Golang course from Programmer Zaman Now, Harisenin's Proclass Golang, GSabadini's Repository, and also my own Created Go Native Clean Architecture. I have customized it to fit my specific needs. Source above this:
- [Amikom Pedia API](https://github.com/benebobaa/amikom-pedia-api).
- [GSabadini Clean Arch](https://github.com/GSabadini/go-clean-architecture).
- [Progammer Zaman Now(Eko Khannedy) Clean Arch](https://github.com/khannedy/golang-clean-architecture).

![Clean Architecture](clean.png)

## ERD Diagram Database
### Source : [diagram.io erd](https://dbdiagram.io/d/65ccfe1dac844320ae26a809)
![ERD Diagram](erd.png)

## Requirements/dependencies
- Docker


## Getting Started Run Locally

- Pull Postgres Images and Run Container with Docker

```sh
make postgres
```

- Run Database Migrations

```sh
make migrateup
```

- Run Go Server API

```sh
go run main.go
```


## API Request

### Currently Published Example Endpoints
### Base URL: `https://go.amikompedia.xyz`

Note: I kindly request you to use or try this API using a valid email, as the push send email feature can function properly.

- Or you can use this user, only these account has provided sample data

| Email                   |  Password   | PIN     |
|-------------------------|:-----------:|:--------|
| `manaada@gmail.com`     | `apaiya123` | `123456` |
| `bensatriya3@gmail.com` | `apaiya123` |      `123456`   |



## Auth
| Endpoint                                 | HTTP Method |       Description        |
|------------------------------------------|:-----------:|:------------------------:|
| `/api/v1/auth/_register`                       |   `POST`    | `Register user ` |
| `/api/v1/auth/_login`                        |   `POST`    |  `Login user `   |
| `/api/v1/auth/_forgot-password` |   `POST`    |  `Forgot password user`  |

## Users

| Endpoint                                 | HTTP Method |        Description        |
|------------------------------------------|:-----------:|:-------------------------:|
| `/api/v1/users/profile`                        |    `GET`    | `Get detail profile user` |
| `/api/v1/users`                        |    `GET`    |    `List user account`    |
| `/api/v1/users/:username` |  `DELETE`   |   `Delete user account`   |

## Transaction

| Endpoint                                 | HTTP Method |                  Description                  |
|------------------------------------------|:-----------:|:---------------------------------------------:|
| `/api/v1/transfer`                        |   `POST`    |     `Create transfer to another account`      |
| `/api/v1/entries/filter`                        |    `GET`    | `Get transaction history and exported to pdf` |
| `/api/v1/expenses-plan` |   `POST`    |            `Create expenses plan`             |

### For more detail api documentation please import `apispec.json` in this repo to your Postman

