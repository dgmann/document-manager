## Validator
docker-compose run --rm migrator /validator -output /data/import.gob

## Importer
docker-compose run --rm migrator /importer -input /data/import.gob