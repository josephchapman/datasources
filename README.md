# datasources

Library of Prometheus Exporters




## Docker compose

Run the example stack:
```bash
cd example-stack/
docker compose up
```

Exporter:
```
http://127.0.0.1:2112/metrics
```
```
# HELP weather_cloud_cover_percent The cloud cover.
# TYPE weather_cloud_cover_percent gauge
weather_cloud_cover_percent{location="London"} 0
weather_cloud_cover_percent{location="New York"} 28
# HELP weather_humidity_relative_percent The relative humidity.
# TYPE weather_humidity_relative_percent gauge
weather_humidity_relative_percent{location="London"} 33
weather_humidity_relative_percent{location="New York"} 37
# HELP weather_precipitation_milimeters The precipitation.
# TYPE weather_precipitation_milimeters gauge
weather_precipitation_milimeters{location="London"} 0
weather_precipitation_milimeters{location="New York"} 0
# HELP weather_rain_milimeters The rain.
# TYPE weather_rain_milimeters gauge
weather_rain_milimeters{location="London"} 0
weather_rain_milimeters{location="New York"} 0
# HELP weather_showers_milimeters The showers.
# TYPE weather_showers_milimeters gauge
weather_showers_milimeters{location="London"} 0
weather_showers_milimeters{location="New York"} 0
# HELP weather_temperature_actual_celcius The actual temperature.
# TYPE weather_temperature_actual_celcius gauge
weather_temperature_actual_celcius{location="London"} 11.4
weather_temperature_actual_celcius{location="New York"} 0.4
# HELP weather_temperature_apparent_celcius The apparent temperature.
# TYPE weather_temperature_apparent_celcius gauge
weather_temperature_apparent_celcius{location="London"} 8.3
weather_temperature_apparent_celcius{location="New York"} -3.4
# HELP weather_wind_direction_degrees The wind direction.
# TYPE weather_wind_direction_degrees gauge
weather_wind_direction_degrees{location="London"} 98
weather_wind_direction_degrees{location="New York"} 292
# HELP weather_wind_gusts_knots The wind gusts.
# TYPE weather_wind_gusts_knots gauge
weather_wind_gusts_knots{location="London"} 5.1
weather_wind_gusts_knots{location="New York"} 8.9
# HELP weather_wind_speed_knots The wind speed.
# TYPE weather_wind_speed_knots gauge
weather_wind_speed_knots{location="London"} 1.4
weather_wind_speed_knots{location="New York"} 4.6
```


Prometheus:
```
http://localhost:9090
```

Sample query:
```
weather_cloud_cover_percent{location="London"}
```
Sample result:
```
100
```


Grafana
```bash
http://localhost:3000
```

Dashboard:

![Screenshot](dashboard_example.png?raw=true "Screenshot of 'Weather' dashboard")