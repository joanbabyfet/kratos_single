package job

import (
	"github.com/go-kratos/kratos/v2/log"
)

type ArticleJob struct {
}

func NewArticleJob() *ArticleJob {
	return &ArticleJob{}
}

func (j *ArticleJob) ClearCache() {
	log.Info("清理文章缓存...")
}