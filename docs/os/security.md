---
sidebar_position: 99
title: Security
sidebar_label: Security
---


## Creating a Self-Signed SSL Certificate

### What is a Self-Signed SSL Certificate? #

A self-signed SSL certificate is a certificate that is signed by the person who created it rather than a trusted certificate authority. Self-signed certificates can have the same level of encryption as the trusted CA-signed SSL certificate.

Web browsers do not recognize the self-signed certificates as valid. When using a self-signed certificate, the web browser shows a warning to the visitor that the web site certificate cannot be verified.

### Prerequisites

+ openssl

### Creating a Self-Signed SSL Certificate

To create a new Self-Signed SSL Certificate, use the `openssl req` command:

```
openssl req -newkey rsa:4096 \
            -x509 \
            -sha256 \
            -days 3650 \
            -nodes \
            -out example.crt \
            -keyout example.key
```

+ `-newkey rsa:4096` - Creates a new certificate request and 4096 bit RSA key. The default one is 2048 bits.
+ `-x509` - Creates a X.509 Certificate.
+ `-sha256` - Use 265-bit SHA (Secure Hash Algorithm).
+ `-days 3650` - The number of days to certify the certificate for. 3650 is ten years. You can use any positive integer.
+ `-nodes` - Creates a key without a passphrase.
+ `-out example.crt` - Specifies the filename to write the newly created certificate to. You can specify any file name.
+ `-keyout example.key` - Specifies the filename to write the newly created private key to. You can specify any file name.

Press enter and answer the questions 

```bash
Generating a RSA private key
......................................................................++++
........++++
writing new private key to 'example.key'
-----
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
```

```
Country Name (2 letter code) [AU]:US
State or Province Name (full name) [Some-State]:Alabama
Locality Name (eg, city) []:Montgomery
Organization Name (eg, company) [Internet Widgits Pty Ltd]:Linuxize
Organizational Unit Name (eg, section) []:Marketing
Common Name (e.g. server FQDN or YOUR name) []:johndoe.com
Email Address []:hello@johndoe.com
```

### Creating Self-Signed SSL Certificate without Prompt

The self-signed SSL certificate can be generated without being prompted for any question with the option `-subj`.

```
openssl req -newkey rsa:4096 \
            -x509 \
            -sha256 \
            -days 3650 \
            -nodes \
            -out example.crt \
            -keyout example.key \
            -subj "/C=SI/ST=Ljubljana/L=Ljubljana/O=Security/OU=IT Department/CN=www.example.com"
```

The fields, specified in -subj line are listed below:

+ `C=` - Country name. The two-letter ISO abbreviation.
+ `ST=` - State or Province name.
+ `L=` - Locality Name. The name of the city where you are located.
+ `O=` - The full name of your organization.
+ `OU=` - Organizational Unit.
+ `CN=` - The fully qualified domain name.