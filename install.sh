#!/bin/sh

useradd -M -r -d /opt/gocmdbd -c 'gocmdbd Service' gocmdbd

mkdir -p /var/log/gocmdbd /opt/gocmdbd/{bin,etc}

go build

cp gocmdbd /opt/gocmdbd/bin
cp config.json /opt/gocmdbd/etc
cp sql/{database,users}.sql /opt/gocmdbd/etc
cp gocmdbd.service /etc/systemd/system

cat /opt/gocmdbd/etc/*.sql | mysql -p

chown -R gocmdbd:gocmdbd /var/log/gocmdbd /opt/gocmdbd

systemctl enable gocmdbd
systemctl start gocmdbd
