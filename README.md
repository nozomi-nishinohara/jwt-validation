# jwt-validation

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
