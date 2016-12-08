# colinmc
Simple slack bot that posts stock prices. It can be build as an RTM slack bot, or a slash command bot that loads into aws lambda

## Build
`go build`
or for aws
`go build -tags AWS`

## Run
`./colinmc [slack-bot-token]`
or for aws
`./aws/zipit.sh`
and upload the stocktopus.zip to lambda

## Usage
The RTM bot will look for any direct messages sent to it and try to pase them as tickers, and respond with stock quotes.
> @stockbotname GOOGL

The aws slash command will respond to slash commands. Single tickers will be a quote and inline graph. 
> /stockbotname GOOGL




for a complete list of commands the bot supports.
> /stockbotname help 

