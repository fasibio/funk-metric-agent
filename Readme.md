# Funk metric Agent

This agent will be used to get Host metrics (CPU / Memory/ Disk / Host livetime)

If you don´t know what funk is please read [here](https://github.com/fasibio/funk-server) first. 

[![coverage report](https://gitlab.com/fasibio/funk-metric-agent/badges/master/coverage.svg)](https://sonar.server2.fasibio.de/dashboard?id=fasibio_metric_master)

You can configure each Agent different on each installation and combine the information to one Elasticsearch Stack.


## Possible Environments to configure the Agent
 Param|Envoirmentname | value | description | require
 --- |---            | ---   | ---         | ---  
--funkserver| FUNK_SERVER | wss://[url]:[port] | Complete Funk Server URL | true
--connectionkey| CONNECTION_KEY | string | The Key to authenticate against [funk-server](https://github.com/fasibio/funk-server). Is declared at your funk-server | true
--insecureSkipVerify|INSECURE_SKIP_VERIFY | false (default) or true | disable ssl verification for server connection | false
--loglevel|LOG_LEVEL | debug or info (default) or warn or error |Which log-level for the agent own logs | false
--statsintervall|STATSINTERVALL | 15 | If LOG_STATS is not no. than the intervall to collect this information | false
--searchindex|SEARCHINDEX | string(default: "default")| the elasticsearch index to write this information (you name will always append by "_metrics_cumlated") | false
--staticcontent|STATICCONTENT | json | extra information wich will added to each document

At the release Tags you can find the linux binaries. 

You have to run them on your server which metrics will be shipped. 

with --help you can find configuration help
