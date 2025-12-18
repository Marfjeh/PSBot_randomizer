package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

type Config struct {
	Events   []Event `json:"events"`
	PsbotURL string  `json:"psbot_url"`
	GuildID  string  `json:"guild"`
}

type Event struct {
	Name			string	 `json:"name"`
	RandomMin string   `json:"random_min"`
	RandomMax string   `json:"random_max"`
	Sounds    []string `json:"sounds"`
	UserAgent string   `json:"useragent"`
}

type PsbotBody struct {
	Guild string `json:"guild"`
	Sound string `json:"sound"`
}

//Calcuate a random time between two values and return that.
func randomDuration(min, max time.Duration) time.Duration {
	r := rand.New(rand.NewPCG(1, uint64(time.Now().UnixNano())))
	return min + time.Duration(r.Int64N(int64(max-min)))
}

//Send an POST request to the server to play the sound
func playSound(ctx context.Context, url string, UserAgent string, Psbot PsbotBody) error {
	b := bytes.NewBuffer([]byte{})
	_ = json.NewEncoder(b).Encode(Psbot)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, b)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		fmt.Printf("Played sound %s in guild: %s\n", Psbot.Sound, Psbot.Guild)
	} else {
		return fmt.Errorf("%w: %q", errors.New("Did not got a sensible response from PSBOT :("), resp.Status)
	}

	return nil
}

//Setup thread
func StartPlaying(ctx context.Context, e Event, guildID string, psbotURL string) error {
	//Seed the randomizer
	r := rand.New(rand.NewPCG(1, uint64(time.Now().UnixNano())))
	randomMin, err := time.ParseDuration(e.RandomMin)
	if err != nil {
		return err
	}

	randomMax, err := time.ParseDuration(e.RandomMax)
	if err != nil {
		return err
	}

X:
	for {
		randomTime := randomDuration(randomMin, randomMax)
		log.Printf("[%s] Waiting for: %q\n", e.Name, randomTime.String())
		select {
		case <-ctx.Done():
			break X
		case <-time.Tick(randomTime):
			randomSoundIndex := r.IntN(len(e.Sounds))
			err := playSound(ctx, psbotURL, e.UserAgent, PsbotBody{Guild: guildID, Sound: e.Sounds[randomSoundIndex]})
			if err != nil {
				log.Printf("Something went fucky wucky when playing a sound %v", err)
			}
		}
	}

	return nil
}

func main() {
	ctx := context.Background()

	eg, ctx2 := errgroup.WithContext(ctx)

	cfg, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatalf("Unable to read config file! %v", err)
	}

	c := Config{}
	err = json.NewDecoder(cfg).Decode(&c)

	for _, e := range c.Events {
		eg.Go(func() error {
			return StartPlaying(ctx2, e, c.GuildID, c.PsbotURL)
		})
	}

	err = eg.Wait()
	if err != nil {
		log.Fatalf("One of the threads got killed %v", err)
	}
}
