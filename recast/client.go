package recast

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	recastAPIURL string = "https://api.recast.ai/v1/request"
)

// Client handles text and voice-file requests to Recast.Ai
type Client struct {
	token       string
	language    string
	hasLanguage bool
}

// NewClient returns a new Recast.Ai client
// The token will be used to authenticate to Recast.AI API.
// The language, if provided will define the mlanguage of the inputs sent to Recast.AI, to use the automatic language detection, an empty string must be provided.
func NewClient(token string, language string) *Client {
	return &Client{token: token, language: language}
}

// SetToken sets the token for the Recast.AI API authentication
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetLanguage sets the language used for the requests
func (c *Client) SetLanguage(language string) {
	c.language = language
}

// TextRequest process a text request to Recast.AI API and returns a Response
// opts is a map of parameters used for the request. Two para,eters can be provided: are "token" and "language". They will be used instead of the client token and language(if one is set).
func (c *Client) TextRequest(text string, opts map[string]string) (*Response, error) {
	var token string
	hasLang := false
	lang := ""

	if c.language != "" {
		hasLang = true
		lang = c.language
	}
	if opts != nil {
		if _, ok := opts["language"]; ok {
			hasLang = true
			lang = opts["language"]
		}
		if _, ok := opts["token"]; ok {
			token = opts["token"]
		} else {
			token = c.token
		}
	} else {
		token = c.token
	}

	form := url.Values{}
	if hasLang {
		form.Add("language", lang)
	}
	form.Add("text", text)

	req, err := http.NewRequest("POST", recastAPIURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
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

// FileRequest handles voice file request to Recast.Ai and returns a Response
// TextRequest process a text request to Recast.AI API and returns a Response
// opts is a map of parameters used for the request. Two parameters can be provided: "token" and "language". They will be used instead of the client token and language.
func (c *Client) FileRequest(filename string, opts map[string]string) (*Response, error) {
	var request *http.Request
	var file *os.File
	var fileContent []byte
	var filePart io.Writer
	var langPart io.Writer
	var resp *http.Response
	var err error
	var token string

	hasLang := false
	lang := ""
	if c.language != "" {
		hasLang = true
		lang = c.language
	}
	if opts != nil {
		if _, ok := opts["language"]; ok {
			hasLang = true
			lang = opts["language"]
		}
		if _, ok := opts["token"]; ok == false {
			token = c.token
		} else {
			token = opts["token"]
		}
	} else {
		token = c.token
	}

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

	if hasLang {
		if langPart, err = writer.CreateFormField("language"); err != nil {
			return nil, err
		}

		if _, err := langPart.Write([]byte(lang)); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	if request, err = http.NewRequest("POST", recastAPIURL, body); err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
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
