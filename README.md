# Event Dispatcher

Application written in Go that can deploy a tracker to track any number of APIs concurrently.
Allows to fetch the API source over pre-determined intervals and be handled appropriately as desired.

## Scheduler Package

While creating the fully generic event dispatcher, a job scheduler package was written which lets functions run periodically at pre-determined interval using a simple, functional syntax.

Heavily inspired from a more powerful pacakge [goCron](https://github.com/jasonlvhit/gocron) but with own needs and ideas. 