# Cloudflare Provider

## Requirements
Setup requires an API Token created with Zone:Zone:Read and Zone:DNS:Edit permissions for all zones in your account.

ref : https://developers.cloudflare.com/api/tokens/create

## Setup 
Add this block configuration in the config.yaml file:

```
cloudflare_account:
  enable: true
  secret:
    api_key: "123456789azerty987654321qwerty159753"
    email: "myemail@sibenj@address.com"
  record:
    name: "example.com"
    ttl: 60
    domain: "example.com"
```

## Use

To use start docker use this command:

```
docker run -itd -v "./config.yaml:/config/config.yaml" ghcr.io/azrod/updateip:latest
```

