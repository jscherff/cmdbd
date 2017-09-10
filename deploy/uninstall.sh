#!/bin/sh

systemctl stop gocmdbd
systemctl disable gocmdbd

userdel gocmdbd

rm -fr /var/log/gocmdbd /opt/gocmdbd
rm /etc/systemd/system/gocmdbd.service 
