# replicator chess

Replicates ELO into a TSDB.

In a system that doesn't support cron, run the binary with the `-nocron <minutes>` flag to keep the process running in a loop that sleeps for `<minutes>` mintues.

## Docker

```bash
PLAYERS='[
    "annacramling",
    "AlexandraBotez"
]'
```

```bash
docker build -t replicator-chess -f replicator-chess/Dockerfile .
$ docker run --network=host replicator-chess
```

(`--network=host` is required while the TSDB host is defined as `http://localhost:8086`)

Note that the `cmn` directory is required in the build container, so the `docker build` command must be run from the top-level directory.

## Output

```
{"time":"2025-04-11T12:54:38.166721107Z","level":"INFO","msg":"Archive found annacramling-2020-01","application":"replicator-chess"}
{"time":"2025-04-11T12:54:38.16676383Z","level":"INFO","msg":"Skipping annacramling-2020-01, already in database","application":"replicator-chess"}
{"time":"2025-04-11T12:54:38.16676745Z","level":"INFO","msg":"Archive found annacramling-2020-02","application":"replicator-chess"}
{"time":"2025-04-11T12:54:38.166769955Z","level":"INFO","msg":"Skipping annacramling-2020-02, already in database","application":"replicator-chess"}
...
rating,player=alexandrabotez,time_class=bullet elo=2453 1743017783000000000
rating,player=alexandrabotez,time_class=bullet elo=2446 1743017870000000000
rating,player=alexandrabotez,time_class=bullet elo=2441 1743135155000000000

{"time":"2025-04-11T12:54:39.038500581Z","level":"INFO","msg":"Query written: https://api.chess.com/pub/player/alexandrabotez/games/2025/03","application":"replicator-chess"}
{"time":"2025-04-11T12:54:39.038514016Z","level":"INFO","msg":"Archive incomplete. Pull not recorded alexandrabotez-2025-03","application":"replicator-chess"}
```