FROM ubuntu:18.04

COPY ./assets/sources.list /etc/apt/sources.list

RUN set -ex \
	&& export DEBIAN_FRONTEND=noninteractive \
	&& chmod -R 777 /tmp \
	&& apt-get clean \
	&& mv /var/lib/apt/lists /tmp \
	&& mkdir -p /var/lib/apt/lists/partial \
	&& apt-get clean \
	&& apt-get update \
	&& apt-get install -y tzdata

RUN set -ex \
	&& apt-get install -y \
	make flex g++ libmysqlclient-dev libmysql++-dev \
	openjdk-8-jdk \
	mysql-client-5.7 \
	nano \
	python3.7 \
	wget \
	busybox

RUN set -ex \
	&& ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
	&& dpkg-reconfigure -f noninteractive tzdata

RUN set -ex \
	&& wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz \
	&&  tar -xvf go1.13.1.linux-amd64.tar.gz \
	&& mv go /usr/lib/go \
	&& rm go1.13.1.linux-amd64.tar.gz

WORKDIR /home/judge/

COPY ./assets/docker-entrypoint.sh /usr/local/bin/
COPY ./core /home/judge/core
COPY ./assets/judge.conf /home/judge/hustoj_core/etc/judge.conf
COPY ./assets/java0.policy /home/judge/hustoj_core/etc/java0.policy
COPY ./assets/install.sh /home/judge/install.sh

RUN set -ex \
	&& chmod a+x install.sh\
	&& sh install.sh

RUN set -ex \
	&& chmod +x /usr/local/bin/docker-entrypoint.sh \
	&& ln -s /usr/local/bin/docker-entrypoint.sh  /docker-entrypoint.sh