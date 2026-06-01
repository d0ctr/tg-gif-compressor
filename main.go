package main

import (
	"d0ctr/tg-gif-compressor/tg"
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug.Level(),
	})))

}

func main() {
	tg.NewTgBot()

	tg.Wait()
}
