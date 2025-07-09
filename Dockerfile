FROM golang:latest

WORKDIR /app

COPY . .

COPY data/.local/share/applications /root/.local/share/applications
COPY data/usr/share/applications /usr/share/applications
COPY data/.config/lbdm /root/.config/lbdm

RUN go install ./cmd/lbdm
