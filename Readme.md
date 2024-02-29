# Widget exporter

## Usage
```shell
Usage of autofaq-exporter:
  -c, --config string   Path to config file
```

## Config
Config format:
```yaml
log-level: info                     # level log of output. One of 'debug', 'info', 'warn' or 'error'
autofaq-url: https://autofaq.local  # URL of AutoFAQ site
listen-address: 0.0.0.0:8080        # address on which to bind and expose metrics [:9901]
services:                           # list of AutoFAQ servicea and widgets of service
  - id: ee3b949d-db82-476c-b5ee-12819e030155
    widgets:
      - 6b40ee48-1ccf-4305-bab1-8f74f1ea4e4d
      - ...
```
### Environment vars
* LOG_LEVEL - level of logs output
* AUTOFAQ_URL - URL of AutoFAQ site
* LISTEN_ADDRESS - address on which to bind and expose metrics

## ToDo
* widget check
* exception with incorrect address of AF