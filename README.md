# raudit - simple audit event logger for Redis Enterprise


```
Usage:
  -A, --append                        append to log file instead of overwriting
  -c, --console                       use console output as well as defined logger
      --exclude-auth-connection       exclude authentication connections
      --exclude-auth-status uints     exclude authentication status by code (default [])
      --exclude-close-connection      exclude closing connections
      --exclude-internal-connection   exclude internal connections
      --exclude-new-connection        exclude new connections
  -F, --facility string               syslog facility to use (default "LOCAL0")
  -L, --listen strings                address:port (s) to listen on (default [127.0.0.1:29001])
  -f, --logfile string                log file to write to (default to stderr)
  -l, --same-log                      use the same log for both internal and external messages
  -s, --syslog                        send logs to syslog
  ```