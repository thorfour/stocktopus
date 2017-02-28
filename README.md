[![Go Report Card](https://goreportcard.com/badge/github.com/thourfor/stocktopus)](https://goreportcard.com/report/github.com/thourfor/stocktopus)
[![WTFPL licensed](https://img.shields.io/badge/license-WTFPL-blue.svg)]
(https://github.com/thourfor/stocktopus/blob/master/LICENSE)

<a href="https://slack.com/oauth/authorize?scope=commands&client_id=15348769670.121517816146"><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>

#stocktopus
Simple Slack bot that posts stock prices. It can be build as an RTM Slack bot, or a slash command bot that loads into aws lambda

## Build
`go build -tags RTM`
or for aws
`go build`

## Run
`./stocktopus [slack-bot-token]`
or for aws
`./aws/zipit.sh`
and upload the stocktopus.zip to lambda

## Usage
The RTM bot will look for any direct messages sent to it and try to pase them as tickers, and respond with stock quotes.
> @stockbotname GOOGL

The aws slash command will respond to slash commands. Single tickers will be a quote and inline graph. 
> /stocktopus GOOGL




for a complete list of commands the bot supports.
> /stocktopus help 

