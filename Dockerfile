# syntax=docker/dockerfile:1

FROM golang:1.22.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /goddy


# Install yt-dlp
RUN apt-get update -y && apt-get install -y \
    ffmpeg \
    python3-pip

RUN python3 -m pip install -U --pre "yt-dlp[default]" --break-system-packages

# Run
CMD ["/goddy"]