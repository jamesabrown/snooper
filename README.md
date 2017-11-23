# snooper

Snooper is an app that allows for quick searching across a large number of servers for data returned via HTTP(S). 

# Usage

```
NAME:
   snooper - Snooper is ready to run

USAGE:
   snooper [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHORS:
   James Brown <jbrown@invoca.com>
   Christian Parkinson <cparkinson@invoca.com>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --filename value     File containing IP's. 1 per line (default: "iplist")
   --port value         Port to send requests to (default: 80)
   --urlPath value      Path for the URL (default: "/is_alive")
   --concurrency value  Number of concurrent requests (default: 10)
   --help, -h           show help
   --version, -v        print the version
```
