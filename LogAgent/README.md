# LogAgent

使用`tail`包实时跟踪日志，放入channel中。从channel中读取日志数据通过`sarama`包发送至kafka

第一版 配置文件解析使用`config`包 后续会将配置信息存储到etcd中
