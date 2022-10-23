#!/bin/bash

USERNAME=$1
GROUP=$2
FOLDER="rbac"

openssl genrsa -out ${FOLDER}/${USERNAME}.key 2048

openssl req -new \
-key ${FOLDER}/${USERNAME}.key -out ${FOLDER}/${USERNAME}.csr \
-subj "/CN=${USERNAME}/O=${GROUP}"

openssl x509 -req -days 365 \
  -in ${FOLDER}/${USERNAME}.csr \
  -CA rbac/ca.crt \
  -CAkey rbac/ca.key \
  -CAcreateserial \
  -out ${FOLDER}/${USERNAME}.crt