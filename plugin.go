package main

import (
	"encoding/json"
	"os"
	//"time"
	"bytes"
	"flag"
	"fmt"
	"github.com/SilphServices/SilphSpout/config"
	"github.com/SilphServices/SilphSpout/formatter"
	"github.com/SilphServices/SilphSpout/model"
	"github.com/SilphServices/SilphSpout/webhook"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func loadConfig(pathToConfigJSON string) (config config.Config, err error) {
	fd, err := os.Open(pathToConfigJSON)
	decoder := json.NewDecoder(fd)
	err = decoder.Decode(&config)
	return
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	configPath := flag.String("config", "config/config.json", "Path to JSON file containing configuration")
	ivFilterPath := flag.String("ivfilter", "config/ivFilter.json", "Path to JSON file containing iv filter")

	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	nameProvider, err := model.NewRemoteNameProvider(config.NamesJSONURL, config.MovesJSONURL)
	if err != nil {
		log.Fatal(err)
	}

	ivFilter, err := model.LoadFilter(*ivFilterPath, nameProvider)
	if err != nil {
		log.Fatal(err)
	}

	poster := webhook.NewPoster(config.OutputWebhookURL)
	/*
		spawn := model.Spawn {
			NameID: 493,
			Move1ID: 1,
			Move2ID: 2,
			IVAttack: 15,
			IVDefense: 14,
			IVStamina: 13,
			DespawnUnixSeconds: int64(time.Now().Unix() + 10),
			Latitude: 46.851648,
			Longitude: -121.761186,
			IsShiny: false,
		}
	*/
	dedupe := model.NewDedupeFilter()

	discordFormatter := formatter.NewDiscordEmbedFormatter(config.NormalThumbnailURLTemplate, config.ShinyThumbnailURLTemplate, nameProvider)

	http.HandleFunc("/spawn", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Print("Failed to dump " + err.Error())
		}
		log.Print(string(dump))

		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		log.Print("Got new POST body: " + buf.String())

		decoder := json.NewDecoder(buf)
		spawns := make([]model.Spawn, 0)
		err = decoder.Decode(&spawns)
		if err != nil {
			log.Print("Failed to decode " + err.Error())
			return
		}

		for _, spawn := range spawns {
			log.Print(fmt.Sprintf("New Spawn: %+v", spawn))

			if ivFilter.Filter(spawn) {
				log.Print("Filtered for low IV")
				return
			}

			if dedupe.Filter(spawn) {
				log.Print("Filtered as duplicate.")
				return
			}

			message := discordFormatter.Format(spawn)
			log.Print(fmt.Sprintf("Posting to webhook: %+v", message))

			poster.Post(message)
		}
		w.WriteHeader(http.StatusOK)
	})

	port := ":" + strconv.Itoa(config.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
