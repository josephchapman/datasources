# README

## Adding packages to the Dockerfile

Base image uses this OS:
```bash
docker run -it --rm python:latest \
    /usr/bin/grep -i 'pretty' /etc/os-release
```
```
PRETTY_NAME="Debian GNU/Linux 13 (trixie)"
```

Packages can be found as follows:
```bash
docker run -it --rm python:latest \
    /bin/bash -c "
            /usr/bin/apt update &> /dev/null \
        &&  /usr/bin/apt policy bsdextrautils 2> /dev/null \
        |   /usr/bin/head -n 3
    "
```
```
bsdextrautils:
  Installed: (none)
  Candidate: 2.41.5-0+deb13u1
```