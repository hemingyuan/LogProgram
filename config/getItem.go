package config

import (
	"fmt"
	"strconv"
)

// GetInt return int config
func (c *Conf) GetInt(key string) (result int, err error) {
	c.lock.RLock()
	data, ok := c.confMap[key]
	c.lock.RUnlock()

	if !ok {
		err = fmt.Errorf("key:[%s] is not exist", key)
		return
	}
	result, err = strconv.Atoi(data)
	return
}

// GetDefaultInt return int config or default config
func (c *Conf) GetDefaultInt(key string, defval int) (result int) {
	c.lock.RLock()
	data, ok := c.confMap[key]
	c.lock.RUnlock()

	if !ok {
		return defval
	}
	result, err := strconv.Atoi(data)
	if err != nil {
		result = defval
	}
	return
}

// GetString return string config
func (c *Conf) GetString(key string) (result string, err error) {
	c.lock.RLock()
	result, ok := c.confMap[key]
	c.lock.RUnlock()

	if !ok {
		err = fmt.Errorf("key:[%s] is not exist", key)
	}
	return
}

// GetDefaultString return string config or defailt config
func (c *Conf) GetDefaultString(key string, defval string) (result string) {
	c.lock.RLock()
	result, ok := c.confMap[key]
	c.lock.RUnlock()

	if !ok {
		result = defval
	}
	return
}
