package clients

import (
	"encoding/json"
	"example.com/authorization"
	"example.com/errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	methodGetUpdates  = "getUpdates"
	methodSendMessage = "sendMessage"
)

func basePath() string {
	token, err := authorization.MustToken()
	if err != nil {
		fmt.Printf("unuccess call mustToken: %s", err)
	}

	return "bot" + token
}

func (c *Client) Update(limit, offset int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(q, methodGetUpdates)
	if err != nil {
		return nil, fmt.Errorf("can't got response body method getUpdates %w", err)
	}

	var resp Response

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("can't get parsed resp.Body %w", err)
	}

	return resp.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	log.Println("Начинаем выполнять SendMessage...")
	log.Printf("ChatID: %d", chatID)
	q := url.Values{}
	q.Add("text", text)
	q.Add("chat_id", strconv.Itoa(chatID))

	_, err := c.doRequest(q, methodSendMessage)
	if err != nil {
		return fmt.Errorf("can't got response body method sendMessage %w", err)
	}

	//var resp Response
	//if err := json.Unmarshal(data, &resp); err != nil {
	//	return fmt.Errorf("can't do Unmarshal %x", err)
	//}
	//if !resp.Ok {
	//	return fmt.Errorf(" failed: resp !Ok %x", err)
	//}
	return nil
}

func (c *Client) doRequest(query url.Values, method string) (data []byte, err error) {
	//log.Println("Начинаем выполнять doRequest...")
	defer func() { err = errors.WrapIfErr("failed func doRequest", err) }()

	u := url.URL{
		Scheme:   "https",
		Host:     c.Host,
		Path:     path.Join(c.BasePath, method),
		RawQuery: query.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
