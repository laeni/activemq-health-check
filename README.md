用于简单 ActiveMQ 的健康检查。

由于 ActiveMQ 经常出问题，出问题后很南快速感知到，所以需要一个简单的健康检查。此命令一般配合 Zabbix 使用。

## 用法

```shell
./activemq-health-check --host <host> --port <port> --queue <queueName>
```

| 参数 | 说明 |
| host | ActiveMQ 的主机名或 IP 地址 |
| port | ActiveMQ 的端口 |
| queueName | 用于健康检查的测试队列名，检查时会发送一条随机字符串到到此队列，然后检查队列中是否存在消息，所以需要使用一个临时队列 |

示例：

```shell
./activemq-health-check --host 127.0.0.1 --port 61613 --queue zabbix.health
```