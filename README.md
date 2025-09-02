# Microsserviço de Autenticação para E-commerce

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)
![Docker](https://img.shields.io/badge/Docker-20.10-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)

## 📖 Sobre o Projeto

Este projeto consiste em um microsserviço de autenticação e gerenciamento de usuários, desenvolvido em Go como parte de um sistema de e-commerce simplificado. O objetivo principal é exercitar conceitos de arquiteturas distribuídas, como comunicação entre serviços, contratos de API, segurança e persistência de dados.

O serviço é totalmente containerizado com Docker, utiliza PostgreSQL para persistência de dados e `golang-migrate` para o versionamento do schema do banco de dados. A ênfase é em práticas profissionais, incluindo uma arquitetura limpa em camadas, um robusto tratamento de erros e uma suíte de testes abrangente.

### ✨ Funcionalidades Principais
* **Cadastro de Usuários:** Endpoint público para criação de novas contas.
* **Autenticação com JWT:** Geração de JSON Web Tokens no login para autenticação stateless.
* **Gerenciamento de Perfil:** Endpoint protegido para consulta de dados do usuário autenticado.
* **Validação Centralizada de Token:** Endpoint interno para que outros microsserviços possam validar tokens.
* **Segurança Serviço-a-Serviço:** Endpoints internos protegidos por API Key.
* **Tratamento de Erros Estruturado:** A API retorna erros em formato JSON com códigos padronizados para facilitar a integração com clientes.
* **Qualidade e Segurança Automatizadas:** Integração com `golangci-lint` (linting), `govulncheck` (análise de vulnerabilidades) e `gitleaks` (detecção de segredos) via `Makefile`.
* **Suíte de Testes Abrangente:** Testes de unidade e integração para garantir a qualidade e a confiabilidade do código.

## 🛠️ Arquitetura e Tecnologias

O projeto segue uma arquitetura em camadas para uma clara separação de responsabilidades (API, Lógica de Negócio, Repositório).

### Tecnologias Utilizadas
* **Linguagem:** Go
* **Banco de Dados:** PostgreSQL
* **Containerização:** Docker & Docker Compose
* **Roteador HTTP:** Chi
* **Migrations:** golang-migrate
* **Automação:** Makefile
* **Testes:** Ginkgo & Gomega, `ory/dockertest`, `stretchr/testify`

### Estrutura de Diretórios

<img width="330" height="884" alt="image" src="https://github.com/user-attachments/assets/52ed3dc3-c7c2-41fb-8307-a35df66d0a14" />

## 📜 Documentação da API

A API utiliza um formato JSON estruturado para respostas de erro.

### Respostas de Erro
Todas as respostas de erro (status `4xx` ou `5xx`) seguem o formato abaixo:
```json
{
  "code": "CODIGO_DO_ERRO",
  "message": "Uma mensagem descritiva do erro."
}
```

**Códigos de Erro Comuns:**

| Status HTTP | Código (`code`) | Descrição |
| :--- | :--- | :--- |
| `400 Bad Request` | `INVALID_REQUEST_BODY` | O corpo da requisição é inválido ou malformado. |
| `400 Bad Request` | `INVALID_INPUT` | Um ou mais campos são inválidos (ex: senha muito curta). |
| `401 Unauthorized`| `INVALID_CREDENTIALS` | E-mail ou senha incorretos. |
| `404 Not Found` | `USER_NOT_FOUND` | O usuário solicitado não foi encontrado. |
| `409 Conflict` | `EMAIL_ALREADY_EXISTS` | O e-mail fornecido no cadastro já está em uso. |
| `500 Internal Server Error` | `INTERNAL_SERVER_ERROR` | Ocorreu uma falha inesperada no servidor. |

### Endpoints

### `POST /register`
* **Descrição:** Cadastra um novo usuário.
* **Autenticação:** Nenhuma
* **Corpo:** `{ "name": "string", "email": "string", "password": "string" }`

### `POST /login`
* **Descrição:** Autentica um usuário e retorna um token JWT. 
* **Autenticação:** Nenhuma
* **Corpo:** `{ "email": "string", "password": "string" }`

### `GET /profile`
* **Descrição:** Retorna o perfil do usuário autenticado. 
* **Autenticação:** JWT Obrigatória (`Authorization: Bearer <token>`)

### `POST /auth/validate`
* **Descrição:** (Uso Interno) Valida um token JWT para outros serviços. 
* **Autenticação:** API Key Interna (`X-Internal-Api-Key: <chave>`)
* **Corpo:** `{ "token": "string" }`

## 🧪 Testes
O projeto adota uma estratégia de testes híbrida para garantir a máxima qualidade e confiança.

### Testes de Unidade
Pacotes: `jwt` e `api`.

Objetivo: Testar a lógica pura e as regras de negócio em completo isolamento, sem a necessidade de banco de dados ou outros serviços externos. A camada api é testada usando mocks para simular o comportamento da camada de serviço.

### Testes de Integração

**Pacotes**: `repository` e `service`.

**Objetivo**: Garantir que a interação com o banco de dados PostgreSQL funcione como o esperado. Estes testes utilizam a biblioteca 

`ory/dockertest` para iniciar um container Docker com um banco de dados real e limpo para cada execução.

**Seeders & Stubs**: Dados de teste são populados usando seeders, e as entidades são geradas com stubs para garantir consistência e realismo.

### Como Rodar os Testes
O `Makefile` já inclui um comando para executar toda a suíte de testes (unidade e integração).

```bash
make test
```

## 🚀 Como Executar o Projeto

Siga os passos abaixo para colocar o ambiente de desenvolvimento no ar.

### Pré-requisitos
* [Go](https://go.dev/doc/install) (versão 1.24+)
* [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/)
* [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Passo a Passo
1.  **Clone o repositório:**
    ```bash
    git clone <url-do-seu-repositorio>
    cd auth-service
    ```

2.  **Configure as Variáveis de Ambiente:**
    Crie um arquivo `.env` na raiz do projeto. Você pode copiar o exemplo abaixo.
    ```env
    # Docker Compose
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=authdb

    # Aplicação (URL para comunicação DENTRO do Docker)
    DATABASE_URL="postgres://postgres:postgres@db:5432/authdb?sslmode=disable"

    # Segredos
    JWT_SECRET="um-segredo-muito-forte-para-jwt"
    INTERNAL_API_KEY="uma-chave-secreta-forte-para-apis-internas"

    # Porta que a aplicação ouve DENTRO do container
    LISTEN_ADDR=":8081"
    ```

3.  **Inicie os Serviços Docker:**
    Este comando irá construir as imagens e iniciar os containers do banco de dados e da aplicação em segundo plano.
    ```bash
    make start
    ```

4.  **Aplique as Migrations:**
    Com o banco de dados no ar, crie as tabelas necessárias.
    ```bash
    make migrate-up
    ```
    Você deve ver uma mensagem de sucesso da migração `create_users_table`.

5.  **Pronto!**
    Sua aplicação está rodando e acessível em `http://localhost:8081`. Você pode acompanhar os logs com `make logs`.

## ⚙️ Comandos do Makefile

* `make start`: Inicia todos os containers em segundo plano.
* `make stop`: Para e remove todos os containers, redes e volumes.
* `make logs`: Exibe os logs do container da aplicação Go.
* `make migrate-up`: Aplica todas as migrações pendentes.
* `make migrate-down`: Reverte a última migração aplicada.
* `make create-migration`: Cria novos arquivos de migração.
* `make lint`: Roda o linter golangci-lint para análise estática do código.
* `make vulncheck`: Roda o govulncheck para buscar vulnerabilidades nas dependências.
* `make gitleaks`: Roda o gitleaks para buscar segredos commitados acidentalmente.

## 🗄️ Acesso ao Banco de Dados

Para visualizar as tabelas e dados, a forma mais fácil é usar o **Adminer**, uma interface gráfica web para bancos de dados.

1.  **Adicione o Serviço ao `docker-compose.yml`:**
    ```yaml
    # ... (dentro de 'services:')
      adminer:
        image: adminer
        container_name: auth-adminer
        restart: always
        ports:
          - "8080:8080" # Usa a porta 8080, pois a app está na 8081
    ```

2.  **Inicie o ambiente com `make start`.**

3.  **Acesse `http://localhost:8080` no seu navegador.**

4.  **Faça login com os seguintes dados:**
    * **System:** `PostgreSQL`
    * **Server:** `db` (nome do serviço do banco no Docker)
    * **Username:** `postgres`
    * **Password:** `postgres`
    * **Database:** `authdb`
