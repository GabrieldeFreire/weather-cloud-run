include .env
$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

COMPOSE_RUN := docker compose up --build -V
COMPOSE_DOWN := docker compose down --volumes
PROJECT_NAME := weather-app
VERSION := latest
IMAGE := gcr.io/$(GCP_PROJECT_NAME)/$(PROJECT_NAME):$(VERSION)

run:
	bash -c "trap '$(COMPOSE_DOWN)' EXIT; $(COMPOSE_RUN)"

test:
	go test -v

build-prod:
	@docker build -t $(IMAGE) .
	@docker push $(IMAGE)
# @docker pull $(IMAGE)

deploy: build-prod
	gcloud run deploy $(PROJECT_NAME) --image $(IMAGE) --project $(GCP_PROJECT_NAME) --allow-unauthenticated --platform managed --region us-east1 --set-env-vars "WEATHER_API_KEY=$(WEATHER_API_KEY)"