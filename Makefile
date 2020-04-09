BASEDIR = $(shell pwd)
APPNAME = terrenceryan


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

build:
	gcloud builds submit --config cloudbuild.yaml .