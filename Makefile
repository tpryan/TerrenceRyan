BASEDIR = $(shell pwd)
APPNAME = terrenceryan
PROJECT = terrenceryan-com
PROJECTNUMBER=$(shell gcloud projects list --filter="$(PROJECT)" \
			--format="value(PROJECT_NUMBER)")

env:
	gcloud config set project $(PROJECT)

angular:
	cd frontend && ng build --prod


prod: clean angular

clean:
	-rm -rf prod/static
	-rm -rf prod/main
	-docker stop $(APPNAME)
	-docker rm $(APPNAME)
	-docker rmi $(APPNAME)		

image: prod
	docker build -t $(APPNAME) "$(BASEDIR)/prod/."
	
serve: 
	docker run --name=$(APPNAME) -d -P -p 8080:8080 $(APPNAME)

dev:
	(trap 'kill 0' SIGINT; \
	cd $(BASEDIR)/prod && go run main.go & \
	cd $(BASEDIR)/frontend && ng serve --open )	

build: env
	gcloud builds submit --config cloudbuild.yaml .

services: env
	-gcloud services enable cloudbuild.googleapis.com
	-gcloud services enable cloudfunctions.googleapis.com

perms: env
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
  	--role roles/cloudfunctions.admin
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECT)@appspot.gserviceaccount.com \
  	--role roles/cloudfunctions.admin  
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECT)@appspot.gserviceaccount.com \
  	--role roles/iam.serviceAccountUser
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
  	--role roles/iam.serviceAccountUser   	 

function: env
	gcloud services enable cloudfunctions.googleapis.com
	-gcloud functions deploy subscribeMailgun --trigger-topic cloud-builds \
	--runtime nodejs10 --set-env-vars GCLOUD_PROJECT=$(PROJECT) \
	--source $(BASEDIR)/functions/email	--allow-unauthenticated	      	