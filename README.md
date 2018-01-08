# CasaBot
Bot that can control the Casaplayer music device from Slack

Utilizes the [nlopes/slack](https://github.com/nlopes/slack) and [inteliquent/casatunes](https://github.com/inteliquent/casatunes) Go libraries.
# Docker
To create the inteliquent/casabot docker image, run the Makefile with:
```
make
```
The Docker image can then be run with:
```bash
export SLACK_TOKEN='xxx'
export CASA_ENDPOINT='yyy'

docker run -d -e SLACK_TOKEN -e CASA_ENDPOINT casabot
```

## Docker ARM7 build
For Raspberry Pi (3), running the commands below will compile a docker image compatible with the Pi's ARM7 architecture:
```bash
make builddocker-arm
```
This will create a docker image called `inteliquent/arm-casabot`. `docker save` and `docker load` can then be used to port the image to a Raspberry Pi.

## ssh into machine
From the office:
`ssh pi@raspberrypi.local`
