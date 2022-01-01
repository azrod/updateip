# Docker

Docker image is stored in [GitHub](https://github.com/users/azrod/packages/container/package/updateip)

Images is available for Arch :

* linux/amd64
* linux/arm
* linux/arm64

## Basic running

```bash

docker run -itd \
    -v "./config.yaml:/config/config.yaml" \
    ghcr.io/azrod/updateip:latest

```

## Advanced running

You can load the configuration via the configuration file and / or by environment variables
[See all available environment variables](envvars.md)

```bash

docker run -itd \
    -v "./config.yaml:/config/config.yaml" \
    -e "LOG_LEVEL=debug" \
    -e "LOG_HUMANIZE=true" \
    -e "METRICS_ENABLE=true" \
    -p "8080:8080" \
    ghcr.io/azrod/updateip:latest

```
