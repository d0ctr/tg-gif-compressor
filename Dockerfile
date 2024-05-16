FROM python:3.9-slim

# Install ffmpeg
RUN apt update && apt install -y ffmpeg && apt clean

# Set the working directory
WORKDIR /app

# Copy the requirements file and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the bot script
COPY main.py .

# Expose the port the app runs on
EXPOSE 80

ENV API_ID=
ENV API_HASH=
ENV TELEGRAM_TOKEN=

CMD ["python", "./main.py"]