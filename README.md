# yDNS Updater

A lightweight appplication which updates a dns entry on [yDNS](https://ydns.eu/) using systemd units provided in `systemd/`.

## Getting started

```bash
# Get the code
go get github.com/wyattjoh/ydns-updater

# Install the systemd files
cp $GOPATH/src/github.com/wyattjoh/ydns-updater/systemd/ydns-updater.{timer,service} /etc/systemd/system

# Edit the systemd file to add in params
# --host "<HOST TO UPDATE>"
# --user "<API USERNAME>"
# --pass "<API PASS>"
vim /etc/systemd/system/ydns-updater.service

# Start and enable timer
systemctl start ydns-updater.timer
systemctl enable ydns-updater.timer
```
