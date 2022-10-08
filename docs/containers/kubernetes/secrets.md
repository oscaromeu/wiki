---
id: secrets
title: Secrets
sidebar_label: Secrets
sidebar_position: 10
---

## Generate a file with encrypted secrets

1. First we need to generate the base64, independently if it's a file or a string line.

+ For strings: `echo -n mipassword |base64- w0`
+ For files: `cat fichero |base64 -w0`

The `-w0` have to be used always to make sure that all is placed on one line. When string line use `echo -n` to make sure that `\n` is not introduced at the end of the line. 

2. Generate the secret file

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: kafka-credentials
  namespace: monitoring
data:
  file_name_or_password_name: base_64_result
```

3. Encrypt previous kubernetes secret with sops and the age public key

```
sops --age=age1helqcqsh9464r8chnwc2fzj8uv7vr5ntnsft0tn45v2xtz0hpfwq98cmsg \ -e sops-test-secret.yaml | tee sops-test-secret-encrypted.yaml
```

