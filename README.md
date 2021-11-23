<center><h1>UpdateIP</h1></center>

UpdateIP is a automatic update of a DNS record based on the External IP address

## How to setup

> Shortly

# How to config

Create the **config.yaml** configuration file

```yaml
log:
  level: debug // Available : trace debug info warn error fatal panic
  humanize: true

aws_account:
  enable: true
  secret:
    access_key_id: "xxx"
    secret_access_key: "xxx"
    region: "eu-west-1"
  record:
    name: "subdomain.domain.com"
    ttl: 60
    domain: "domain.com"

ovh_account:
  enable: true
  secret:
    application_key: "xxx"
    application_secret: "xxx"
    consumer_key: "xxx"
    region: "eu-west-1"
  record:
    name: "subdomain.domain.com"
    ttl: 60
    zone: "domain.com"

cloudflare_account:
  enable: true
  secret:
    api_key: "xxx"
    email: "xxx"
  record:
    name: "subdomain.domain.com"
    ttl: 60
    domain: "domain.com"
```
