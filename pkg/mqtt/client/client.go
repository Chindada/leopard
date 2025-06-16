package client

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	qos = 1
)

type Client interface {
	Connected() bool
	Disconnect()
	Unsubscribe(topic string)
	Subscribe(topic string, callback mqtt.MessageHandler)
	Publish(topic string, retained bool, payload any)
}

type client struct {
	connected     bool
	connectedLock sync.RWMutex

	uid      string
	password string
	host     string
	port     string

	mq     mqtt.Client
	logger Logger

	onConnHandler mqtt.OnConnectHandler
	onLostHandler mqtt.ConnectionLostHandler

	handler     map[string]mqtt.MessageHandler
	handlerLock sync.RWMutex
}

const (
	retryConnectInterval = 15 * time.Second
	connectTimeout       = 10 * time.Second
)

func NewClient(opts ...Option) (Client, error) {
	client := &client{
		handler: make(map[string]mqtt.MessageHandler),
	}
	for _, opt := range opts {
		opt(client)
	}
	if err := client.validate(); err != nil {
		return nil, err
	}
	mqOpts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("mqtts://%s", net.JoinHostPort(client.host, client.port))).
		SetClientID(client.uid).
		SetUsername(client.uid).
		SetPassword(client.password).
		SetKeepAlive(5 * time.Minute).
		SetConnectRetryInterval(retryConnectInterval).
		SetMaxReconnectInterval(retryConnectInterval).
		SetConnectRetry(true).
		SetCleanSession(true).
		SetConnectTimeout(connectTimeout).
		SetOrderMatters(false).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetOnConnectHandler(client.onConn).
		SetConnectionLostHandler(client.onLost)
	client.mq = mqtt.NewClient(mqOpts)
	client.mq.Connect()
	return client, nil
}

func (c *client) validate() error {
	if c.uid == "" {
		return ErrUIDRequired
	}
	if c.password == "" {
		return ErrPasswordRequired
	}
	if c.host == "" {
		return ErrHostRequired
	}
	if c.port == "" {
		return ErrPortRequired
	}
	if c.logger == nil {
		c.logger = NewLogger(nil)
	}
	c.onConnHandler = c.onConn
	c.onLostHandler = c.onLost
	return nil
}

func (c *client) onConn(mqtt.Client) {
	c.logger.Infof("MQTT broker connected to %s:%s", c.host, c.port)
	c.handlerLock.RLock()
	for topic, handler := range c.handler {
		c.mq.Subscribe(topic, qos, handler)
	}
	c.handlerLock.RUnlock()

	c.connectedLock.Lock()
	c.connected = true
	c.connectedLock.Unlock()
}

func (c *client) onLost(_ mqtt.Client, err error) {
	c.connectedLock.Lock()
	c.connected = false
	c.connectedLock.Unlock()
	c.logger.Errorf("MQTT broker connection lost: %v", err)
}

func (c *client) Connected() bool {
	c.connectedLock.RLock()
	defer c.connectedLock.RUnlock()
	return c.connected
}

func (c *client) Disconnect() {
	c.mq.Disconnect(0)
}

func (c *client) Subscribe(topic string, callback mqtt.MessageHandler) {
	c.handlerLock.Lock()
	c.handler[topic] = callback
	c.handlerLock.Unlock()
	if c.Connected() {
		c.mq.Subscribe(topic, qos, callback)
	}
}

func (c *client) Unsubscribe(topic string) {
	c.handlerLock.Lock()
	delete(c.handler, topic)
	c.handlerLock.Unlock()
	if c.Connected() {
		c.mq.Unsubscribe(topic)
	}
}

func (c *client) Publish(topic string, retained bool, payload any) {
	if c.Connected() {
		c.mq.Publish(topic, qos, retained, payload)
	}
}
