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
	"io"
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

type Output struct {
	filter model.IVFilter
	poster webhook.Poster
}

func main() {
	logFile, err := os.OpenFile("SilphSpout.log", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.Print("Starting Silph Spout!")

	configPath := flag.String("config", "config/config.json", "Path to JSON file containing configuration")
	ivFilterPath := flag.String("ivfilter", "config/ivFilter.json", "Path to JSON file containing iv filter")

	flag.Parse()

	log.Print("Loading config " + *configPath)
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Print("Couldn't load config. It's formatted incorrectly or in the wrong spot.")
		log.Fatal(err)
	}

	nameProvider, err := model.NewRemoteNameProvider(cfg.NamesJSONURL, cfg.MovesJSONURL)
	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.Outputs) == 0 {
		log.Print("Using legacy output configuration.")
		cfg.Outputs = []config.OutputConfig{
			{
				Service:    "discord",
				WebhookURL: cfg.OutputWebhookURL,
				FilterPath: *ivFilterPath,
			},
		}
	}

	outputs := make([]Output, len(cfg.Outputs))
	for i, outconfig := range cfg.Outputs {
		num := strconv.Itoa(i)

		if outconfig.Service != "discord" {
			log.Fatal("Output " + num + ": Unsupported output service: " + outconfig.Service)
		}

		log.Print("Output " + num + ": Loading ivFilter " + outconfig.FilterPath)
		ivFilter, err := model.LoadFilter(outconfig.FilterPath, nameProvider)
		if err != nil {
			log.Print("Output " + num + ": Couldn't load IV Filter. It's probably in the wrong spot or formatted incorrectly.")
			log.Fatal(err)
		}

		output := Output{
			filter: ivFilter,
			poster: webhook.NewPoster(outconfig.WebhookURL),
		}
		outputs[i] = output
	}

	dedupe := model.NewDedupeFilter()

	discordFormatter := formatter.NewDiscordEmbedFormatter(cfg.NormalThumbnailURLTemplate, cfg.ShinyThumbnailURLTemplate, nameProvider)

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

			if dedupe.Filter(spawn) {
				log.Print("Filtered as duplicate.")
				return
			}

			message := discordFormatter.Format(spawn)

			for i, output := range outputs {
				log.Print(fmt.Sprintf("Posting to webhook: " + strconv.Itoa(i)))

				if output.filter.Filter(spawn) {
					log.Print("Filtered for low IV.")
				} else {
					output.poster.Post(message)
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	})

	port := ":" + strconv.Itoa(cfg.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
