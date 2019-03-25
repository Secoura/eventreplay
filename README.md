# Event Replay

Event Replay is a plugin for Secoura that reads in an input file, performs some replacements and outputs the result.

This can be used to simulate a continuously streaming set of log files.

# Usage

Create a config file to specify the configuration to pass to the plugin:

```yaml
identifier: bro-logs
delimiter: "\n"
earliest_time: -20m
latest_time: now
replacements:
  - token: ^\d{10}
    type: timestamp
    replacement: "UNIX"
  - token: \w+\t(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\t\d+\t\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\t\d+
    type: file
    replacement: /path/to/ip_address.sample
  - token: \w+\t\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\t\d+\t(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\t\d+
    type: file
    replacement: /path/to/internal_ips.sample
```