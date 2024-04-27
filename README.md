<p align="center"><img src="https://raw.githubusercontent.com/dgmann/document-manager/master/apps/frontend/src/assets/icons/icon-512x512.png" alt="DocumentManager" height="265"></p>
<h1 align="center">DocumentManager</h1>
<p align="center">A system for managing medical records.</p>

## Quickstart
```shell
wget https://github.com/dgmann/document-manager/releases/latest/download/docker-compose.yaml
touch .env # Fill it with your config. See sectin below. 
docker compose up -d
# Access at http://localhost/
```
## Configuration

Configuration is done primarily through environment variables. For an overview, check the [./deployment/docker-compose.yml](./deployment/docker-compose.yml).
The mandatory configuration values are provided by a `.env` file with the following values:

| Env            | Key                            |
|----------------|--------------------------------|
| M1_HOST        | IP or hostname to M1 Oracle DB |
| M1_DB_USERNAME | username for the M1 Oracle DB  |
| M1_DB_PASSWORD | username for the M1 Oracle DB  |
