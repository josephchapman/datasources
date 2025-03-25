# logger speedrun

Produces log entries (to stdout) describing the current records in specified categories.
Logs can be queried via logql and displayed in a dashboard.

In a system that doesn't support cron, run the binary with the `-nocron <minutes>` flag to keep the process running in a loop that sleeps for `<minutes>` mintues.

## Docker

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
docker build -t logger-speedrun -f logger-speedrun/Dockerfile .
docker run logger-speedrun
```

Note that the `cmn` directory is required in the build container, so the `docker build` command must be run from the top-level directory.

## Output

```
{"time":"2025-03-25T17:17:21.15492674Z","level":"INFO","msg":"Current record retrieved","current_record":{"game":"Super Mario 64","category":"120 Star","player":"Suigi","time":"01h 35m 28s","performed":{"date":"2024-11-03","days_ago":142}},"application":"logger-speedrun"}
{"time":"2025-03-25T17:17:24.244658268Z","level":"INFO","msg":"Current record retrieved","current_record":{"game":"Super Mario 64","category":"70 Star","player":"Suigi","time":"00h 46m 26s","performed":{"date":"2024-11-17","days_ago":128}},"application":"logger-speedrun"}
```