# Signing Updates

Inspired by [Sparkle](https://sparkle-project.org)
**Warning** When you lose the private key you can do no longer updates.
Don't forget to document the signature format, here it is ecdsa-SHA1

```
# generate an ES256 key pair using the Eliptic Curve algorithm
openssl ecparam -genkey -name prime256v1 -noout -out ec_private.pem
openssl ec -in ec_private.pem -pubout -out ec_public.pem

# The public key is bundled with the app so it can verify the updates
touch test-app
tar -czvf update.tar.gz test-app ec_public.pem

# create base64 signature, which is published with the update
openssl dgst -ecdsa-with-SHA1 -sign ec_private.pem -out signature.bin update.tar.gz
openssl base64 -in signature.bin -out signature.sha1

# Publish the update and the signature (ecdsa-SHA1), make sure to use https with at least TLS 1.2
# https://yourserver.com/update.html

# unpack update
mkdir /tmp/update
tar -xzf update.tar.gz -C /tmp/update

# the client downloads update and verifys it
openssl base64 -d -in signature.sha1 -out /tmp/signature.bin
openssl dgst -ecdsa-with-SHA1 -verify /tmp/update/ec_public.pem -signature /tmp/signature.bin update.tar.gz
rm /tmp/signature.bin
rm -r /tmp/update

```

## Why is it here SHA1 and not SHA256?
ecdsa-with-SHA256 is not supported by Openssl/pam_pkcs11

## Should you encrypt the update?
If the attacker has access to the system, encrypting the update doest not prevent him from extracting the update contents. But it helps in reducing humans errors, for example Users can't accidentally unpack the update.

> Don't go overboard or you may lock yourself out. I did it two times.

## Update tips
* Host the update on a separate subdomain.
* Backup your private keys really well, just having them on a dev computer or server is not enough. (maybe print as QR?)
* The changelog should be a format with textsize and color, so you can use big red text if the user needs to do an extra step before or after the update.
* The changelog maybe needs to be translated.
* Ship delta updates with `bsdiff`






