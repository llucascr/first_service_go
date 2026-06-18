FORMAT: 1A
HOST: http://localhost:8080

# First Service Go

API REST em Go (gorilla/mux + PostgreSQL) para gerenciamento de **notebooks**,
**tags** e seus **conteúdos**. Esta documentação segue o formato
[API Blueprint](https://apiblueprint.org/).

## Convenções

Estas convenções valem para **toda** a API.

### Base URL

Todos os endpoints são relativos a `http://localhost:8080` (porta `8080`,
definida em `cmd/api/main.go`).

### Identificadores via HTTP Header

> ⚠️ **Atenção — convenção incomum.** Diferente da maioria das APIs REST, este
> serviço **não recebe identificadores via path ou query string**. Todos os IDs
> (`user_id`, `notebook_id`, `content_id`, `tag_id`, `node_id`) são enviados
> como **HTTP Headers**. As URIs são fixas (ex.: `/notebook`), e o registro
> afetado é determinado pelos headers da requisição.

Todos os IDs são UUIDs (formato `8-4-4-4-12`). Se um header de ID obrigatório
estiver ausente ou for um UUID inválido, o `uuid.Parse` falha e a resposta é
`500 Internal Server Error` (ver "Modelo de erros").

### Autenticação

Atualmente **não há camada de autenticação**. A identidade do usuário é
assumida a partir do header `user_id`, que é confiado sem validação. Não há
verificação de propriedade entre o `user_id` e os recursos acessados.

### Corpo e tipo de conteúdo

Requisições com corpo usam `application/json`. Respostas de sucesso são
`application/json`. Alguns endpoints de escrita (`POST`/`PUT` de
`nodescontent` e `metatagcontent`) **não leem corpo** — montam a entidade
inteiramente a partir dos headers.

### Soft delete

`DELETE` é lógico: preenche a coluna `deleted_at` em vez de remover a linha.
Retorna `204 No Content` sem corpo.

### Modelo de erros

Os erros **não** são JSON — usam `http.Error`, retornando `text/plain`:

+ `400 Bad Request` — corpo JSON malformado (`JSON inválido: ...`).
+ `500 Internal Server Error` — UUID ausente/inválido em header, ou erro de
  service/banco (`Erro ao salvar notebook: ...`, `Erro ao fazer o parse do
  uuid: ...`, etc.).

Não existem hoje respostas `401`, `404` ou `422`: um registro inexistente ou um
header faltando resultam em `500`, não em `404`/`400`.

# Group Health

Endpoint de verificação de saúde do serviço.

## Health [/health]

### Verificar saúde [GET]

Não exige headers nem autenticação.

+ Response 200 (application/json)

    + Attributes (HealthResponse)

    + Body

            {
                "status": "ok",
                "message": "Service is healthy"
            }

# Group Notebooks

Cadernos pertencentes a um usuário.

## Notebook [/notebook]

Recurso de notebook individual. A ação depende do método HTTP e dos headers.

### Criar notebook [POST]

Cria um notebook para o usuário do header `user_id`.

+ Request (application/json)

    + Headers

            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6

    + Attributes
        + icon: `📓` (string, optional)
        + name: `Estudos de Go` (string, required)
        + image: `https://cdn.exemplo.com/capa.png` (string, optional)
        + description: `Tudo sobre a linguagem Go` (string, optional)

    + Body

            {
                "icon": "📓",
                "name": "Estudos de Go",
                "image": "https://cdn.exemplo.com/capa.png",
                "description": "Tudo sobre a linguagem Go"
            }

+ Response 200 (application/json)

    + Attributes (Notebook)

    + Body

            {
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "icon": "📓",
                "name": "Estudos de Go",
                "image": "https://cdn.exemplo.com/capa.png",
                "description": "Tudo sobre a linguagem Go",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 400 (text/plain)

    + Body

            JSON inválido: unexpected end of JSON input

+ Response 500 (text/plain)

    + Body

            Erro ao salvar notebook: ...

### Buscar notebook por ID [GET]

Retorna o notebook identificado pelo header `notebook_id`.

+ Request

    + Headers

            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90

+ Response 200 (application/json)

    + Attributes (Notebook)

    + Body

            {
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "icon": "📓",
                "name": "Estudos de Go",
                "image": "https://cdn.exemplo.com/capa.png",
                "description": "Tudo sobre a linguagem Go",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao procurar notebook: ...

### Atualizar notebook [PUT]

Atualiza o notebook do header `notebook_id`. O header `user_id` também é lido.

+ Request (application/json)

    + Headers

            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90
            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6

    + Attributes
        + icon: `📗` (string, optional)
        + name: `Go avançado` (string, required)
        + image: `https://cdn.exemplo.com/capa2.png` (string, optional)
        + description: `Concorrência e generics` (string, optional)

    + Body

            {
                "icon": "📗",
                "name": "Go avançado",
                "image": "https://cdn.exemplo.com/capa2.png",
                "description": "Concorrência e generics"
            }

+ Response 200 (application/json)

    + Attributes (Notebook)

    + Body

            {
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "icon": "📗",
                "name": "Go avançado",
                "image": "https://cdn.exemplo.com/capa2.png",
                "description": "Concorrência e generics",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T12:10:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao atualizar notebook: ...

### Deletar notebook [DELETE]

Soft delete do notebook do header `notebook_id`.

+ Request

    + Headers

            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90

+ Response 204

+ Response 500 (text/plain)

    + Body

            Erro ao deletar notebook: ...

## Notebooks do usuário [/notebook/list]

### Listar notebooks do usuário [GET]

Retorna todos os notebooks do usuário do header `user_id`.

+ Request

    + Headers

            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6

+ Response 200 (application/json)

    + Attributes (array[Notebook])

    + Body

            [
                {
                    "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "icon": "📓",
                    "name": "Estudos de Go",
                    "image": "https://cdn.exemplo.com/capa.png",
                    "description": "Tudo sobre a linguagem Go",
                    "deleted_at": null,
                    "created_at": "2026-06-18T11:45:00Z",
                    "updated_at": "2026-06-18T11:45:00Z"
                }
            ]

+ Response 500 (text/plain)

    + Body

            Erro ao procuar lista de Notebooks do user: ...

# Group Tags

Tags coloridas pertencentes a um usuário.

## Tag [/tag]

### Criar tag [POST]

+ Request (application/json)

    + Headers

            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6

    + Attributes
        + name: `Importante` (string, required)
        + color: `FF5733` (string, required) - Cor em HEX sem `#` (até 6 caracteres).

    + Body

            {
                "name": "Importante",
                "color": "FF5733"
            }

+ Response 200 (application/json)

    + Attributes (Tag)

    + Body

            {
                "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                "name": "Importante",
                "color": "FF5733",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 400 (text/plain)

    + Body

            JSON inválido: unexpected end of JSON input

+ Response 500 (text/plain)

    + Body

            Erro ao salvar tag: ...

### Buscar tag por ID [GET]

+ Request

    + Headers

            tag_id: 1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d

+ Response 200 (application/json)

    + Attributes (Tag)

    + Body

            {
                "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                "name": "Importante",
                "color": "FF5733",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao procurar tag: ...

### Atualizar tag [PUT]

+ Request (application/json)

    + Headers

            tag_id: 1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d
            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6

    + Attributes
        + name: `Urgente` (string, required)
        + color: `C0392B` (string, required)

    + Body

            {
                "name": "Urgente",
                "color": "C0392B"
            }

+ Response 200 (application/json)

    + Attributes (Tag)

    + Body

            {
                "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                "name": "Urgente",
                "color": "C0392B",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T12:10:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao atualizar tag: ...

### Deletar tag [DELETE]

+ Request

    + Headers

            tag_id: 1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d

+ Response 204

+ Response 500 (text/plain)

    + Body

            Erro ao deletar tag: ...

## Tags do usuário [/tag/list]

### Listar tags do usuário [GET]

+ Request

    + Headers

            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6

+ Response 200 (application/json)

    + Attributes (array[Tag])

    + Body

            [
                {
                    "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                    "name": "Importante",
                    "color": "FF5733",
                    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "deleted_at": null,
                    "created_at": "2026-06-18T11:45:00Z",
                    "updated_at": "2026-06-18T11:45:00Z"
                }
            ]

+ Response 500 (text/plain)

    + Body

            Erro ao procuar lista de Tags do user: ...

# Group Meta Contents

Conteúdos (seções) que pertencem a um notebook.

## Meta Content [/metacontent]

### Criar meta content [POST]

Cria um meta content no notebook do header `notebook_id`, para o usuário do
header `user_id`.

+ Request (application/json)

    + Headers

            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6
            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90

    + Attributes
        + icon: `📄` (string, optional)
        + name: `Introdução` (string, required)

    + Body

            {
                "icon": "📄",
                "name": "Introdução"
            }

+ Response 200 (application/json)

    + Attributes (MetaContent)

    + Body

            {
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "icon": "📄",
                "name": "Introdução",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 400 (text/plain)

    + Body

            JSON inválido: unexpected end of JSON input

+ Response 500 (text/plain)

    + Body

            Erro ao salvar meta content: ...

### Buscar meta content por ID [GET]

+ Request

    + Headers

            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c

+ Response 200 (application/json)

    + Attributes (MetaContent)

    + Body

            {
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "icon": "📄",
                "name": "Introdução",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao procurar meta content: ...

### Atualizar meta content [PUT]

Atualiza o meta content do header `content_id`. O corpo aceita apenas `icon` e
`name`.

+ Request (application/json)

    + Headers

            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c

    + Attributes
        + icon: `📝` (string, optional)
        + name: `Visão geral` (string, optional)

    + Body

            {
                "icon": "📝",
                "name": "Visão geral"
            }

+ Response 200 (application/json)

    + Attributes (MetaContent)

    + Body

            {
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "icon": "📝",
                "name": "Visão geral",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T12:10:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao atualizar meta content: ...

### Deletar meta content [DELETE]

+ Request

    + Headers

            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c

+ Response 204

+ Response 500 (text/plain)

    + Body

            Erro ao deletar meta content: ...

## Meta Contents do notebook [/metacontent/list]

### Listar meta contents do notebook [GET]

+ Request

    + Headers

            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90

+ Response 200 (application/json)

    + Attributes (array[MetaContent])

    + Body

            [
                {
                    "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                    "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "icon": "📄",
                    "name": "Introdução",
                    "deleted_at": null,
                    "created_at": "2026-06-18T11:45:00Z",
                    "updated_at": "2026-06-18T11:45:00Z"
                }
            ]

+ Response 500 (text/plain)

    + Body

            Erro ao procuar lista de Meta Contents do notebook: ...

# Group Nodes Contents

Nós de conteúdo associados a um meta content e a um notebook.

> Observação: os endpoints de escrita deste grupo **não leem corpo JSON** — a
> entidade é montada inteiramente a partir dos headers.

## Nodes Content [/nodescontent]

### Criar nodes content [POST]

Monta o nó a partir dos headers `user_id`, `notebook_id` e `content_id`. Não há
corpo de requisição.

+ Request

    + Headers

            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6
            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90
            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c

+ Response 200 (application/json)

    + Attributes (NodesContent)

    + Body

            {
                "node_id": "2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f",
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao salvar nodes content: ...

### Buscar nodes content por ID [GET]

+ Request

    + Headers

            node_id: 2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f

+ Response 200 (application/json)

    + Attributes (NodesContent)

    + Body

            {
                "node_id": "2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f",
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao procurar nodes content: ...

### Atualizar nodes content [PUT]

Reassocia o nó do header `node_id` ao meta content do header `content_id`. Não
há corpo de requisição.

+ Request

    + Headers

            node_id: 2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f
            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c

+ Response 200 (application/json)

    + Attributes (NodesContent)

    + Body

            {
                "node_id": "2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f",
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T12:10:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao atualizar nodes content: ...

### Deletar nodes content [DELETE]

+ Request

    + Headers

            node_id: 2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f

+ Response 204

+ Response 500 (text/plain)

    + Body

            Erro ao deletar nodes content: ...

## Nodes Contents do notebook [/nodescontent/list/notebook]

### Listar nodes contents do notebook [GET]

+ Request

    + Headers

            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90

+ Response 200 (application/json)

    + Attributes (array[NodesContent])

    + Body

            [
                {
                    "node_id": "2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f",
                    "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                    "deleted_at": null,
                    "created_at": "2026-06-18T11:45:00Z",
                    "updated_at": "2026-06-18T11:45:00Z"
                }
            ]

+ Response 500 (text/plain)

    + Body

            Erro ao procurar lista de Nodes Contents do notebook: ...

## Nodes Contents do meta content [/nodescontent/list/content]

### Listar nodes contents do meta content [GET]

+ Request

    + Headers

            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c

+ Response 200 (application/json)

    + Attributes (array[NodesContent])

    + Body

            [
                {
                    "node_id": "2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f",
                    "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                    "deleted_at": null,
                    "created_at": "2026-06-18T11:45:00Z",
                    "updated_at": "2026-06-18T11:45:00Z"
                }
            ]

+ Response 500 (text/plain)

    + Body

            Erro ao procurar lista de Nodes Contents do content: ...

# Group Meta Tag Contents

Associação (N:N) entre um meta content e uma tag. A chave primária é o par
`(content_id, tag_id)`.

> Observação: os endpoints de escrita deste grupo **não leem corpo JSON** — a
> entidade é montada inteiramente a partir dos headers.

## Meta Tag Content [/metatagcontent]

### Criar meta tag content [POST]

Vincula a tag (`tag_id`) ao meta content (`content_id`). Monta tudo a partir dos
headers; sem corpo de requisição.

+ Request

    + Headers

            user_id: 3fa85f64-5717-4562-b3fc-2c963f66afa6
            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90
            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c
            tag_id: 1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d

+ Response 200 (application/json)

    + Attributes (MetaTagContent)

    + Body

            {
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao salvar meta tag content: ...

### Buscar meta tag content [GET]

Busca pelo par `content_id` + `tag_id`.

+ Request

    + Headers

            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c
            tag_id: 1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d

+ Response 200 (application/json)

    + Attributes (MetaTagContent)

    + Body

            {
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T11:45:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao procurar meta tag content: ...

### Atualizar meta tag content [PUT]

Atualiza o vínculo do par `content_id` + `tag_id`, reassociando ao notebook do
header `notebook_id`. Sem corpo de requisição.

+ Request

    + Headers

            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c
            tag_id: 1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d
            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90

+ Response 200 (application/json)

    + Attributes (MetaTagContent)

    + Body

            {
                "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "deleted_at": null,
                "created_at": "2026-06-18T11:45:00Z",
                "updated_at": "2026-06-18T12:10:00Z"
            }

+ Response 500 (text/plain)

    + Body

            Erro ao atualizar meta tag content: ...

### Deletar meta tag content [DELETE]

Remove (soft delete) o vínculo do par `content_id` + `tag_id`.

+ Request

    + Headers

            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c
            tag_id: 1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d

+ Response 204

+ Response 500 (text/plain)

    + Body

            Erro ao deletar meta tag content: ...

## Meta Tag Contents do meta content [/metatagcontent/list/content]

### Listar por meta content [GET]

+ Request

    + Headers

            content_id: 5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c

+ Response 200 (application/json)

    + Attributes (array[MetaTagContent])

    + Body

            [
                {
                    "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                    "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                    "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "deleted_at": null,
                    "created_at": "2026-06-18T11:45:00Z",
                    "updated_at": "2026-06-18T11:45:00Z"
                }
            ]

+ Response 500 (text/plain)

    + Body

            Erro ao procurar lista de Meta Tag Contents do content: ...

## Meta Tag Contents da tag [/metatagcontent/list/tag]

### Listar por tag [GET]

+ Request

    + Headers

            tag_id: 1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d

+ Response 200 (application/json)

    + Attributes (array[MetaTagContent])

    + Body

            [
                {
                    "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                    "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                    "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "deleted_at": null,
                    "created_at": "2026-06-18T11:45:00Z",
                    "updated_at": "2026-06-18T11:45:00Z"
                }
            ]

+ Response 500 (text/plain)

    + Body

            Erro ao procurar lista de Meta Tag Contents da tag: ...

## Meta Tag Contents do notebook [/metatagcontent/list/notebook]

### Listar por notebook [GET]

+ Request

    + Headers

            notebook_id: 9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90

+ Response 200 (application/json)

    + Attributes (array[MetaTagContent])

    + Body

            [
                {
                    "content_id": "5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c",
                    "tag_id": "1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d",
                    "notebook_id": "9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90",
                    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "deleted_at": null,
                    "created_at": "2026-06-18T11:45:00Z",
                    "updated_at": "2026-06-18T11:45:00Z"
                }
            ]

+ Response 500 (text/plain)

    + Body

            Erro ao procurar lista de Meta Tag Contents do notebook: ...

# Group Users [ToDo]

> 🚧 **Não implementado.** Existe `model.User`, `service.UserService`,
> `repository.UserRepository` e a tabela `users`, mas o `UserHandler`
> (`handler/user_handler.go`) está comentado e **não é montado** no router. Os
> endpoints abaixo são um **planejamento** e não respondem hoje.

Estrutura planejada (CRUD), seguindo as convenções da API:

## User [/user] (planejado)

### Criar usuário [POST] (planejado)

+ Attributes (User)

### Buscar usuário por ID [GET] (planejado)

### Atualizar usuário [PUT] (planejado)

### Deletar usuário [DELETE] (planejado)

# Data Structures

## HealthResponse (object)
+ status: `ok` (string, required)
+ message: `Service is healthy` (string, required)

## Notebook (object)
+ notebook_id: `9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90` (string, required) - UUID.
+ user_id: `3fa85f64-5717-4562-b3fc-2c963f66afa6` (string, required) - UUID.
+ icon: `📓` (string, optional)
+ name: `Estudos de Go` (string, required)
+ image: `https://cdn.exemplo.com/capa.png` (string, optional)
+ description: `Tudo sobre a linguagem Go` (string, optional)
+ deleted_at: `null` (string, nullable) - Timestamp RFC3339 ou `null` se ativo.
+ created_at: `2026-06-18T11:45:00Z` (string, required) - Timestamp RFC3339.
+ updated_at: `2026-06-18T11:45:00Z` (string, required) - Timestamp RFC3339.

## Tag (object)
+ tag_id: `1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d` (string, required) - UUID.
+ name: `Importante` (string, required)
+ color: `FF5733` (string, required) - Cor HEX sem `#` (até 6 caracteres).
+ user_id: `3fa85f64-5717-4562-b3fc-2c963f66afa6` (string, required) - UUID.
+ deleted_at: `null` (string, nullable)
+ created_at: `2026-06-18T11:45:00Z` (string, required)
+ updated_at: `2026-06-18T11:45:00Z` (string, required)

## MetaContent (object)
+ content_id: `5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c` (string, required) - UUID.
+ notebook_id: `9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90` (string, required) - UUID.
+ user_id: `3fa85f64-5717-4562-b3fc-2c963f66afa6` (string, required) - UUID.
+ icon: `📄` (string, optional)
+ name: `Introdução` (string, required)
+ deleted_at: `null` (string, nullable)
+ created_at: `2026-06-18T11:45:00Z` (string, required)
+ updated_at: `2026-06-18T11:45:00Z` (string, required)

## NodesContent (object)
+ node_id: `2f7b3c4d-5e6a-4b1c-9d8e-3a2b1c0d9e8f` (string, required) - UUID.
+ content_id: `5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c` (string, required) - UUID.
+ user_id: `3fa85f64-5717-4562-b3fc-2c963f66afa6` (string, required) - UUID.
+ notebook_id: `9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90` (string, required) - UUID.
+ deleted_at: `null` (string, nullable)
+ created_at: `2026-06-18T11:45:00Z` (string, required)
+ updated_at: `2026-06-18T11:45:00Z` (string, required)

## MetaTagContent (object)
+ content_id: `5e2c9a1b-7d3f-4e6a-8b2c-1f0d9e8a7b6c` (string, required) - UUID. Parte da PK composta.
+ tag_id: `1d4a8f6c-2b3e-4d5a-9c7b-0e1f2a3b4c5d` (string, required) - UUID. Parte da PK composta.
+ notebook_id: `9c1b6d2e-4f7a-4c3b-8e1d-2a5f6b7c8d90` (string, required) - UUID.
+ user_id: `3fa85f64-5717-4562-b3fc-2c963f66afa6` (string, required) - UUID.
+ deleted_at: `null` (string, nullable)
+ created_at: `2026-06-18T11:45:00Z` (string, required)
+ updated_at: `2026-06-18T11:45:00Z` (string, required)

## User (object)
+ user_id: `3fa85f64-5717-4562-b3fc-2c963f66afa6` (string, required) - UUID.
+ name: `Lucas` (string, required)
+ email: `lucas@exemplo.com` (string, required)
+ phone: `+55 11 99999-0000` (string, optional)
+ deleted_at: `null` (string, nullable)
+ created_at: `2026-06-18T11:45:00Z` (string, required)
+ updated_at: `2026-06-18T11:45:00Z` (string, required)
