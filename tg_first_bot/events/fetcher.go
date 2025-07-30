package events

import (
	"example.com/clients"
	"example.com/errors"
)

func (p *EventProcessor) Fetch(limit int) (res []Event, err error) {
	defer func() { err = errors.WrapIfErr("Can't convert updateToEvent", err) }()
	upd, err := p.Client.Update(limit, p.offset)
	if err != nil {
		return nil, err
	}

	if len(upd) == 0 {
		return nil, nil
	}
	p.offset = upd[len(upd)-1].UpdateID + 1

	res = make([]Event, 0, len(upd))
	for _, u := range upd {
		res = append(res, fetchUpdate(u))
	}

	return res, nil
}

func fetchUpdate(u clients.Update) Event {
	event := Event{
		TextPage: u.Message.Text,
		Type:     typeEvent(u),
		Meta:     meta(u),
	}
	return event
}

func typeEvent(u clients.Update) EventType {
	switch u.Message.Text {
	case "":
		return UnknownType
	}
	return Message
}

func meta(u clients.Update) Meta { //добавить логику что сообщение с флагом delete true или false
	m := Meta{
		ChatID:   u.ChatID,
		UpdateID: u.UpdateID,
		//UserID:    u.UserID,
		//Username:  u.Username,
		//MessageID: u.MessageID,
	}
	return m
}
