FROM ubuntu:jammy

RUN apt update
RUN apt-get -y install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

RUN mkdir -p /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
RUN echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update
RUN apt-get -y install docker-ce docker-ce-cli containerd.io docker-compose-plugin

RUN curl -LO https://go.dev/dl/go1.18.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.18.4.linux-amd64.tar.gz

ENV PATH "$PATH:/usr/local/go/bin"

RUN service docker start

WORKDIR /mnt/data
WORKDIR /calzone

COPY . .

RUN go build .

RUN useradd -rm -d /home/cali -s /bin/bash -g docker -G sudo -u 1001 cali
RUN usermod -aG docker cali

EXPOSE 61895

CMD service docker start && su -c "/calzone/service" -s /bin/sh cali