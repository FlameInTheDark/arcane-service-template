package discord

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

// NewEmbedStruct generated embed
type NewEmbedStruct struct {
	ms        *discordgo.MessageSend
	embLength int
}

func truncateText(text string, length int) string {
	if length > 3 && len(text) > length && len(text) > 3 {
		text = text[:length-3] + "..."
	}
	return text
}

// NewEmbed creates new embed
func NewEmbed(title string) *NewEmbedStruct {
	title = truncateText(title, 256)
	return &NewEmbedStruct{&discordgo.MessageSend{Embed: &discordgo.MessageEmbed{Title: title}}, len(title)}
}

// CheckLength returns true if length of embed chars less then 6000
func (emb *NewEmbedStruct) checkLength(newLength int) bool {
	if emb.embLength+newLength <= 6000 {
		return true
	}
	return false
}

// Field adds field to embed
func (emb *NewEmbedStruct) Field(name, value string, inline bool) *NewEmbedStruct {
	if len(name) > 0 && len(value) > 0 {
		name = truncateText(name, 256)
		value = truncateText(value, 1024)
		newLength := len(name + value)
		if emb.checkLength(newLength) {
			emb.ms.Embed.Fields = append(emb.ms.Embed.Fields,
				&discordgo.MessageEmbedField{
					Name:   name,
					Value:  value,
					Inline: inline})
			emb.embLength += newLength
		}
	}
	return emb
}

// TimeStamp adds timestamp to footer of embed
func (emb *NewEmbedStruct) TimeStamp(ts string) *NewEmbedStruct {
	emb.ms.Embed.Timestamp = ts
	return emb
}

// Author adds author to embed
func (emb *NewEmbedStruct) Author(name, url, iconURL string) *NewEmbedStruct {
	name = truncateText(name, 256)
	newLength := len(name)
	if emb.checkLength(newLength) {
		emb.ms.Embed.Author = &discordgo.MessageEmbedAuthor{URL: url, Name: name, IconURL: iconURL}
		emb.embLength += newLength
	}
	return emb
}

// Desc adds description to embed
func (emb *NewEmbedStruct) Desc(desc string) *NewEmbedStruct {
	if len(desc) > 0 {
		desc = truncateText(desc, 2048)
		newLength := len(desc)
		if emb.checkLength(newLength) {
			emb.ms.Embed.Description = desc
			emb.embLength += newLength
		}
	}
	return emb
}

// URL adds url to embed description
func (emb *NewEmbedStruct) URL(url string) *NewEmbedStruct {
	emb.ms.Embed.URL = url
	return emb
}

// Footer adds footer text
func (emb *NewEmbedStruct) Footer(text string) *NewEmbedStruct {
	text = truncateText(text, 2048)
	newLength := len(text)
	if emb.checkLength(newLength) {
		emb.ms.Embed.Footer = &discordgo.MessageEmbedFooter{Text: text}
		emb.embLength += newLength
	}
	return emb
}

// Color adds color to embed
func (emb *NewEmbedStruct) Color(color int) *NewEmbedStruct {
	emb.ms.Embed.Color = color
	return emb
}

// AttachImg adds attached image to embed from io.Reader
func (emb *NewEmbedStruct) AttachImg(name string, file io.Reader) *NewEmbedStruct {
	emb.ms.Embed.Image = &discordgo.MessageEmbedImage{URL: "attachment://" + name}
	emb.ms.Files = append(emb.ms.Files, &discordgo.File{Name: name, Reader: file})
	return emb
}

// AttachImgURL adds attached image to embed from url
func (emb *NewEmbedStruct) AttachImgURL(url string) *NewEmbedStruct {
	emb.ms.Embed.Image = &discordgo.MessageEmbedImage{URL: url}
	return emb
}

// AttachThumbURL adds attached thumbnail to embed from url
func (emb *NewEmbedStruct) AttachThumbURL(url string) *NewEmbedStruct {
	emb.ms.Embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: url}
	return emb
}

// Send send embed message to Discord
func (emb *NewEmbedStruct) GetMessageSend() *discordgo.MessageSend {
	return emb.ms
}

// GetEmbed returns discords embed
func (emb *NewEmbedStruct) GetEmbed() *discordgo.MessageEmbed {
	return emb.ms.Embed
}
