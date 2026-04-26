package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"kratos_single/internal/biz"
	"kratos_single/internal/conf"
	"kratos_single/internal/data"
	"kratos_single/internal/job"
	"kratos_single/internal/pkg/i18n"
	"kratos_single/internal/pkg/utils"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

//注入 cronJob
func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, cronJob *job.CronJob, cli *clientv3.Client, mq *data.MQ) *kratos.App {
	
	// 确认 etcd client 已建立
	reg := etcd.New(cli)
	log.NewHelper(logger).Info("连接 etcd 成功")

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Registrar(reg), // 注册到 etcd

		kratos.Server(
			gs,
			hs,
		),

		// 启动前
		kratos.BeforeStart(func(ctx context.Context) error {

			log.NewHelper(logger).Info("连接 etcd 成功")

			// cron
			cronJob.Start()
			log.NewHelper(logger).Info("Cron Job 已启动")

			// 2. 启动 RabbitMQ Consumer
			err := mq.Consume("test_queue")
			if err != nil {
				return err
			}
			log.NewHelper(logger).Info("RabbitMQ Consumer 已启动")

			// 3. 发送测试消息（上线可删）
			usecase := biz.NewMQUsecase(mq)

			err = usecase.SendWelcomeMessage(ctx, 1001)
			if err != nil {
				return err
			}

			log.NewHelper(logger).Info("测试消息已发送")

			return nil
		}),

		// 停止前
		kratos.BeforeStop(func(ctx context.Context) error {
			log.NewHelper(logger).Info("停止 Cron Job")
			cronJob.Stop()

			log.NewHelper(logger).Info("关闭 RabbitMQ")
			mq.Close()

			return nil
		}),
	)
}

func main() {
	flag.Parse()
	// logger := log.With(log.NewStdLogger(os.Stdout),  //输出到文件改为 log.NewStdLogger(file)
	// 	"ts", log.DefaultTimestamp,
	// 	"caller", log.DefaultCaller,
	// 	"service.id", id,
	// 	"service.name", Name,
	// 	"service.version", Version,
	// 	"trace.id", tracing.TraceID(),
	// 	"span.id", tracing.SpanID(),
	// )

	//初始化 i18n（必须最先）
	i18n.InitI18n()

	//获取环境变量
	env := os.Getenv("KRATOS_ENV")
	if env == "" {
		env = "dev"
	}
	root := utils.RootPath()
	fmt.Println("当前环境:", env)
	fmt.Println("当前目录:", root)

	configPath := filepath.Join(
		root,
		"configs",
		fmt.Sprintf("config.%s.yaml", env),
	)
	fmt.Println("配置路径:", configPath)

	c := config.New(
		config.WithSource(
			//file.NewSource(flagconf),
			file.NewSource(configPath),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	
	//wire 注入（已包含 cronJob）
	app, cleanup, err := wireApp(bc.Server, bc.Data)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 启动
	if err := app.Run(); err != nil {
		panic(err)
	}
}
