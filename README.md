### Before running app:

Create config yaml file
```
storage_path: 
secret_key: 
salt: 
expiration_time: 
http_server:
  address: 
  timeout: 
  idle_timeout: 
```

... and export env variable ```CONFIG_PATH``` with path to config like:
```
export CONFIG_PATH=/home/user/config.yaml
```

### Run migrations:
up:
```
cd migrations
goose sqlite3 ../storage.db up
```

down:
```
cd migrations
goose sqlite3 ../storage.db down
```

### Working with makefile:

Seed data:
```
make seed
```

Build app:
```
make build
```

Run app:
```
make run
```