package main

import (
	"context"
	"strings"
)

type Command struct {
	Name      string
	Arguments []string
	Context   context.Context
}

func NewCommand(content string) *Command {
	ss := strings.Fields(content)
	return &Command{
		Name:      ss[0],
		Arguments: ss[1:],
		Context:   context.TODO(),
	}
}

func (c *Command) Argument(index int) string {
	if index >= len(c.Arguments) || index < 0 {
		return ""
	}
	return c.Arguments[index]
}
