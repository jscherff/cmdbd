#!/bin/sh

useradd -M -r -d /opt/gohttpd -c 'gohttpd Service' gohttpd

mkdir -p /var/log/gohttpd /opt/gohttpd/{bin,etc}

cp gohttpd /opt/gohttpd/bin
cp config.json mysql.json /opt/gohttpd/etc
cp gohttpd.service /etc/systemd/system

chown -R gohttpd:gohttpd /var/log/gohttpd /opt/gohttpd

systemctl enable gohttpd
systemctl start gohttpd
