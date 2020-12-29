# jwt-validation

## Overview

```yaml:oauth.yaml
# example to yaml file
oauth:
  - domain: localhost
#    aud:
#      - "aaaa"
    iss: "http://localhost/"
    jwk-set-uri: http://localhost/.well-known/jwks.json
cache:
  name: inmemory
  time: 30
```

## Environments

| Key                    | Overview             | default |
| :--------------------- | :------------------- | :------ |
| REDIS_CLUSTER_ENDPOINT | redis host address   |         |
| REDIS_CLUSTER_PORT     | redis port number    |         |
| VALIDATION_FILE_NAME   | チェック用ファイル名 | oauth   |
