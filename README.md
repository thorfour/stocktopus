[![Go Report Card](https://goreportcard.com/badge/github.com/thorfour/stocktopus)](https://goreportcard.com/report/github.com/thorfour/stocktopus)
[![WTFPL licensed](https://img.shields.io/badge/license-WTFPL-blue.svg)](https://github.com/thorfour/stocktopus/blob/master/LICENSE)
[![CircleCI](https://circleci.com/gh/thorfour/stocktopus.svg?style=svg)](https://circleci.com/gh/thorfour/stocktopus)
[![Docker Repository on Quay](https://quay.io/repository/thorfour/stocktopus/status "Docker Repository on Quay")](https://quay.io/repository/thorfour/stocktopus)

<a href="https://slack.com/oauth/authorize?scope=commands&client_id=15348769670.121517816146"><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>

<img src="stocktopus_cropped.png" height="285" width="132"/>

# stocktopus
Simple Slack bot that posts stock prices, manages play money portfolios and watch lists. 

## Download

`go get -t -d github.com/thorfour/stocktopus/...`

#### or
`docker pull quay.io/thorfour/stocktopus:v1.0.0`

## Deploy

If you'd like to deploy your own version of stocktopus to DigitalOcean cloud there are [terraform](https://www.terraform.io/) files provided in `deploy/terraform`, simply run `terraform apply`. You'll need to provide a do api key and provide your own hostname. 

## Build

make

### Docker:

make docker

### Binaries should be output to bin/ directories.

## Run
`docker run -d -p 80:80 -p 443:443 -e REDISADDR=<redis endpoint> -e REDISPW=<redis password> quay.io/thorfour/stocktopus:v1.0.0`

## Usage
The slash command will respond to slash commands. Single tickers will be a quote and inline graph. 
> /stocktopus GOOGL

for a complete list of commands the bot supports.
> /stocktopus help 

In addition to what is covered in the help menu, stocktopus also supports team-wide watchlists. To utilize these you use the same format as you would for your personal watch list but simple add a name after a `#` character immediately after the command.

For example:
`/stocktopus watch #funlist GOOG`
Would add GOOG to a watch list called funlist, and then anyone in your Slack team can access that same list Ex. `/stocktopus list #funlist`

You might also want to watch/buy or lookup securities listed on non-US exchanges. To do so simply add the exchange followed by a colon(:) before the ticker name. 
Ex. `/stocktopus tse:are` to list the ARE stock from the Toronto Stock Exchange 

### Note: other exchanges are no longer supported by the application. To enable this feature on a local build, you can build with alphavantage instead of IEX using `make docker-alpha`

## Attribution

Data provided for free by [IEX](https://iextrading.com/developer/)

## Privacy Policy

Stocktopus does not collect any personal identifying information. It does not store a history of your requests. The only data it does store is a unique ID received from Slack for a user if they opt to use the list or play money features. If you use the list or play money features it also stores the list of stocks a user has bought or added to their watch list. 

## Contact

If you have questions comments or concerns about the app please email help@stocktopus.io

#### Special Thanks to /u/shirokarasu over at r/DrawForMe for creating the icon
