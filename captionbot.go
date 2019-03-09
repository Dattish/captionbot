package captionbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"
)

const (
	fileRequestUrl = "https://www.captionbot.ai/api/upload"
	urlRequestUrl  = "https://captionbot.azurewebsites.net/api/messages"
	rateRequestUrl = "https://captionbot.azurewebsites.net/api/messages"
)

type request struct {
	Type    string
	Content string
}

func URLCaption(imageUrl string) (string, error) {
	body, err := json.Marshal(request{"CaptionRequest", imageUrl})
	if err != nil {
		return "", err
	}

	response, err := http.Post(urlRequestUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	responseBody := &bytes.Buffer{}
	_, err = responseBody.ReadFrom(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	caption := responseBody.String()
	return strconv.Unquote(caption)
}

func FileCaption(filePath string) (string, error) {
	imageUrl, err := uploadFile(filePath)
	if err != nil {
		return "", err
	}
	return URLCaption(imageUrl)
}

func Rate(rating int) (string, error) {
	body, err := json.Marshal(request{"Feedback", strconv.Itoa(rating)})
	if err != nil {
		return "", err
	}

	response, err := http.Post(rateRequestUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	responseBody := &bytes.Buffer{}
	_, err = responseBody.ReadFrom(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	caption := responseBody.String()
	return strconv.Unquote(caption)
}

func uploadFile(filePath string) (string, error) {
	var buffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&buffer)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="file"; filename="%s"`,
			strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(filePath)))

	contentType, err := getContentType(filePath)
	if err != nil {
		return "", err
	}

	h.Set("Content-Type", contentType)
	fileWriter, err := multipartWriter.CreatePart(h)

	if err != nil {
		return "", err
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return "", err
	}

	multipartWriter.Close()

	request, err := http.NewRequest("POST", fileRequestUrl, &buffer)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(response.Status)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(response.Body)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()

		imageUrl := body.String()
		return strconv.Unquote(imageUrl)
	}
}

func getContentType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)

	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}
