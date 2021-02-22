BASEDIR = $(shell pwd)
APPNAME = terrenceryan
PROJECT = terrenceryan-com
REDISNAME=terrenceryan-com-cache
REGION=us-central1
PROJECTNUMBER=$(shell gcloud projects list --filter="$(PROJECT)" \
			--format="value(PROJECT_NUMBER)")
REDISIP=$(shell gcloud redis instances describe $(REDISNAME) --region us-central1 --format='value(host)' )			

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

main:
	cd prod && GOOS=linux CGO_ENABLED=0 go build -o main main.go 	

image: prod angular 
	docker build -t $(APPNAME) "$(BASEDIR)/prod/."
	
serve: 
	docker run --name=$(APPNAME) -d -P -p 8080:8080 $(APPNAME)


build: env
	gcloud builds submit --config cloudbuild.yaml .

services: env
	-gcloud services enable cloudbuild.googleapis.com
	-gcloud services enable cloudfunctions.googleapis.com
	-gcloud services enable redis.googleapis.com
	-gcloud services enable vpcaccess.googleapis.com
	-gcloud services enable secretmanager.googleapis.com 

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
	-gcloud projects add-iam-policy-binding $(PROJECT) \
  	--member serviceAccount:$(PROJECTNUMBER)@cloudbuild.gserviceaccount.com \
  	--role roles/secretmanager.secretAccessor	 

function: env
	gcloud services enable cloudfunctions.googleapis.com
	-gcloud functions deploy subscribeMailgun --trigger-topic cloud-builds \
	--runtime nodejs10 --set-env-vars GCLOUD_PROJECT=$(PROJECT) \
	--source $(BASEDIR)/functions/email	--allow-unauthenticated	      

memorystore: env
	-gcloud redis instances create $(REDISNAME) --size=1 --region=$(REGION)
	-gcloud compute networks vpc-access connectors create \
	redisconnector --network default --region $(REGION) \
	--range 10.8.0.0/28		


redis: redisclean
	docker run --name some-redis -p 6379:6379 -d redis	

redisclean:
	-docker stop some-redis
	-docker rm some-redis	


dev: redis
	(trap 'kill 0' SIGINT; \
	cd $(BASEDIR)/prod && \
	export REDISHOST=127.0.0.1 && \
	export REDISPORT=6379 && \
	go run main.go & \
	cd $(BASEDIR)/frontend && ng serve --open )		

backend: redis
	cd $(BASEDIR)/prod && \
	export REDISHOST=127.0.0.1 && \
	export REDISPORT=6379 && \
	go run main.go 		

secretredisip:
	printf $(REDISIP) | gcloud secrets create redisip --data-file=-
