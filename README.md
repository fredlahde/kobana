[![Build Status](https://travis-ci.org/fredlahde/kobana.svg?branch=master)](https://travis-ci.org/fredlahde/kobana)

# Kobana

Generates a yaml defined chroot environment in a blaze. Build opon [Alpine](https://alpinelinux.org/).

## Important

This is by no means secure. Please read `man 2 chroot`. *Do not run in production*. Use Docker or containerd for this. One can easily escape this chroot, I do not prevent this by design.  
Build only for linux. No plans for macos yet.

## How it works

Kobana mounts a ramfs and unpacks a minimal [Alpine](https://alpinelinux.org/) Linux into it. It then chroots into it and ensures the given enrionment

## Requirements

* Linux
* Some free ram
* Support for ramfs: `cat /proc/filesystems | grep ramfs`

## TODO

* Design Doc about the sructure of the config file
