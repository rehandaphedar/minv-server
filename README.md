# Introduction

This is the official server for [minv](https://sr.ht/~rehandaphedar/minv), built with golang, sqlc, sqlite, chi, and ffmpeg.

# Install Instructions

## Docker + Compose

Navigate to an empty directory

``` shell
mkdir minv-server
cd minv-server
mkdir data
```

Download required files and edit `data/config.toml`:

``` shell
wget https://git.sr.ht/~rehandaphedar/minv-server/blob/main/docker-compose.yaml
wget https://git.sr.ht/~rehandaphedar/minv-server/blob/main/data/config.toml -O data/config.toml
```

Note: If you are editing `port` in `data/config.toml`, make sure to change it in `docker-compose.yaml` as well.

Run the container

``` shell
docker compose up
```


## Manual

### Development Dependencies

- go
- [sqlc](https://sqlc.dev)

### Runtime Dependencies

- ffmpeg

### Compiling

Clone the source code

``` shell
git clone https://git.sr.ht/~rehandaphedar/minv-server
cd minv-server
```

Generate SQL and build

``` shell
sqlc generate
go build .
```

### Deploying

Copy `minv-server` and `data/config.toml`. Then, edit `config.toml` and run:

``` shell
./minv-server
```

# Misc

- You can use the [nohup](https://linux.die.net/man/1/nohup) command to easily redirect logs to a file. Also helpful when trying to keep the app alive after exiting an SSH session.
