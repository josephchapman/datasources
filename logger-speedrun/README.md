# logger speedrun

Produces log entries (to stdout) describing the current records in specified categories.
Logs can be queried via logql and displayed in a dashboard.

In a system that doesn't support cron, run the binary with the `-nocron <minutes>` flag to keep the process running in a loop that sleeps for `<minutes>` mintues.

# Docker

```bash
export SPEEDRUN_LEADERBOARDS='[
  {
    "game": "j1l9qz1g",
    "category": "z275w5k0",
    "values": {
      "p854r2vl": "5q85yy6q"
    }
  },
  {
    "game": "76rljl68",
    "category": "jdzmjxkv",
    "values": {
      "6nj463vn": "01388y3q"
    }
  },
  {
    "game": "o1y9wo6q",
    "category": "wkpoo02r",
    "values": {
      "e8m7em86": "9qj7z0oq"
    }
  },
  {
    "game": "o1y9wo6q",
    "category": "7dgrrxk4",
    "values": {
      "e8m7em86": "9qj7z0oq"
    }
  }
]'
```

```bash
cd logger-speedrun/
docker build -t logger-speedrun .
docker run logger-speedrun
```


