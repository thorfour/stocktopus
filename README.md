[![Go Report Card](https://goreportcard.com/badge/github.com/thourfor/stocktopus)](https://goreportcard.com/report/github.com/thourfor/stocktopus)
[![WTFPL licensed](https://img.shields.io/badge/license-WTFPL-blue.svg)](https://github.com/thourfor/stocktopus/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/thourfor/stocktopus.svg?branch=master)](https://travis-ci.org/thourfor/stocktopus)

<a href="https://slack.com/oauth/authorize?scope=commands&client_id=15348769670.121517816146"><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>

<img src="stocktopus_cropped.png" height="285" width="132"/>

# stocktopus
Simple Slack bot that posts stock prices. It can be build as an RTM Slack bot, or a slash command bot that loads into aws lambda

## Download

`go get github.com/thourfor/stocktopus`

## Build
### AWS Lambda:

`make aws`

### GCP Cloud Function (under development):

`make gcp`

### RTM CLient (no longer actively developed):

`make rtm`

### Files should be output to bin/ directories.

### Serverless:

`bin/aws or bin/gcp` which will contain the binary and the zip file of the nodejs handler and binary.

### Rtm:

`bin/rtm` for the rtm client

## Run
`./stocktopus [slack-bot-token]`
or for aws
upload the `stocktopus.zip` file in `/bin/aws` as a lambda function
or for gcp
upload the `stocktopus.zip` file in `bin/gcp` as a cloud function

## Usage
The RTM bot will look for any direct messages sent to it and try to pase them as tickers, and respond with stock quotes.
> @stockbotname GOOGL

The aws slash command will respond to slash commands. Single tickers will be a quote and inline graph. 
> /stocktopus GOOGL




for a complete list of commands the bot supports.
> /stocktopus help 

In addition to what is covered in the help menu, stocktopus also supports team-wide watchlists. To utilize these you use the same format as you would for your personal watch list but simple add a name after a `#` character immediately after the command.

For example:
`/stocktopus watch #funlist GOOG`
Would add GOOG to a watch list called funlist, and then anyone in your Slack team can access that same list Ex. `/stocktopus list #funlist`

You might also want to watch/buy or lookup securities listed on non-US exchanges. To do so simply add the exchange followed by a colon(:) before the ticker name. 
Ex. `/stocktopus tse:are` to list the ARE stock from the Toronto Stock Exchange 

## Attribution

Data provided for free by [IEX](https://iextrading.com/developer/)

## Privacy Policy

Stocktopus does not collect any personal identifying information. It does not store a history of your requests. The only data it does store is a unique ID received from Slack for a user if they opt to use the list or play money features. If you use the list or play money features it also stores the list of stocks a user has bought or added to their watch list. 

## Contact

If you have questions comments or concerns about the app please email thorvald@protonmail.ch

## Donations

If you like the app and would like to donate to help us pay for server costs we accept Bitcoin `18Q4Nnyc3AxrHis5ATioWSZ6hn5SqyzVyj`

#### Special Thanks to /u/shirokarasu over at r/DrawForMe for creating the icon
