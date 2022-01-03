# Metrics

Different metrics are available. Metrics are organized into different categories. 

Categories :

* Global
* Per-Providers

## How to enable and configure metrics server

You have some options to enable metrics server.

### Config File

| Options        | Default  | Required                 | Actions                                 |
| -------------- | -------- | ------------------------ | --------------------------------------- |
| metrics.enable | false    | :heavy_multiplication_x: | Define if start metrics web server      |
| metrics.host   | 0.0.0.0  | :heavy_multiplication_x: | Set IP address for metrics web server   |
| metrics.port   | 8080     | :heavy_multiplication_x: | Set port for metrics web server         |
| metrics.path   | /metrics | :heavy_multiplication_x: | Path for acceding to metrics web server |
|                |          |                          |                                         |

```yaml title=config.yaml"
metrics:
  enable: true # Default: false
  port: 8080 # Default : 8080
  host: 0.0.0.0 # Default: 0.0.0.0
  path: /metrics # Default: /metrics

```

## Env Variables


| Options        | Actions                                 |
| -------------- | --------------------------------------- |
|                |                                         |
| METRICS_ENABLE | Define if start metrics web server      |
| METRICS_HOST   | Set IP address for metrics web server   |
| METRICS_PORT   | Set port for metrics web server         |
| METRICS_PATH   | Path for acceding to metrics web server |

```bash title="exemple"
METRICS_ENABLE=true ./updateip
```

## Metrics details

### Global

| Metrics Name                 | Description          |
| ---------------------------- | -------------------- |
|                              |                      |
| updateip_count_event_receive | Total event received |

### AWS Provider

| Metrics Name           | Description                          |
| ---------------------- | ------------------------------------ |
|                        |                                      |
| updateip_aws_func_time | Execution time of each function      |
| updateip_aws_status    | Return Status of AWS Provider        |
| updateip_aws_update    | Number of DNS record validity checks |
