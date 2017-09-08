#!/bin/sh

useradd -M -r -d /opt/gocmdbd -c 'gocmdbd Service' gocmdbd

mkdir -p /var/log/gocmdbd /opt/gocmdbd/{bin,etc}

cp gocmdbd /opt/gocmdbd/bin
cp config.json mysql.json /opt/gocmdbd/etc
cp gocmdbd.service /etc/systemd/system

chown -R gocmdbd:gocmdbd /var/log/gocmdbd /opt/gocmdbd

systemctl enable gocmdbd
systemctl start gocmdbd
