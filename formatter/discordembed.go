package formatter

import (
	"fmt"
	"github.com/SilphServices/SilphSpout/discord"
	"github.com/SilphServices/SilphSpout/model"
	"strconv"
	"time"
)

type DiscordEmbed struct {
	normalThumbnailURLTemplate string
	shinyThumbnailURLTemplate  string
	nameProvider               model.NameProvider
}

func NewDiscordEmbedFormatter(normalThumbnailURLTemplate, shinyThumbnailURLTemplate string, nameProvider model.NameProvider) DiscordEmbed {
	return DiscordEmbed{
		normalThumbnailURLTemplate: normalThumbnailURLTemplate,
		shinyThumbnailURLTemplate:  shinyThumbnailURLTemplate,
		nameProvider:               nameProvider,
	}
}

func (formatter DiscordEmbed) Format(spawn model.Spawn) (message discord.Message) {
	embed := discord.Embed{}

	nameString := formatter.nameProvider.GetName(spawn.NameID)
	latString := strconv.FormatFloat(spawn.Latitude, 'f', -1, 64)
	lngString := strconv.FormatFloat(spawn.Longitude, 'f', -1, 64)

	embedDescription := ""
	if spawn.IsShiny {
		embedDescription = "**Shiny**\n"
	}

	if spawn.Move1ID >= 0 && spawn.Move2ID >= 0 {
		move1String := formatter.nameProvider.GetMove(spawn.Move1ID)
		move2String := formatter.nameProvider.GetMove(spawn.Move2ID)
		embedDescription = embedDescription + move1String + ", " + move2String
	}

	embedTitle := nameString
	if spawn.HasIV() {
		percent := spawn.IVPercent()
		embed.Color = ivPercentToColor(percent)

		ivString := fmt.Sprintf(`%d%% (%d/%d/%d)`, percent, spawn.IVAttack, spawn.IVDefense, spawn.IVStamina)
		embedTitle = embedTitle + " " + ivString
	}
	embed.Title = embedTitle

	mapTitle := nameString
	if spawn.DespawnUnixSeconds > 0 {
		now := time.Now()
		despawn := time.Unix(spawn.DespawnUnixSeconds, 0)

		despawnString := despawn.Format("15:04:05")
		duration := despawn.Sub(now)
		durationString := formatDuration(duration)
		embedDescription = embedDescription + "\nUntil " + despawnString + " (" + durationString + ")"

		mapTitle = mapTitle + "%20until%20" + despawnString
	}
	embed.URL = "https://maps.google.com/maps?q=" + mapTitle + "+%40" + latString + "," + lngString
	embed.Description = embedDescription

	thumbURLTemplate := formatter.normalThumbnailURLTemplate
	if spawn.IsShiny {
		thumbURLTemplate = formatter.shinyThumbnailURLTemplate
	}
	embed.Thumbnail = discord.Thumbnail{
		URL: fmt.Sprintf(thumbURLTemplate, spawn.NameID),
	}

	message = discord.Message{
		Embeds: []discord.Embed{
			embed,
		},
	}

	return
}

func formatDuration(duration time.Duration) string {
	seconds := int(duration.Seconds()) % 60
	str := strconv.Itoa(seconds) + "s"

	minutes := int(duration.Minutes()) % 60
	if minutes > 0 {
		str = strconv.Itoa(minutes) + "m " + str
	}

	hours := int(duration.Hours())
	if hours > 0 {
		str = strconv.Itoa(hours) + "h " + str
	}
	return str
}

func ivPercentToColor(percent int) int {
	if percent == 100 {
		return 0xff8000 // Legendary Orange
	} else if percent > 90 {
		return 0xa335ee // Epic Purple
	} else if percent > 80 {
		return 0x0070dd // Rare Blue
	} else if percent > 50 {
		return 0x1eff00 // Uncommon Green
	} else if percent > 25 {
		return 0xffffff // Common White
	} else {
		return 0x9d9d9d // Poor Gray
	}
}
