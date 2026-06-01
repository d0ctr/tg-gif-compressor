package tg

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"d0ctr/tg-gif-compressor/tg/handlers"
	"d0ctr/tg-gif-compressor/utils"

	"github.com/AshokShau/gotdbot"
)

var (
	TELEGRAM_API_HASH = utils.Getenv("API_HASH")
	TELEGRAM_API_ID   = utils.GetenvAs("API_ID", strconv.Atoi)
	TELEGRAM_TOKEN    = utils.Getenv("TELEGRAM_TOKEN")
	TDLIB_PATH        = os.Getenv("TDLIB_PATH")
) 

var tg struct {
	client *gotdbot.Client
	dispatcher *gotdbot.Dispatcher
	commands []gotdbot.BotCommand
}

func getTdlib() string {
	const libPrefix = "/usr/local/lib"
	libdir, err := os.ReadDir(libPrefix)
	if err != nil {
		log.Panicf("Failed to read '%s' to find 'libtdjson.so', fix that or specify 'TDLIB_PATH' in the environment", libPrefix)
	}

	for _, entry := range libdir {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "libtdjson.so") {
			return libPrefix + "/" + entry.Name()
		}
	}

	log.Panicf("Failed to find 'libtdjson.so' in '%s', put it there or specify 'TDLIB_PATH' in the environment", libPrefix)
	return ""
}

func init() {
	if TDLIB_PATH == "" {
		TDLIB_PATH = getTdlib()
	}
}

func NewTgBot() {
	client, err := gotdbot.NewClient(int32(TELEGRAM_API_ID), TELEGRAM_API_HASH, TELEGRAM_TOKEN, &gotdbot.ClientOpts{
		LibraryPath: TDLIB_PATH,
		Logger: slog.Default().With("component", "gotdbot"),
	})
        if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

        client.Dispatcher.ErrorHandler = func(c *gotdbot.Client, ctx *gotdbot.Context, err error) error {
		slog.With("component", "dispatcher-error").Error(fmt.Sprintf("got an error: %v", err))
		return nil
	};

	tg.client, tg.dispatcher = client, client.Dispatcher

	handlers.Register(tg.dispatcher)

	if err := client.Start(); err != nil {
		log.Fatalf("failed to start polling %v", err)
	} else {
		slog.Debug("started polling")
	}

	me, err := client.GetMe()
	if err != nil {
		log.Fatalf("failed to get persona: %v", err)
	} else {
		log.Printf("impersonating '%s'", me.Usernames.ActiveUsernames[0])
	}

	// if err := tg.client.SetCommands(tg.commands, "", nil); err != nil {
	// 	log.Fatalf("failed to set commands %v", err)
	// } else {
	// 	log.Printf("registered commands %v", tg.commands)
	// }
}

func Wait() {
	tg.client.Idle()
}

func Stop() {
	// empty? tg.client doesn't have a stop
}

