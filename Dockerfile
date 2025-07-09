FROM ubuntu:rolling

WORKDIR /app

COPY . .

COPY data/.local/share/applications /root/.local/share/applications
COPY data/usr/share/applications /usr/share/applications


RUN apt-get update && apt-get install -y golang
RUN go build -o /bin/lbdm
