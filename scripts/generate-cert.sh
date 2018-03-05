#!/bin/bash -e

openssl req -new -x509 -days 365 \
            -key private-key.pem \
            -config configuration.cnf \
            -out certificate.pem