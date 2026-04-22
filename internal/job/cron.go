package job

import (
	"log"

	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

var ProviderSet = wire.NewSet(
	NewUserJob,
	NewArticleJob,
	NewCronJob,
)

type CronJob struct {
	c *cron.Cron
}

func NewCronJob(userJob *UserJob, articleJob *ArticleJob) *CronJob {

	c := cron.New(
		cron.WithSeconds(), // 支援秒
	)

	// 每10秒执行
	c.AddFunc("*/10 * * * * *", userJob.SyncUser)

	// 每30秒执行
	c.AddFunc("*/30 * * * * *", articleJob.ClearCache)

	return &CronJob{
		c: c,
	}
}

func (j *CronJob) Start() {
	log.Println("cron start...")
	j.c.Start()
}

func (j *CronJob) Stop() {
	log.Println("cron stop...")
	ctx := j.c.Stop()
	<-ctx.Done()
}