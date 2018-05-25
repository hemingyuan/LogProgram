package config

// Notifyer 用于配置更新时，回调通知应用层更新配置信息
type Notifyer interface {
	Callback(*Conf)
}
