### Local db setup (for local development using)

Dir `local-db-setup` must contain files needed to set up the database  
used for local development.  


### Preconditions
- Local db starts in docker container so required precondition: installed Docker  
- `local-db-setup` dir must contain `init.sql` file with database and user credentials creating,  
user privileges initialization.  Example:
```
CREATE DATABASE calendar;
\connect calendar;
CREATE SCHEMA IF NOT EXISTS calendar;
CREATE USER calendar_user WITH ENCRYPTED PASSWORD 'calendar_pswd';
GRANT CONNECT ON DATABASE calendar TO calendar_user;
GRANT USAGE, CREATE ON SCHEMA calendar TO calendar_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA calendar TO calendar_user;
```

### Run setup
To set up local db for calendar application just run command  
from `local-db-setup` dir:

```
docker compose up -d
```