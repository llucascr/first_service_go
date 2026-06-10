# Modelagem do Banco de Dados

```sql
CREATE TABLE users (
    user_id UUID PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notebooks (
    notebook_id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(user_id),

    icon VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    image VARCHAR(900),
    description VARCHAR(900),

    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE meta_contents (
    content_id UUID PRIMARY KEY
    notebook_id UUID REFERENCES notebooks(notebook_id), 

    icon VARCHAR(50),
    name VARCHAR(255) NOT NULL,

    user_id UUID NOT NULL,

    deleted_at TIMESTAMP
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE nodes_contents (
    node_id UUID
    content_id UUID REFERENCES meta_contents(content_id), 

    user_id UUID NOT NULL,
    notebook_id UUID NOT NULL,

    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE meta_tag_contents (
    content_id UUID REFERENCES meta_content(content_id),

    notebook_id UUID,
    user_id UUID NOT NULL,

    tag_id UUID REFERENCES tags(tag_id) NOT NULL,

    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tags (
    tag_id UUID PRIMARY KEY,
    
    name VARCHAR(255)
    color VARCHAR(6) NOT NULL,

    user_id UUID REFERENCES users(user_id) NOT NULL,

    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

# Criar Banco de Dados
```bash
docker exec -it postgres_db psql -U admin -d first_service_go -c "CREATE DATABASE surfbook_dev;"
```

# Verificando se o rodou a criação da DATABASE
```bash
docker exec postgres_db psql -U admin -d first_service_go -tAc "SELECT 1 FROM pg_database WHERE datname='surfbook_dev';" 
```

# Rodar a Migration dentro do container Docker
```bash
docker-compose exec -T postgres psql -U admin -d first_service_go < ./migrations/00001-create-tables.up.sql
```