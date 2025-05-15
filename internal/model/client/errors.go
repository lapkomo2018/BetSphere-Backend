package client

import "errors"

var (
	ErrClientClosed        = errors.New("client is closed")
	ErrClientReadChanFull  = errors.New("read chan is full")
	ErrClientWriteChanFull = errors.New("write chan is full")
)
