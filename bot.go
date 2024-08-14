package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/wader/goutubedl"
)

const (
	InstagramMainURL = "https://www.instagram.com/"
	XMainURL         = "https://x.com/"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("No TELEGRAM_TOKEN found in env")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(MainHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	log.Info("Bot is started 🎉")
	b.Start(ctx)
}

func CreateTempDirectoryForChat(chatID int64) string {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("%v/voice/", chatID))
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Info("Failed to create temp directory, maybe it's exist: %v", err)
	}
	log.Info("Directory for chat %v is created: %v", chatID, dir)
	return dir
}

func DownloadFile(url string, filepath string) error {
	log.Info("Downloading audio from: %v \t to: %v", url, filepath)
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Info("File successfully downloaded from: %v \t to: %v", url, filepath)
	return nil
}

func MainHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil && update.Message.Text != "" {
		text := update.Message.Text
		if strings.Contains(text, XMainURL) {
			processMessage(ctx, b, update.Message, text, XMainURL)
		} else if strings.Contains(text, InstagramMainURL) {
			processMessage(ctx, b, update.Message, text, InstagramMainURL)
		} else {
			log.Info("Text does not contain an Instagram or X link.")
		}
	}
}

func processMessage(ctx context.Context, b *bot.Bot, msg *models.Message, text, mainURL string) {
	log.Info("Found %s link in the text: %s", mainURL, text)
	b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
		ChatID:    msg.Chat.ID,
		MessageID: msg.ID,
		Reaction: []models.ReactionType{
			{Type: models.ReactionTypeTypeEmoji,
				ReactionTypeEmoji: &models.ReactionTypeEmoji{
					Emoji: "👀",
				}},
		},
	})
	if url := extractLink(text, mainURL); url != "" {
		downloadedVideo := download(url)
		fileName := downloadedVideo.Name()
		videoData, err := os.ReadFile(fileName)

		if err != nil {
			log.Info("Error reading video file: %v", err)
			return
		}

		_, err = b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:  msg.Chat.ID,
			Caption: fmt.Sprintf("✍️ @%s \n💌 %s", msg.From.Username, msg.Text),
			Video:   &models.InputFileUpload{Filename: "video.mp4", Data: bytes.NewReader(videoData)}, // Wrap []byte in bytes.NewReader
		})

		if err != nil {
			log.Info("Error sending video: %v", err)
		}

		err = os.Remove(fileName)
		if err != nil {
			log.Info("Error removing file: %v", err)
		}

		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    msg.Chat.ID,
			MessageID: msg.ID,
		})
	}
}
func extractLink(message, mainURL string) string {
	start := strings.Index(message, mainURL)
	if start != -1 {
		t := message[start:]
		end := strings.Index(t, " ")
		if end == -1 {
			end = len(t)
		}
		link := t[:end]
		log.Info("Link extracted: %s", link)
		return link
	}
	return ""
}

func download(url string) *os.File {
	result, err := goutubedl.New(context.Background(), url, goutubedl.Options{})
	if err != nil {
		log.Fatal(err)
	}
	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		log.Fatal(err)
	}
	defer downloadResult.Close()
	f, err := os.Create("output.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(f, downloadResult)
	return f
}
