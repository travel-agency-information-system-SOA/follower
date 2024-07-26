# Travel agency information system (service-oriented architecture)  

This project represents an upgrade in software design, where we transitioned from a monolithic application we previously designed to a service-oriented architecture. 

You can find more information about the project's basis [here](https://github.com/travel-agency-information-system/back-end).

We divided the basic monolith into several sections that became separate applications, each with its own repository, distinct language, and database. 

We ran them together using Docker.

# Followers

OSNOVNO O OVOM MODULU

The Followers microservice was added to our application at a later stage (it was not part of the initial monolithic application).

It is also written in Go, and the back-end controllers invoke its methods from the Go file using gRPC.

We integrated everything using Docker containers and ran all services together. Docke enabled us to easily orchestrate and manage different services, improving the efficiency of development and testing.

## Transition to gRPC

The transition to gRPC represents a significant step in enhancing our microservice architecture. This remote procedure call (RPC) framework offers numerous advantages over traditional RESTful APIs, especially in complex systems with a large number of microservices.

gRPC uses Protocol Buffers (Protobuf) for data serialization, which allows for faster and more efficient encoding and decoding compared to JSON used in REST APIs. This significantly reduces message size and improves communication speed between services.

Services and methods are defined in .proto files, providing clear and consistent API documentation. This makes it easier to maintain and enhance communication between services.

The transition to gRPC has made our architecture more efficient, scalable, and easier to maintain, improving the overall performance and quality of communication between microservices.

## Observability (Logging, Tracing, Metrics)

To gain a better understanding of our system's functioning and to simplify error resolution, we introduced logging, tracing, and metrics (the three pillars of observability).

## Technologies

- ***Server platform***: Go (Golang) 

- ***Client platform***: Angular (TypeScript, HTML, CSS) with RESTful services for the front-end interface

- ***Database***: Graph Databse (Neo4j)

- ***Back-end***: C# (ASP.NET) which interfaces with the Go-based microservice, managing requests and orchestrating the overall application flow

## Getting started

To set up the project locally using Docker, follow these steps:

```
1. Clone the repository
2. Run the entire setup using Docker Compose
```
- Ensure Docker and Docker Compose are installed and running on your machine
- Use the provided docker-compose.yml file and docker-compose-migrations.yml file to manage the services
- To build and start all services:
```
docker-compose up --build
```
- To stop all services:
```
docker-compose down
```
Use appropriate tools like pgAdmin and MongoDB Compass to interact with the database.

## Configuration

Make sure to review the docker-compose.yml file and Dockerfiles for specific configuration details, such as environment variables, volume mounts, and network settings.

Check the port mappings in docker-compose.yml to ensure that services are properly exposed and accessible.

You can download the files from [this link](https://ufile.io/f/ud3nw). If you are unable to do so, please contact me via email.

## Authors
Contributors to this project:
- Ana Radovanović
- Kristina Zelić
- Milica Petrović
- Petar Kovačević
