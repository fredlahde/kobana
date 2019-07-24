# Kobana

Generates a yaml defined chroot environment in a blaze. Buil opon alpine.

## Note

This is by no means secure. Please read `man 2 chroot`. *Do not run in production*. Use Docker or containerd for this. One can easily escape this chroot, I do not prevent this by design.  
Build only for linux. No plans for macos yet.
