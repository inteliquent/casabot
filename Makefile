default: builddocker

setup:
	go get github.com/nlopes/slack
	go get github.com/inteliquent/casatunes
	go get github.com/inteliquent/casabot

buildgo:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o /casabot github.com/inteliquent/casabot

builddocker:
	docker build -t inteliquent/casabot -f ./Dockerfile.build .
	docker run casabot /bin/true
	docker cp `docker ps -q -n=1`:/casabot .
	docker rm -f `docker ps -q -n=1`
	chmod 755 ./casabot
	docker build --rm=true -t inteliquent/casabot -f Dockerfile.static .
