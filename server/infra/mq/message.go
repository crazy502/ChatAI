package mq

import "errors"

const MessageQueueName = "Message"

var ErrDropMessage = errors.New("drop message")
