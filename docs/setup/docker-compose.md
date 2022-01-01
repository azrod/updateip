# Docker Compose

Docker image is stored in [GitHub](https://github.com/users/azrod/packages/container/package/updateip)

Images is available for Arch :

* linux/amd64
* linux/arm
* linux/arm64

## Setup compose file

```bash
version: '3.8'

services:
  updateip:
    container_name: updateaip
    image: ghcr.io/azrod/updateip:latest
    volumes:
      - ./config.yaml:/config/config.yaml:ro
    environnement:
      - LOG_LEVEL=debug
      - LOG_HUMANIZE=true
      - METRICS_ENABLE=true
      - METRICS_PORT=8080
      - METRICS_HOST=0.0.0.0
    port:
      - "8080:8080"
    restart: always
```

[See all available environment variables](envvars.md)