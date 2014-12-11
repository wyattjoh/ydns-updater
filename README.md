# yDNS Updater

[![Gobuild Download](http://gobuild.io/badge/github.com/codeskyblue/gobuild/download.png)](http://gobuild.io/github.com/wyattjoh/ydns-updater)

A lightweight appplication which updates a dns entry on [https://ydns.eu/](https://ydns.eu/) using systemd unit provided in `systemd/ydns-updater.service`.

## Getting started

### From source

```bash
# Get the code
go get github.com/wyattjoh/ydns-updater

# Install the systemd files
cp $GOPATH/src/github.com/wyattjoh/ydns-updater/systemd/ydns-updater.service /etc/systemd/system

# Edit the systemd file
#
# Adjust /root/go/bin to where your $GOPATH/bin directory is for your user
#
# --host "<HOST TO UPDATE>"
# --user "<API USERNAME>" Found https://ydns.eu/api/
# --pass "<API PASS>" Found https://ydns.eu/api/
vim /etc/systemd/system/ydns-updater.service

# Start and enable service
systemctl start ydns-updater.service
systemctl enable ydns-updater.service
```

### Precompiled

1. Visit http://gobuild.io/github.com/wyattjoh/ydns-updater and download the binary
2. Download systemd unit file and install into `/etc/systemd/system`
3. Start service `systemctl start ydns-updater.service`
4. Enable service `systemctl enable ydns-updater.service`