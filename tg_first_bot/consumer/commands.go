package consumer

import (
	"example.com/errors"
	"example.com/storage"
	"log"
	"net/url"
	"strings"
)

const (
	commandStart      = "/start"
	commandPickRandom = "/random"
	commandDelete     = "/delete"
	msgStart          = "Привет! Сейчас бот умеет сохранять ссылки и возвращать рандомно одну из сохраненных с помощью команды /random \n\n" +
		"Если бот прислал ссылку, она автоматически будет удалена. \n\n" +
		"Если хочешь удалить сохраненную ссылку, отправь команду /delete (в разработке)"
	msgOk = "Ссылка успешно сохранена"
	//msgDeleteOk    = "Ссылка успешно удалена"
	msgDeleteInfo  = "Пришлите ссылку которую нужно удалить"
	msgUnknownType = "Не сохранено, это не ссылка. Пришлите ссылку в формате https://www.ww"
	msgIsExist     = "Такая ссылка уже сохранена"
	msgNotExist    = "Нет сохраненных ссылок. Сначала пришлите ссылку для сохранения"
)

func (c *Consumer) doCMD(chatID int, textPage string) error {
	tPage := strings.TrimSpace(textPage)
	log.Printf("Получено новое сообщение: %s", tPage)

	switch tPage {
	case commandStart:
		_ = c.sendMessage(chatID, msgStart)

	case commandPickRandom:
		page, err := c.internalBasePath.PickRandom(chatID)
		if err != nil {
			return errors.WrapIfErr("failed PickRandom", err)
		}
		if page == nil {
			return c.sendMessage(chatID, msgNotExist)
		}
		_ = c.sendMessage(chatID, page.TextPage)
		err = c.internalBasePath.Remove(page)
		if err != nil {
			return errors.WrapIfErr("failed Remove", err)
		}

	case commandDelete:
		_ = c.sendMessage(chatID, msgDeleteInfo)
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
		log.Printf("Получена ссылка: %s", tPage)
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
		log.Println("Это url")
		resIsExist, err := c.internalBasePath.IsExist(&page)
		if err != nil {
			return errors.WrapIfErr("failed isExist", err)
		}

		if resIsExist == true {
			log.Println("Такой url уже есть")
			_ = c.sendMessage(chatID, msgIsExist)
		}
		if resIsExist == false {
			log.Println("Такого url еще нет, пробуем сохранить")
			if err := c.internalBasePath.Save(page); err != nil {
				return errors.WrapIfErr("failed save doCMD", err)
			}
			log.Println("Сохранили url, пробуем отправить сообщение что успешно сохранено")
			_ = c.sendMessage(chatID, msgOk)
		}
	}
	if resIsUrl == false {
		log.Println("Это не url")
		_ = c.sendMessage(chatID, msgUnknownType)
	}
	return nil
}

func isUrl(textPage string) (bool, error) {
	u, err := url.Parse(textPage)
	if err != nil {
		return false, nil
		//return false, errors.WrapIfErr("failed url.Parse", err)
	}
	return (u.Scheme == "http" || u.Scheme == "Http" || u.Scheme == "https" || u.Scheme == "Https") && u.Host != "", nil
}

func (c *Consumer) sendMessage(chatID int, msg string) error {
	if err := c.client.Client.SendMessage(chatID, msg); err != nil {
		return errors.WrapIfErr("failed sendMessage", err)
	}
	return nil
}
