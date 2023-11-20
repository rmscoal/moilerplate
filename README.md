# Moilerplate
> A monolith boilerplate for Golang backend applications with built-in strong security in mind.

> [!WARNING]<br>
> In development

## What is this project
It is a boilerplate for monolithic backend application that prioritizes security. I created this project to serves as boilerplates for my other backend applications. One of my examples project that uses this boilerplate is [SinarLog's backend](https://github.com/SinarLog/backend). It also follows Uncle Bob's Clean Architecture concepts and is inspired by some of the best clean architecture golang app out there.

## Folder Structure.
- `cmd` consists of bootstraping the app as well as starting the server.

- `config` loading the applications config by reading from .env files.

- `pkg` consists of all external/in-house packages for the application to use. Usually consists of the infrastructure or service initializations.

- `testdata` consists of mocks structs for testing.

- `internal` where all the fun begins<br>
  - `internal/domain` stores the domain of the app. I'm trying to follow Domain Driven Design as much as possible here.<br>
  - `internal/delivery` consists of the delivery methods to communicate, like the http endpoints, middlewares, and routers.<br>
  - `internal/utils` consists of application's utility functions, like primitive type manipulations.<br>
  - `internal/app` stores the application layer.<br>
    - `internal/app/usecase` consists of the application logic and orchestration of its infrastructure and service layer.<br>
    - `internal/app/repo` consists of interfaces that the infrastructure layer has to follow.<br>
    - `internal/app/service` consistes of interfaces that the service layer has to follow.<br>
  - `internal/adapter` consists of implementations to fullfil the application's infrastructure and service contracts.<br>
  - `internal/composer` acts as the manager to store all usecase, infrastructure and service layer for easier management.<br>

Currently in development
