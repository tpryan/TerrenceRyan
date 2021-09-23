BASEDIR = $(shell pwd)
APPNAME=color
ZONE=us-central1-a
REGION=us-central1
PROJECT = terrenceryan-com

env:
	@echo "Making sure project settings all in order"
	@gcloud config set project $(PROJECT)
	@gcloud config set compute/zone $(ZONE)
	@gcloud config set run/region $(ZONE)


deploy: env
	-gcloud functions deploy subscribeMailgun --trigger-topic cloud-builds --runtime nodejs10

build:
	gcloud builds submit --config cloudbuild.yaml .