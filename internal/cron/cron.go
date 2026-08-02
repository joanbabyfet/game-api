package cron

import (
	"game-api/internal/config"

	robfigcron "github.com/robfig/cron/v3"
)

func Start(recoveryJob *GameOrderRecoveryJob) *robfigcron.Cron {
	c := robfigcron.New(
		robfigcron.WithChain(
			robfigcron.SkipIfStillRunning(robfigcron.DefaultLogger),
			robfigcron.Recover(robfigcron.DefaultLogger),
		),
	)
	if _, err := c.AddJob(config.Cfg.Cron.OrderRecover, recoveryJob); err != nil {
		panic(err)
	}
	c.Start()
	return c
}
