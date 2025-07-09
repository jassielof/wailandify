FROM golang:latest

WORKDIR /app

COPY . .

COPY data/.local/share/applications /root/.local/share/applications
COPY data/usr/share/applications /usr/share/applications

RUN go build -o /bin/lbdm
