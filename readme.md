# MEWS API

[![Go](https://github.com/fahminlb33/devoria1-wtc-backend/actions/workflows/go.yml/badge.svg)](https://github.com/fahminlb33/devoria1-wtc-backend/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/fahminlb33/devoria1-wtc-backend/branch/master/graph/badge.svg?token=hRNbJKqQgM)](https://codecov.io/gh/fahminlb33/devoria1-wtc-backend)

> Great news comes from great backend service :wink:

This is my submission for DEVORIA WTC.

In this repository you'll find my preferred codebase style which are highly influenced by my other project, [TASI](https://github.com/fahminlb33/tasi-backend). This project structure is also used in one of my squad service (LGC-Visibility). I hope I can demonstrate a clear and concise clean architecture in functional way (I'm more of an OOP guy).

More explanation will come later on after I finish testing this project.

## What you'll get

It ain't much, but honest work :smile:

- MEWS API - User management API: register, login, profile
- MEWS API - Article management API: create, update, delete, find (with pagination), get single
- Integrated with Elastic APM
- Using PostgreSQL with GORM
- Docker support
- Unit tests
- Pretty much overkill coverage report lol

## Running this Project

For simplicity sake, you can run this projet in Docker. Clone this repo and then run `docker-compose up` at this project root directory. You'll get ELK already setup (Elasticsearch, Kibana, Elastic APM), PostgreSQL, and MEWS API all in single docker-compose.

If you prefer running this app locally, you'll need Go 1.17 and PostgreSQL running on your local machine. Rename the `.env.example` to `.env` and adjust the contents according to your environment. If you want to use APM, you'll have to cofigure it yourself too.

To run this project:

``` bash
# install swag to generate swagger
go install github.com/swaggo/swag/cmd/swag@latest

# if you're on PowerShell terminal,
./make swagger
./make run

# if you're on other terminal,
swag init
go run main.go
```

## Configuration

The `.env` file has all the required environment variables you'll need to run this project. One small note about private and public key pair for JWT verification, you'd have to generate a new key WITHOUT password and stored in PCKS1 format.

To generate new key, you can run `./make generate_keypair`. After you've got your keys, encode your keys as Base64 and then store it in the `.env` file. Why I don't place the key as files? Because it's not safe to store keys in your git. You'd be better off using Vault services or just place it in the environment variables.

Other than that, I think the settings is already self explanatory.

## Unit Tests and Coverage

This repo also includes many handy script to run and create coverage reports.

- `./make test`, run all unit tests
- `./make coverage`, run all unit tests and generate HTML coverage report
- `./make coverage_pretty`, run all unit test and generate HTML coverage report using Cobertura output and ReportGenerator exporter (nicer UI for coverage report)

Before you can run the `./make coverage_pretty` command, you'll need to install .NET 6 and `gocov` tool to create Cobertura XML. Here's how:

```bash
dotnet tool install -g dotnet-reportgenerator-globaltool
go install github.com/axw/gocov/gocov
go install github.com/AlekSi/gocov-xml
go install github.com/matm/gocov-html
```

Those tools will export the coverage as Cobertura XML and HTML based report. The XML file is used for ReportGenerator. You can then use this XML file as coverage report for SonarQube or Jenkins.
