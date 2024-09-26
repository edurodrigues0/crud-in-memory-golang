# CRUD IN MEMORY GOLANG

## Apresentação

Neste projeto desafio é criar uma API RESTful que irá realizar operações CRUD in-memory.

## Como Rodar

### Pré-requisitos

- **Go** instalado na máquina.

### Passos para Configuração e Execução

1. **Clone o repositório:**

```bash
git clone git@github.com:edurodrigues0/crud-in-memory-golang.git
# Naveguete até o repositorio
```

5. **Acesse a aplicação:**
O backend estará rodando na porta 8080 usando o comando ``go run ./cmd/api/main.go``. Você pode acessar a API através de http://localhost:8080.

## Endpoints da API

### Create User
- Método: POST
- Endpoint: /api/users
- Descrição: Permite criar um usuario in memory.
- Body: "firstname", "lastname", "biography"

### Get User
- Método: GET
- Endpoint: /api/users/:id
- Descrição: Permite buscar usuario pelo id
- Parâmetros: userID tipo uuid na url param

### Get Users
- Método: GET
- Endpoint: /api/users
- Descrição: Permite listar os usuarios in memory.

### Update User
- Método: PUT
- Endpoint: /api/users/:id
- Descrição: Permite editar um usuario in memory.
- Body: "firstname", "lastname", "biography"

### Delete User
- Método: DELETE
- Endpoint: /api/users/:id
- Descrição: Permite deletar um usuario in memory.