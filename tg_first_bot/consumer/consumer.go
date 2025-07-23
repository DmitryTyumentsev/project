package consumer

import (
	"example.com/errors"
	"example.com/events"
	"time"
)

const limit = 100

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.client.Fetch(limit)
		if err != nil {
			return errors.WrapIfErr("failed start Fetch", err)
		}
		if len(gotEvents) == 0 {
			time.Sleep(2 * time.Duration(time.Second))
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			return errors.WrapIfErr("failed handleEvents", err)
		}
	}
}

func (c *Consumer) handleEvents(ev []events.Event) error {
	for _, e := range ev {
		switch e.Type {
		case events.Message:
			//if e.Delete == true {
			//	if err := c.doCMD(e.ChatID, e.TextPage); err != nil {
			//		return errors.WrapIfErr("failed doCMD", err)
			//	}
			//} else {
			if err := c.doCMD(e.ChatID, e.TextPage); err != nil {
				return errors.WrapIfErr("failed doCMD", err)
			}
		default:
			if err := c.client.Client.SendMessage(e.ChatID, msgUnknownType); err != nil {
				return errors.WrapIfErr("failed sendMessage", err)
			}
		}
	}
	return nil
}
