# First Service Go

## Comandos

### Criar Banco de Dados
```bash
docker exec -it postgres_db psql -U admin -d first_service_go -c "CREATE DATABASE surfbook_dev;"
```

### Verificando se o rodou a criação da DATABASE
```bash
docker exec postgres_db psql -U admin -d first_service_go -tAc "SELECT 1 FROM pg_database WHERE datname='surfbook_dev';" 
```

### Rodar a Migration dentro do container Docker
```bash
docker-compose exec -T postgres psql -U admin -d first_service_go < ./migrations/00001-create-tables.up.sql
```

### Rodar os testes da camada de Service
```bash
go test -v service/*.go 
```

```bash
bash test_services.sh
```

### Snyk job de seguraça
[https://snyk.io/pt-BR/](https://snyk.io/pt-BR/)