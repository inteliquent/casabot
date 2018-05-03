default: builddocker-amd64

setup:
	go get -f -u github.com/nlopes/slack
	go get -f -u github.com/inteliquent/casatunes
	go get -f -u github.com/inteliquent/casabot

buildgo-arm:
	CGO_ENABLED=0 GOOS=linux GOARM=7 go build -ldflags "-s" -a -installsuffix cgo -o /casabot github.com/inteliquent/casabot

buildgo-amd64:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o /casabot github.com/inteliquent/casabot

builddocker-amd64:
	docker build --pull -t inteliquent/casabot -f ./Dockerfile.build --build-arg GOARCH=amd64 .
	docker run inteliquent/casabot /bin/true
	docker cp `docker ps -q -n=1`:/casabot .
	docker rm -f `docker ps -q -n=1`
	chmod 755 ./casabot
	docker build --pull --rm=true -t layered/casabot -f Dockerfile.amd64-static .

builddocker-arm:
	docker build --pull -t inteliquent/casabot -f ./Dockerfile.build --build-arg GOARCH=arm .
	docker run inteliquent/casabot /bin/true
	docker cp `docker ps -q -n=1`:/casabot .
	docker rm -f `docker ps -q -n=1`
	chmod 755 ./casabot
	docker build --pull --rm=true -t layered/arm-casabot -f Dockerfile.arm-static .
