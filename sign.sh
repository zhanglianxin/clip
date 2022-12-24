#!/bin/bash

# ref: https://blog.csdn.net/yxh265/article/details/50519810
#openssl smime -sign \
#    -in unsigned.mobileconfig \
#    -out signed.mobileconfig \
#    -signer cert.pem \
#    -inkey privkey.pem \
#    -certfile fullchain.pem \
#    -outform der -nodetach

openssl smime -sign \
    -in "$1" \
    -out "$2" \
    -signer certs/cert.pem \
    -inkey certs/privkey.pem \
    -certfile certs/chain.pem \
    -outform der -nodetach

# echo file size
wc -c "$2"
