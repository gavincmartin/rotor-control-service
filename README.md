# rotor-control-service

This rotor control service allows for remote control of a rotor system via a RESTful interface for the purpose of tracking satellites across the sky.

The software is written in Golang and runs as a single Linux binary with MongoDB.

## Features
- Manual rotor control
- Scheduling of future tracking passes
- Automatic execution of previously scheduled passes
- Complex querying of past and future passes
- (Coming Soon) Slack messaging for daily pass schedules and imminent passes
- (Coming Soon) Web interface for pass schedules


## Running Locally

You have two options for running the service on your local computer:
1) Running with Docker-Compose (Recommended)
2) Running the binary + your own MongoDB server


### Running with Docker Compose
1) Clone this repository (`git clone https://github.com/gavincmartin/rotor-control-service.git`)
2) Install [Docker](https://www.docker.com/products/docker-desktop)
3) Run the service with Compose (`docker-compose up`)

The service should be running at localhost on port 8080.


### Running the Binary + Your Own MongoDB Server
1) Clone this repository (`git clone https://github.com/gavincmartin/rotor-control-service.git`)
2) Install [MongoDB](https://docs.mongodb.com/manual/installation/)
3) Run MongoDB on your local machine (see the installation docs above for the guide for your OS)
4) Run the service (`go run service.go`)

The service should be running at localhost on port 8080.


## API Documentation
Coming soon...
