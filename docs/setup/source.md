# Running UpdateIP from source

UpdateIP is develloped in golang and can be executed from source. **Golang 1.17 or later is required**.

```bash

cd /tmp/
git clone git@github.com:azrod/updateip.git && cd updateip

go mod download 
go build

PATH_CONFIG_DIRECTORY=$(pwd) ./updateip
```

[See all available environment variables](envvars.md)