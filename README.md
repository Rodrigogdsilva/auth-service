# Microsserviço de Autenticação para E-commerce

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)
![Docker](https://img.shields.io/badge/Docker-20.10-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)

## 📖 Sobre o Projeto

Este projeto consiste em um microsserviço de autenticação e gerenciamento de usuários, desenvolvido em Go como parte de um sistema de e-commerce simplificado. O objetivo principal é exercitar conceitos de arquiteturas distribuídas, como comunicação entre serviços, contratos de API, segurança e persistência de dados.

O serviço é totalmente containerizado com Docker, utiliza PostgreSQL para persistência de dados e `golang-migrate` para o versionamento do schema do banco de dados.

### ✨ Funcionalidades Principais
* **Cadastro de Usuários:** Endpoint público para criação de novas contas. 
* **Autenticação com JWT:** Geração de JSON Web Tokens no login para autenticação stateless. 
* **Gerenciamento de Perfil:** Endpoint protegido para consulta de dados do usuário autenticado.
* **Validação Centralizada de Token:** Endpoint interno para que outros microsserviços possam validar tokens. 
* **Segurança Serviço-a-Serviço:** Endpoints internos protegidos por API Key.

## 🛠️ Arquitetura e Tecnologias

O projeto segue uma arquitetura em camadas para uma clara separação de responsabilidades (API, Lógica de Negócio, Repositório).

### Tecnologias Utilizadas
* **Linguagem:** Go
* **Banco de Dados:** PostgreSQL
* **Containerização:** Docker & Docker Compose
* **Roteador HTTP:** Chi
* **Migrations:** golang-migrate
* **Automação:** Makefile

### Estrutura de Diretórios

<img width="580" height="408" alt="image" src="https://github.com/user-attachments/assets/513e61d5-a3e4-4d4e-b1d9-a63cad2bc380" />

## 📜 Documentação da API

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
