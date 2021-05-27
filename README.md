# UnicornKitchen
An unusual recipe book built with a rainbow of modern tech!

The purpose of this project is to explore new programming and data languages that I don't use at my day job or in other personal projects.  All good devs need to keep up to date!

Currently the application is storing data in a postgres database and then acts as an HTTP server so that a user may interact with it via GraphQL queries.  In the future, a flask application will be included that will place a proper UI on top of the GraphQL layer.


Here is an up to date list of the different technologies used within:

+ Go
+ Postgres
+ GraphQL
+ Docker
+ Shell Scripts
+ AWS
+ Azure

This project is setup so that it can run on your local machine, either by go-run or in a docker container.  For local instances, a shell script in included that will run Postgres in a container, or if you already have it installed you may set the appropriate environment variables.

The included docker file will build an image that can be run on either AWS (I personally use ECS) or Azure.
