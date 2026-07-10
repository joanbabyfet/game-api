package cron

import "github.com/robfig/cron/v3"

func Start() {

    c := cron.New()

    c.Start()

}