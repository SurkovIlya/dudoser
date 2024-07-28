package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Value string `json:"value"`
}

const host = "%host%:8080/v1/receiving"
const rps = 2

func main() {
	var (
		data Message
		err  error
	)

	for {
		go func() {
			data.Value, err = RandomString(10)
			if err != nil {
				log.Println(err)
			}

			json_data, err := json.Marshal(data)
			if err != nil {
				log.Printf("неудалось конвертировать объект в байты: %s", err)
			}

			buf := bytes.NewBuffer(json_data)

			resp, err := http.Post(host, "application/json", buf)
			if err != nil {
				log.Printf("неудалось отправить POST запрос: %s", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Printf("status code: %d", resp.StatusCode)
			}
		}()

		time.Sleep((1000 / rps) * time.Millisecond)
	}

}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) (string, error) {
	byteArray := make([]byte, length)
	_, err := rand.Read(byteArray)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		byteArray[i] = letters[int(byteArray[i])%len(letters)]
	}

	return string(byteArray), nil
}
