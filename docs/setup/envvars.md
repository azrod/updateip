# Environnement variables

All parameters can be overloaded by environment variables

| Options                   | Actions                                                         |
| ------------------------- | --------------------------------------------------------------- |
| PATH_CONFIG_DIRECTORY     | Directory where the configuration file is located               |
| PATH_CONFIG_FILE          | Configuration file name                                         |
|                           |                                                                 |
| LOG_LEVEL                 | Set Log Level *(accept trace, debug, info, warn, error, fatal)* |
| LOG_HUMANIZE              | Set human log format                                            |
|                           |                                                                 |
| METRICS_ENABLE            | Define if start metrics web server                              |
| METRICS_HOST              | Set IP address for metrics web server                           |
| METRICS_PORT              | Set port for metrics web server                                 |
| METRICS_PATH              | Path for acceding to metrics web server                         |
| METRICS_LOGGING           | Logging request http on endpoint                                |
|                           |                                                                 |
| AWS_ACCOUNT_ENABLE        | Enable AWS Route 53 Provider                                    |
| AWS_ACCESS_KEY_ID         | AccessKey for AWS Account                                       |
| AWS_SECRET_ACCESS_KEY     | SecretKey for AWS Account                                       |
| AWS_REGION                | Region for your domain                                          |
| AWS_RECORD_NAME           | FQDN record *(ex.domain.com)*                                   |
| AWS_RECORD_DOMAIN         | Domain Name *(domain.com)*                                      |
| AWS_HOSTED_ZONE_ID        | HostedZoneID of your domain                                     |
|                           |                                                                 |
| OVH_ACCOUNT_ENABLE        | Enable OVH Provider                                             |
| OVH_APPLICATION_KEY       | Application Key for OVH Account                                 |
| OVH_APPLICATION_SECRET    | Application Secret for OVH Account                              |
| OVH_CONSUMER_KEY          | Consumer Key for OVH Account                                    |
| OVH_REGION                | Region for your domain                                          |
| OVH_RECORD_NAME           | FQDN record *(ex.domain.com)*                                   |
| OVH_RECORD_ZONE           | DNS Zone *(domain.com)*                                         |
|                           |                                                                 |
| CLOUDFLARE_ACCOUNT_ENABLE | Enable Cloudflare Provider                                      |
| CLOUDFLARE_API_KEY        | API Key for cloudflare Account                                  |
| CLOUDFLARE_EMAIL          | Email for Cloudfalre Account                                    |
| CLOUDFLARE_RECORD_NAME    | FQDN record *(ex.domain.com)*                                   |
| CLOUDFLARE_RECORD_DOMAIN  | DNS Zone *(domain.com)*                                         |
| CLOUDFLARE_RECORD_ZONEID  | ID DNS Zone                                                     |
