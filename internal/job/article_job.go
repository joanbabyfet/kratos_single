package job

import (
	"github.com/go-kratos/kratos/v2/log"
)

//共享同一套 biz/data
type ArticleJob struct {
}

func NewArticleJob() *ArticleJob {
	return &ArticleJob{}
}

func (j *ArticleJob) ClearCache() {
	log.Info("清理文章缓存...")
}