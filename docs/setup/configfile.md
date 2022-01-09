# Configuration File

See all options available in configuration file. 

By default Update loads the configuration file into the `/config` directory and the configuration file is called `config.yaml`.
it is possible to change the directory and the configuration file which are called when the program starts. 

* `PATH_CONFIG_DIRECTORY` : directory where the configuration file is located.
* `PATH_CONFIG_FILE` : configuration file name.


## All options available

| Options                                 | Default  | Required                 | Actions                                                         |
| --------------------------------------- | -------- | ------------------------ | --------------------------------------------------------------- |
| log.level                               | info     | :heavy_multiplication_x: | Set Log Level *(accept trace, debug, info, warn, error, fatal)* |
| log.humanize                            | false    | :heavy_multiplication_x: | Set human log format                                            |
|                                         |          |                          |                                                                 |
| metrics.enable                          | false    | :heavy_multiplication_x: | Define if start metrics web server                              |
| metrics.host                            | 0.0.0.0  | :heavy_multiplication_x: | Set IP address for metrics web server                           |
| metrics.port                            | 8080     | :heavy_multiplication_x: | Set port for metrics web server                                 |
| metrics.path                            | /metrics | :heavy_multiplication_x: | Path for acceding to metrics web server                         |
| metrics.logging                         | false    | :heavy_multiplication_x: | Logging request http endpoint                                   |
|                                         |          |                          |                                                                 |
| providers.aws.enable                    | false    | :heavy_multiplication_x: | Enable AWS Route 53 Provider                                    |
| providers.aws.secret.access_key_id      | ""       | :heavy_check_mark:       | AccessKey for AWS Account                                       |
| providers.aws.secret.secret_access_key  | ""       | :heavy_check_mark:       | SecretKey for AWS Account                                       |
| providers.aws.secret.region             | ""       | :heavy_check_mark:       | Region for your domain                                          |
| providers.aws.record.name               | ""       | :heavy_check_mark:       | FQDN record *(ex.domain.com)*                                   |
| providers.aws.record.domain             | ""       | :heavy_check_mark:       | Domain Name *(domain.com)*                                      |
| providers.aws.record.hosted_zone_id     | ""       | :heavy_multiplication_x: | HostedZoneID of your domain                                     |
|                                         |          |                          |                                                                 |
| providers.ovh.enable                    | false    | :heavy_multiplication_x: | Enable OVH Provider                                             |
| providers.ovh.secret.application_key    | ""       | :heavy_check_mark:       | Application Key for OVH Account                                 |
| providers.ovh.secret.application_secret | ""       | :heavy_check_mark:       | Application Secret for OVH Account                              |
| providers.ovh.secret.consumer_key       | ""       | :heavy_check_mark:       | Consumer Key for OVH Account                                    |
| providers.ovh.secret.region             | ""       | :heavy_check_mark:       | Region for your domain                                          |
| providers.ovh.record.name               | ""       | :heavy_check_mark:       | FQDN record *(ex.domain.com)*                                   |
| providers.ovh.record.zone               | ""       | :heavy_check_mark:       | DNS Zone *(domain.com)*                                         |
|                                         |          |                          |                                                                 |
| providers.cloudflare.enable             | false    | :heavy_multiplication_x: | Enable Cloudflare Provider                                      |
| providers.cloudflare.secret.api_key     | ""       | :heavy_check_mark:       | API Key for cloudflare Account                                  |
| providers.cloudflare.secret.email       | ""       | :heavy_check_mark:       | Email for Cloudfalre Account                                    |
| providers.cloudflare.record.name        | ""       | :heavy_check_mark:       | FQDN record *(ex.domain.com)*                                   |
| providers.cloudflare.record.domain      | ""       | :heavy_check_mark:       | DNS Zone *(domain.com)*                                         |
| providers.cloudfalre.record.zone_id     | ""       | :heavy_check_mark:       | ID DNS Zone                                                     |

```yaml
log:
  level: debug # Available : trace debug info warn error fatal panic
  humanize: true # Default: false

metrics:
  enable: true # Default: false
  port: 8080 # Default : 8080
  host: 0.0.0.0 # Default: 0.0.0.0
  path: /metrics # Default: /metrics
  logging: true # Default: false

providers:
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
