package recast

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	recastAPIURL string = "https://api.recast.ai/v1/request"
)

// Client handles text and voice-file requests to Recast.Ai
type Client struct {
	token string
}

// NewClient returns a new Recast.Ai client initialized with a token
func NewClient(token string) *Client {
	return &Client{token: token}
}

// SetToken sets the token for the recast authentication
func (c *Client) SetToken(token string) {
	c.token = token
}

// TextRequest handles request to Recast.AI and returns a Response struct
func (c *Client) TextRequest(text string) (*Response, error) {
	textParam := fmt.Sprintf("text=%s", text)
	req, err := http.NewRequest("POST", recastAPIURL, bytes.NewBuffer([]byte(textParam)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request failed: %s", resp.Status)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	response, err := newResponse(string(body))
	if err != nil {
		return nil, err
	}
	return response, nil
}

// FileRequest handles voice file request to Recast.Ai and returns a Response struct
func (c *Client) FileRequest(filename string) (*Response, error) {
	var request *http.Request
	var file *os.File
	var fileContent []byte
	var filePart io.Writer
	var resp *http.Response
	var err error

	if file, err = os.Open(filename); err != nil {
		return nil, err
	}

	if fileContent, err = ioutil.ReadAll(file); err != nil {
		return nil, err
	}

	defer file.Close()
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	if filePart, err = writer.CreateFormFile("voice", file.Name()); err != nil {
		return nil, err
	}

	if _, err := filePart.Write(fileContent); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	if request, err = http.NewRequest("POST", recastAPIURL, body); err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Token %s", c.token))
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	if resp, err = client.Do(request); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request failed: %s (%s)", resp.Status, string(responseBody))
	}
	response, err := newResponse(string(responseBody))
	if err != nil {
		return nil, err
	}
	return response, nil
}
