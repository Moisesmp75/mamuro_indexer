package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mamuro_indexer/models"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	url_api_bulk = "http://localhost:4080/api/_bulkv2"
	username     = "admin"
	password     = "Complexpass#123"
)

var (
	errors = []string{}
)

func ConvertDate(date string) time.Time {
	const format = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
	dt, err := time.Parse(format, strings.TrimSpace(date))
	if err != nil {
		return time.Now()
	}
	return dt
}

func readFile(file *os.File) (*mail.Message, error) {
	msg, err := mail.ReadMessage(file)
	if err != nil {
		return &mail.Message{}, err
	}
	return msg, nil
}

func readBody(msg *mail.Message) ([]byte, error) {
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func processData(header mail.Header, body []byte) models.Email {
	mail := models.Email{
		Message_ID:                header.Get("Message-ID"),
		Date:                      ConvertDate(header.Get("Date")),
		From:                      header.Get("From"),
		To:                        strings.Split(header.Get("To"),","),
		Subject:                   header.Get("Subject"),
		Cc:                        strings.Split(header.Get("Cc"),","),
		Mime_version:              header.Get("Mime-Version"),
		Content_Type:              header.Get("Content-Type"),
		Content_Transfer_Encoding: header.Get("Content-Transfer-Encoding"),
		Bcc:                       strings.Split(header.Get("Bcc"),","),
		X_from:                    header.Get("X-From"),
		X_to:                      strings.Split(header.Get("X-To"),","),
		X_cc:                      strings.Split(header.Get("X-cc"),","),
		X_bcc:                     strings.Split(header.Get("X-bcc"),","),
		X_folder:                  header.Get("X-Folder"),
		X_origin:                  header.Get("X-Origin"),
		X_filename:                header.Get("X-FileName"),
		Content:                   string(body),
	}

	return mail
}

func generateMail(pathForlder string) (models.Email, error) {
	file, err := os.Open(pathForlder)
	if err != nil {
		return models.Email{}, nil
	}
	defer file.Close()

	msg, err := readFile(file)
	if err != nil {
		return models.Email{}, nil
	}

	headers := msg.Header
	body, err := readBody(msg)
	if err != nil {
		return models.Email{}, err
	}

	return processData(headers, body), nil
}

func listFolders(path string) []string {
	listFolder := []string{}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	for _, file := range files {
		listFolder = append(listFolder, file.Name())
	}
	return listFolder
}

func walkDirectory(path string, isDir bool, bulkData *models.DocZinc) {
	if !isDir {
		// fmt.Println(path)
		mail, err := generateMail(path)
		if err != nil {
			errors = append(errors, err.Error())
		}
		bulkData.ListEmail = append(bulkData.ListEmail, mail)
		return
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		walkDirectory(path+"\\"+file.Name(), file.IsDir(), bulkData)
	}
}

func postMailBulkV2(bulkData models.DocZinc) {
	json, err := json.Marshal(bulkData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req, err := http.NewRequest(http.MethodPost, url_api_bulk, bytes.NewBuffer(json))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return
	}
	defer resp.Body.Close()
}

func inspectDirectory(pathDirectory string, wg *sync.WaitGroup) {
	bulkData := models.DocZinc{}
	bulkData.Document = "enron_mail"
	walkDirectory(pathDirectory, true, &bulkData)
	fmt.Println(pathDirectory)
	postMailBulkV2(bulkData)
	defer wg.Done()
}

func Indexer(path string) {
	folders := listFolders(path)
	// workers := 10

	var wg sync.WaitGroup

	for _, dir := range folders {
		wg.Add(1)
		go inspectDirectory(path+"\\"+dir, &wg)
	}

	wg.Wait()
}