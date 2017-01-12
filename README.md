# yDNS Updater

[![Build Status](https://travis-ci.org/wyattjoh/ydns-updater.svg?branch=master)](https://travis-ci.org/wyattjoh/ydns-updater)

A lightweight appplication which updates a dns entry on [https://ydns.eu/](https://ydns.eu/) using a systemd unit provided in `systemd/ydns-updater.service`. This will take the servers current ip address and update the associated DNS entry. Perfect for connecting to your home network if it has a dynamic ip address.

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

1. Download pre-compiled binary on the [Releases Page](https://github.com/wyattjoh/ydns-updater/releases/latest) for your Arch/OS
2. Download systemd unit file and install into `/etc/systemd/system`
3. Start service `systemctl start ydns-updater.service`
4. Enable service `systemctl enable ydns-updater.service`
