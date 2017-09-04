#!/bin/sh

systemctl stop gohttpd
systemctl disable gohttpd

userdel gohttpd

rm -fr /var/log/gohttpd /opt/gohttpd
rm /etc/systemd/system/gohttpd.service 
