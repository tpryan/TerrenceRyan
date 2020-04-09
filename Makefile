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

	