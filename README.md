# Microsserviço de Autenticação para E-commerce

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)
![Docker](https://img.shields.io/badge/Docker-20.10-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)

## 📖 Sobre o Projeto

[cite_start]Este projeto consiste em um microsserviço de autenticação e gerenciamento de usuários, desenvolvido em Go como parte de um sistema de e-commerce simplificado. [cite: 4] [cite_start]O objetivo principal é exercitar conceitos de arquiteturas distribuídas, como comunicação entre serviços, contratos de API, segurança e persistência de dados. [cite: 5]

O serviço é totalmente containerizado com Docker, utiliza PostgreSQL para persistência de dados e `golang-migrate` para o versionamento do schema do banco de dados.

### ✨ Funcionalidades Principais
* [cite_start]**Cadastro de Usuários:** Endpoint público para criação de novas contas. [cite: 32]
* [cite_start]**Autenticação com JWT:** Geração de JSON Web Tokens no login para autenticação stateless. [cite: 32]
* [cite_start]**Gerenciamento de Perfil:** Endpoint protegido para consulta de dados do usuário autenticado. [cite: 32]
* [cite_start]**Validação Centralizada de Token:** Endpoint interno para que outros microsserviços possam validar tokens. [cite: 15, 32]
* [cite_start]**Segurança Serviço-a-Serviço:** Endpoints internos protegidos por API Key. [cite: 20]

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