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
