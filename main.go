package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func randomDuration(min, max time.Duration) time.Duration {
	return min + time.Duration(rand.Int63n(int64(max-min)))
}

func makePostRequest(url string, payload []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Smoke Detector/1.0 battery is low")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Send POST request, the people are safe again OwO:")
	}

	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Printf("Your smoke detector has low battery.")
	rand.New(rand.NewSource(time.Now().UnixNano()))
 
	url := "URL HERE"
	payload := []byte(`{"guild":"GUILD_HERE","sound":"black_noise"}`)

	for {
		duration := randomDuration(1*time.Minute, 2*time.Hour)
		fmt.Printf("Waiting for %v\n", duration)

		time.Sleep(duration)

		fmt.Println("Low battery alert!111")

		if err := makePostRequest(url, payload); err != nil {
			fmt.Println("Error making POST request:", err)
		}

		fmt.Println("Going to sleep... Zzz...")
	}

}
