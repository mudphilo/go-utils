package Tinify

import (
	"log"
	"os"
)

const VERSION = "1.0"

var (
	key    string
	client *Client
)

func SetKey(set_key string) {
	key = set_key
}

func GetClient() *Client {

	if len(key) == 0 {

		key = os.Getenv("TINIIFY_API_KEY")
	}

	if len(key) == 0 {

		log.Printf("Provide an API key with Tinify.setKey(key string)")
		return nil
	}

	if client == nil {
		c, err := NewClient(key)
		if err != nil {

			log.Printf("Provide an API key with Tinify.setKey(key string) got error %v", err.Error())
			return nil
		}

		client = c
	}
	return client
}
