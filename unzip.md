# Unpacking a protected file without using memory

In this example we copy `raspbian.img` image onto a SD card, without putting it into memory. Examples are for macOS but can easily ported to linux.

## Zip
```bash
# pipe zip file to dd (on mac)
diskutil unmountDisk /dev/disk3
unzip -p raspbian.zip raspbian.img | dd of=/dev/diskX bs=4m

# zip img with password
rm raspbian.zip # remove existing zip first to prevent updating it
zip -er raspbian.zip raspbian.img

# pipe zip with password to dd
diskutil unmountDisk /dev/diskX
unzip -P password -p raspbian.zip raspbian.img | dd of=/dev/diskX bs=4m
```

## Zip Verify (optional)
Include a txt file with all the important files and hashes. After dd mount the volume and check them.
```bash
unzip -P password -p raspbian.zip hashes.txt | ./verify_boot.sh
```

## OpenSSL
OpenSSL gives you stronger encryption than zip, but that is usually not needed for updates, because the attacker has access to the system.
```bash
openssl enc -aes-128-cbc -salt -in raspbian.img -out encrypted.enc -k password
# openssl enc -d -aes-128-cbc -in encrypted.enc -out raspbian.img -k password

diskutil unmountDisk /dev/diskX
openssl enc -d -aes-128-cbc -in encrypted.enc -k password | dd of=/dev/diskX bs=4m
```

