# Getting started


## Docker

The easiest and fastest way to get started with updateip is to use its docker image.
You can run the image with the following command :

```bash

docker run -itd \
    -v "./config.yaml:/config/config.yaml" \
    ghcr.io/azrod/updateip:latest

```

## Basic configuration

Create the **config.yaml** configuration file

```yaml
log:
  level: debug # Available : trace debug info warn error fatal panic
  humanize: true # Default: false

# >> Here Setup your provider <<
```

## Setup your provider

Go to the [provider documentation](providers.md) to setup your provider.
