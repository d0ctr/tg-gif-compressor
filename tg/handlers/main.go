package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AshokShau/gotdbot"
	"github.com/AshokShau/gotdbot/handlers"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Register(dispatcher *gotdbot.Dispatcher) {
	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotdbot.Client, ctx *gotdbot.Context) error {
		ctx.EffectiveMessage.ReplyText(b, "Если я вдруг получу ещё одну гифку как документ, клянусь, я потеряю терпение", nil)

		return nil
	}))

	dispatcher.AddHandler(handlers.NewMessage(
		func(msg *gotdbot.Message) bool {
			if msg.Content == nil {
				return false
			}

			document, ok := msg.Content.(*gotdbot.MessageDocument)
			if !ok {
			return false
			}

			if document.Document != nil && document.Document.MimeType == "image/gif" {
				return true
			}

			return false
		},
		func(b *gotdbot.Client, ctx *gotdbot.Context) error {
			log.Print("processing gif document")

			document := ctx.EffectiveMessage.Content.(*gotdbot.MessageDocument)

			file, err := document.Document.Document.Download(b, 0, 0, 1, &gotdbot.DownloadFileOpts{ Synchronous: true })

			if err != nil {
				return fmt.Errorf("failed to download a file: %w", err)
			}

			local := file.Local
			if local == nil {
				return fmt.Errorf("local copy is absent")
			}

			if !local.IsDownloadingCompleted {
				return fmt.Errorf("download has silently failed")
			}

			if local.Path == "" {
				return fmt.Errorf("local file path is empty")
			}

			defer func() {
				file.Delete(b)
			}()

			outputPath := fmt.Sprintf("%s.mp4", local.Path)

			log.Printf("converting '%s' to '%s'", local.Path, outputPath)
			{
				if action, err := ctx.EffectiveMessage.Action(b, gotdbot.ChatActionRecordingVideo{}.GetType(), nil); err == nil {
					action.Start()
					defer action.Stop()
				}

				errorOutput := strings.Builder{}
				err = ffmpeg.
					Input(local.Path).
					Output(outputPath).
					WithErrorOutput(&errorOutput).
					Run()

				if err != nil {
					log.Printf("error stream from ffmpeg: %s", errorOutput.String())
					return fmt.Errorf("ffmpeg errord: %w", err)
				}
			}
			defer func() {
				os.Remove(outputPath)
			}()

			log.Printf("sending '%s'", outputPath)
			{
				if action, err := ctx.EffectiveMessage.Action(b, gotdbot.ChatActionUploadingVideo{}.GetType(), nil); err == nil {
					action.Start()
					defer action.Stop()
				}


				output := gotdbot.GetInputFile(outputPath)
				if _, err = ctx.EffectiveMessage.ReplyAnimation(b, output, nil); err != nil {
					return fmt.Errorf("failed to reply: %w", err)
				}
			}

			return nil
		}, 
	))
}

