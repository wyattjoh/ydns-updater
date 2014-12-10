# yDNS Updater

[![Gobuild Download](http://gobuild.io/badge/github.com/codeskyblue/gobuild/download.png)](http://gobuild.io/github.com/wyattjoh/ydns-updater)

A lightweight appplication which updates a dns entry on [yDNS](https://ydns.eu/) using systemd units provided in `systemd/`.

## Getting started

### From source

```bash
# Get the code
go get github.com/wyattjoh/ydns-updater

# Install the systemd files
cp $GOPATH/src/github.com/wyattjoh/ydns-updater/systemd/ydns-updater.{timer,service} /etc/systemd/system

# Edit the systemd file
#
# Adjust /root/go/bin to where your $GOPATH/bin directory is for your user
#
# --host "<HOST TO UPDATE>"
# --user "<API USERNAME>" Found https://ydns.eu/api/
# --pass "<API PASS>" Found https://ydns.eu/api/
vim /etc/systemd/system/ydns-updater.service

# Start and enable timer
systemctl start ydns-updater.timer
systemctl enable ydns-updater.timer
```

### Precompiled

1. Visit http://gobuild.io/github.com/wyattjoh/ydns-updater and download the binary
2. Download systemd unit files and install into `/etc/systemd/system`
3. Start timer `systemctl start ydns-updater.timer`
4. Enable timer `systemctl enable ydns-updater.timer`
