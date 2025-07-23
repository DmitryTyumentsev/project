package consumer

import (
	"example.com/errors"
	"example.com/storage"
	"net/url"
	"strings"
)

const (
	commandHello      = "/hello"
	commandPickRandom = "/random"
	commandDelete     = "/delete"
	msgHello          = "Привет! Сейчас бот умеет сохранять ссылки и возвращать рандомно одну из сохраненных раз в неделю. \n " +
		"Если хочешь получить ссылку не дожидаясь недели, отправь команду /random . \n" +
		"Если хочешь удалить сохраненную ссылку, отправь команду /delete . \n" +
		"Если бот прислал ссылку, она автоматически будет удалена"
	msgOk = "Ссылка успешно сохранена"
	//msgDeleteOk    = "Ссылка успешно удалена"
	msgDeleteInfo  = "Пришлите ссылку которую нужно удалить"
	msgUnknownType = "Не сохранено, это не ссылка. Пришлите ссылку в формате https://www.ww"
	msgIsExist     = "Такая ссылка уже сохранена"
)

func (c *Consumer) doCMD(chatID int, textPage string) error {
	tPage := strings.TrimSpace(textPage)

	switch tPage {
	case commandHello:
		if err := c.client.Client.SendMessage(chatID, msgHello); err != nil {
			return errors.WrapIfErr("failed sendMessage", err)
		}

	case commandPickRandom:
		page, err := c.internalBasePath.PickRandom() // chatID ?
		if err != nil {
			return errors.WrapIfErr("failed PickRandom", err)
		}
		if err := c.client.Client.SendMessage(chatID, page.TextPage); err != nil {
			return errors.WrapIfErr("failed sendMessage", err)
		}
		err = c.internalBasePath.Remove(page)
		if err != nil {
			return errors.WrapIfErr("failed Remove", err)
		}

	case commandDelete:
		if err := c.client.Client.SendMessage(chatID, msgDeleteInfo); err != nil {
			return errors.WrapIfErr("failed sendMessage", err)
		}
		//event := events.Event{
		//	TextPage: textPage,
		//	Type: events.Message,
		//	Meta: events.Meta{
		//		ChatID: chatID,
		//		Delete: true,
		//	},
		//}
		//if err := c.handleEvents(event)
		//err := c.internalBasePath.Remove()
		//if err != nil {
		//	return errors.WrapIfErr("failed Remove", err)
		//}

	default:
		if err := c.savePage(tPage, chatID); err != nil {
			return errors.WrapIfErr("failed savePage", err)
		}
	}
	return nil
}

func (c *Consumer) savePage(textPage string, chatID int) error {
	resIsUrl, err := isUrl(textPage)
	if err != nil {
		return errors.WrapIfErr("failed isUrl", err)
	}
	page := storage.Page{
		ChatID:   chatID,
		TextPage: textPage,
	}
	if resIsUrl == true {
		resIsExist, err := c.internalBasePath.IsExist(page)
		if err != nil {
			return errors.WrapIfErr("failed isExist", err)
		}

		if resIsExist == true {
			if err := c.client.Client.SendMessage(chatID, msgIsExist); err != nil {
				return errors.WrapIfErr("failed sendMessage", err)
			}
		}
		if resIsExist == false {
			if err := c.internalBasePath.Save(page); err != nil {
				return errors.WrapIfErr("failed save doCMD", err)
			}
			if err := c.client.Client.SendMessage(chatID, msgOk); err != nil {
				return errors.WrapIfErr("failed sendMessage", err)
			}
		}
	}
	if resIsUrl == false {
		if err := c.client.Client.SendMessage(chatID, msgUnknownType); err != nil {
			return errors.WrapIfErr("failed sendMessage", err)
		}
	}
	return nil
}

func isUrl(textPage string) (bool, error) {
	u, err := url.Parse(textPage)
	if err != nil {
		return false, errors.WrapIfErr("failed url.Parse", err)
	}
	return (u.Scheme == "http" || u.Scheme == "Http" || u.Scheme == "https" || u.Scheme == "Https") && u.Host != "", nil
}
