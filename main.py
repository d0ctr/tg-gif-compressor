import os
import asyncio
import subprocess
from telethon import TelegramClient, events
from ffmpy import FFmpeg

# Environment Variables from Dockerfile for security
api_id = os.getenv('API_ID')
api_hash = os.getenv('API_HASH')
bot_token = os.getenv('TELEGRAM_TOKEN')

dev_null = open(os.devnull, 'w')

client = TelegramClient('bot', api_id, api_hash).start(bot_token=bot_token)


@client.on(events.NewMessage(pattern='/start'))
async def start(event):
    await event.reply('Если я вдруг получу ещё одну гифку как документ, клянусь, я потеряю терпение')

@client.on(events.NewMessage)
async def handle_message(event: events.NewMessage.Event):
    if event.message.media:
        if event.message.document.mime_type == 'image/gif':
            print("received gif as document")

            async with client.action(event.chat_id, "video") as action:
                file_path = await client.download_media(event.message.document)
                output_path = f"{file_path}.mp4"

                print("converting")
                convert_gif_to_mp4(file_path, output_path)

                print("sending")
                await client.send_file(event.chat_id, file=output_path, reply_to=event.message.id, progress_callback=action.progress)

                os.remove(file_path)
                os.remove(output_path)

def convert_gif_to_mp4(input: str, output: str):
    ffmpeg = FFmpeg(inputs={input: None}, outputs={output: None})
    ffmpeg.run()

async def main():
    await client.start(bot_token=bot_token)
    print("Bot is running...")
    await client.run_until_disconnected()

loop = asyncio.get_event_loop()
loop.run_until_complete(main())