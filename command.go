package main

// Command todo
type Command struct {
	action func(*Context)
}

// NewCommand todo
func NewCommand(action func(*Context)) *Command {
	return &Command{action: action}
}

// Run todo
func (c *Command) Run(ctx *Context) {
	c.action(ctx)
}
