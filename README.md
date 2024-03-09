用于简单 ActiveMQ 的健康检查。

由于 ActiveMQ 经常出问题，出问题后很南快速感知到，所以需要一个简单的健康检查。此命令一般配合 Zabbix 使用。

## 用法

```shell
./activemq-health-check --host <host> --port <port> --queue <queueName>
```

| 参数 | 必须 | 说明 |
| --- | --- | --- |
| host | 是 | ActiveMQ 的主机名或 IP 地址 |
| port | 是 | ActiveMQ 的端口 |
| username | 否 | 连接 MQ 服务器的用户名 |
| password | 否 | 连接 MQ 服务器的密码 |
| queueName | 是 | 用于健康检查的测试队列名，检查时会发送一条随机字符串到到此队列，然后检查队列中是否存在消息，所以需要使用一个临时队列 |
| persistent | 否 | 发送的消息是否需要持久化 |

示例：

```shell
$ ./activemq-health-check --host 10.10.1.1 --port 61613 --queue zabbix.health
2024/03/05 23:59:09 /home/xxx/activemq-health-check/cmd/root.go:73: [I] - 连接到 ActiveMQ 服务器成功: 10.10.1.1:61613
2024/03/05 23:59:09 /home/xxx/activemq-health-check/cmd/root.go:86: [I] - 发送消息到 ActiveMQ 队列成功！
2024/03/05 23:59:09 /home/xxx/activemq-health-check/cmd/root.go:104: [I] - ActiveMQ 健康检查通过！
```
