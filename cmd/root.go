/*
Copyright © 2024 NAME HERE m@laeni.cn
*/

package cmd

import (
	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// rootCmd 表示在没有任何子命令的情况下调用时的基本命令
var rootCmd = &cobra.Command{
	Use:   "activemq-health-check",
	Short: "Zabbix 插件，用于检查 ActiveMQ 的健康状态",
	Long: `此插件主要有两个作用：
  1. 检查 MQ 本身的健康状态：为了确保 ActiveMQ 正常运行，此插件会像 MQ 发送消息并消费消息，以检查 ActiveMQ 的运行状态.
  2. 检查是否有大量消息对接。

示例:
  activemq-health-check --host 127.0.0.1 --port 61613 --queue zabbix.health`,
	Run: run,
}

// Execute 将所有子命令添加到根命令中，并相应地设置标志。
// 这是由 main.main() 调用的。它只需要对 rootCmd 调用一次。
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	host       string // Active Mq 服务器主机名
	port       string // Active Mq 服务器端口
	username   string // Active Mq 用户名（可选）
	password   string // Active Mq 密码（可选）
	queue      string // 用于测试的队列名
	persistent bool   // 发送的消息是否需要持久化
)

var (
	logInfo    = log.New(os.Stdout, "[I] - ", log.LstdFlags|log.Lmsgprefix|log.Llongfile)
	logWarning = log.New(os.Stderr, "[W] - ", log.LstdFlags|log.Lmsgprefix|log.Llongfile)
	logError   = log.New(os.Stderr, "[E] - ", log.LstdFlags|log.Lmsgprefix|log.Llongfile)
)

func init() {
	// 在这里，您将定义您的标志和配置设置。
	// Cobra 支持持久标志，如果在此处定义，则对于您的应用程序来说，这些标志将是全局的。

	rootCmd.PersistentFlags().StringVarP(&host, "host", "", "", "主机名")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "", "", "端口")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "", "", "连接 MQ 服务器的用户名（可选）")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "", "", "连接 MQ 服务器的密码（可选）")
	rootCmd.PersistentFlags().StringVarP(&queue, "queue", "", "", "队列名")
	rootCmd.PersistentFlags().BoolVarP(&persistent, "persistent", "", false, "发送的消息是否需要持久化。由于作为健康检查，默认情况下消息不会持久化")

	// Cobra 还支持本地标志，这些标志仅在直接调用此操作时运行。
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, _ []string) {
	if host == "" || port == "" {
		cmd.Help()
		os.Exit(1)
	}

	// 连接到 ActiveMQ 服务器
	stompConn, err := stomp.Dial("tcp", net.JoinHostPort(host, port), stomp.ConnOpt.Login(username, password))
	if err != nil {
		logError.Println("连接到 ActiveMQ 服务器失败: " + err.Error())
		os.Exit(1)
	}
	defer stompConn.Disconnect()
	logInfo.Println("连接到 ActiveMQ 服务器成功: " + host + ":" + port)

	var wg sync.WaitGroup
	wg.Add(1)

	// 生成一个随机数
	random := strconv.Itoa(rand.Int())

	// 将消息发送到 ActiveMQ 中
	opts := make([]func(*frame.Frame) error, 0)
	if persistent {
		opts = append(opts, func(f *frame.Frame) error {
			f.Header.Set("persistent", "true")
			return nil
		})
	}
	err = stompConn.Send(queue, "text/plain", []byte(random), opts...)
	if err != nil {
		logError.Println("发送消息到 ActiveMQ 失败: " + err.Error())
	}
	logInfo.Println("发送消息到 ActiveMQ 队列成功！")

	// 消费消息
	sub, err2 := stompConn.Subscribe(queue, stomp.AckAuto)
	if err2 != nil {
		logError.Println("订阅队列失败: " + err2.Error())
	}
	defer sub.Unsubscribe()
	go func() {
		for {
			msg := <-sub.C
			if msg == nil {
				return
			}
			if msg.Err != nil {
				logError.Println("订阅队列失败: " + msg.Err.Error())
			}
			if msg.ContentType == "text/plain" && string(msg.Body) == random {
				logInfo.Println("ActiveMQ 健康检查通过！")
				wg.Done()
			} else {
				logWarning.Println("没有消费到预期消息！预期: " + random + " 得到: " + string(msg.Body))
			}
		}
	}()

	// 等待 10s，如果超时则退出
	go func() {
		time.Sleep(time.Second * 10)
		logError.Println("10 秒无法消费到期望消息，请检查 ActiveMQ 的运行状态")
		wg.Done()
	}()

	wg.Wait()
}
