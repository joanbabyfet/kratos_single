package data

import (
	"context"
	"encoding/json"
	"fmt"
	"kratos_single/internal/conf"
	"net/smtp"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MailRepo struct {
	log *log.Helper
	conf *conf.Data
	conn *amqp.Connection
	ch   *amqp.Channel
}

type MailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewMailRepo(c *conf.Data, logger log.Logger) (*MailRepo, error) {
	// return &MailRepo{
	// 	conf: c,
	// 	log:  log.NewHelper(logger),
	// }

	conn, err := amqp.Dial(c.Rabbitmq.Url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 宣告队列
	_, err = ch.QueueDeclare(
		"test_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	r := &MailRepo{
		conf: c,
		log:  log.NewHelper(logger),
		conn: conn,
		ch:   ch,
	}

	// 启动消费者
	go r.startConsumer()

	return r, nil
}

//原本 Send() 改成：发布到 RabbitMQ
func (r *MailRepo) Send(
	ctx context.Context,
	to string,
	subject string,
	body string,
) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	msg := MailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	//发消息（Producer）
	err = r.ch.PublishWithContext(
		ctx,
		"",
		"test_queue",
		false,
		false,
		amqp.Publishing{ //消息本体 + 属性
			ContentType: "application/json",
			Body:        data,
		},
	)

	if err != nil {
		r.log.Errorf("rabbitmq publish fail: %v", err)
		return err
	}

	r.log.Infof(
		"rabbitmq publish success to=%s subject=%s",
		to,
		subject,
	)

	return nil
}

//发送邮件, 保留原本 Gmail SMTP 写法
func (r *MailRepo) SendSMTP(
	ctx context.Context,
	to string,
	subject string,
	body string,
) error {

	select {
		case <-ctx.Done():
			return ctx.Err()
		default:
	}

	//获取邮件配值内容
	host := r.conf.Mail.Host
	port := r.conf.Mail.Port
	user := r.conf.Mail.User
	pass := r.conf.Mail.Password
	from := r.conf.Mail.From

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
			body,
	)

	auth := smtp.PlainAuth("", user, pass, host)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		from,
		[]string{to},
		msg,
	)

	if err != nil {
		r.log.Errorf(
			"gmail send fail to=%s subject=%s err=%v",
			to,
			subject,
			err,
		)
		return err
	}

	r.log.Infof(
		"gmail send success to=%s subject=%s",
		to,
		subject,
	)

	return nil
}

//////////////////////////////////////////////////////
// 消费者：真正寄 Gmail 邮件
//////////////////////////////////////////////////////

func (r *MailRepo) startConsumer() {
	
	msgs, err := r.ch.Consume(
		"test_queue",
		"",
		true, // auto ack
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		r.log.Errorf("rabbitmq consume fail: %v", err)
		panic(err)
	}
	r.log.Info("mail consumer started...")

	for d := range msgs {
		
		var msg MailMessage
		err := json.Unmarshal(d.Body, &msg)
		if err != nil {
			r.log.Errorf("json decode fail: %v", err)
			continue
		}
		r.log.Infof(
			"consumer receive mail to=%s subject=%s",
			msg.To,
			msg.Subject,
		)
		
		//发送邮件
		_ = r.SendSMTP(
			context.Background(),
			msg.To,
			msg.Subject,
			msg.Body,
		)
	}
}