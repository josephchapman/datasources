# prometheus exporter weather

Provides an OpenMetrics (`/metrics`) page containing weather data for multiple locations


# Docker

```bash
export WEATHER_LOCATIONS='[
    {
        "name": "New York",
        "latitude": 40.79,
        "longitude": -73.96,
        "tzdata": "America/New_York"
    },
    {
        "name": "London",
        "latitude": 51.50,
        "longitude": -0.12,
        "tzdata": "Europe/London"
    }
]'
```

```bash
docker build -t prometheus-exporter-weather -f prometheus-exporter-weather/Dockerfile .
docker run -it --rm -p 2112:2112 prometheus-exporter-weather
```

Note that the `cmn` directory is required in the build container, so the `docker build` command must be run from the top-level directory.
