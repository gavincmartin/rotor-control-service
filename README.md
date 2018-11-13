# rotor-control-service

This rotor control service allows for remote control of a rotor system via a RESTful interface for the purpose of tracking satellites across the sky.

The software is written in Golang and runs as a single Linux binary with MongoDB.

## Features
- Manual rotor control
- Scheduling of future tracking passes
- Automatic execution of previously scheduled passes
- Complex querying of past and future passes
- Slack notifications for daily pass schedules and imminent passes
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

## Configuration Options
The service can be given a custom configuration by altering the environment variables in the [docker-compose.yml](https://github.com/gavincmartin/rotor-control-service/blob/master/docker-compose.yml) file (if running with Compose) or by manually setting environment variables prior to running it (if running locally).

The configuration options are:
1) `PORT`: the port on which the service listens
2) `MONGO_SERVER`: the hostname of the MongoDB server you want to connect to. This should match the name of the MongoDB container you are running with Compose (`mongodb` by default), or can connect to your local MongoDB server at `localhost` or another specified address.
3) `MONGO_DB_NAME`: the name of the tracking pass database used for this application (`tracking_passes_db` by default).
4) `SLACK_POST_URL`: the URL of your Slack webhook that should receive POST requests from the service. This is where daily schedules and "pass starting" notifications will be sent. For information on configuring this for your Slack workspace, you can look at [Slack's API Documentation on Incoming Webhooks](https://api.slack.com/incoming-webhooks).
5) `SLACK_SCHEDULE_POST_TIME`: the time that you'd like a daily schedule sent to the `SLACK_POST_URL` specified above. It should be in `HH:MM` format. It will schedule based upon what the local timezone of the machine running it is--if you are running in a container with Compose, you should supply the desired time in UTC.

## API Documentation
Postman-generated documentation with example requests can be found [here](https://documenter.getpostman.com/view/5438849/RzZAkdf5).
