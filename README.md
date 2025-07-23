# 🚦 Rate-Limiter em Go

## 🎯 Objetivo

Desenvolver um **Rate Limiter** em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um **endereço IP específico** ou em um **token de acesso**.

---

## 📄 Descrição

Este projeto implementa um rate limiter como **middleware** para servidores web em Go, permitindo controlar o tráfego de requisições por:

- 🌐 **Endereço IP**: limita o número de requisições de um mesmo IP em um intervalo de tempo.
- 🔐 **Token de Acesso**: limita requisições baseadas em um token único enviado no header `API_KEY`. O limite do token, se existir, sempre se sobrepõe ao limite por IP.

🗃️ As informações de controle são armazenadas no **Redis**, facilitando a troca futura por outro mecanismo de persistência.

---

## ✅ Requisitos

- ✅ Middleware injetável no servidor web  
- ✅ Configuração de limites via `.env` ou variáveis de ambiente  
- ✅ Limitação por IP e/ou por token de acesso  
- ✅ Resposta HTTP `429` com mensagem específica ao exceder o limite  
- ✅ Armazenamento das informações no Redis (via Docker Compose)  
- ✅ Estratégia de persistência desacoplada  
- ✅ Lógica do limiter separada do middleware  
- ✅ Servidor web na porta `8080`  

---

## ⚙️ Configuração

### 1. Clone o repositório

```bash
git clone https://github.com/Eliezer2000/Rate-Limiter.git
cd rate_limiter
```

### 2. Configure o arquivo `.env`

```env
RATE_LIMIT_IP=10
RATE_LIMIT_TOKEN=100
BLOCK_TIME=300
REDIS_ADDR=localhost:6379
SERVER_PORT=8080
```

### 3. Suba o Redis com Docker Compose

```bash
docker-compose up -d
```

### 4. Instale as dependências Go

```bash
go mod tidy
```

### 5. Rode o servidor

```bash
go run main.go
```

---

## 🔍 Como funciona

### ✅ Por IP

Se não houver header `API_KEY`, o limite é aplicado ao **IP do cliente**.

### ✅ Por Token

Se o header `API_KEY` estiver presente, o limite do **token** é usado (se existir), caso contrário aplica-se o limite padrão de token.

### ⛔ Bloqueio

Ao exceder o limite, o IP/token é **bloqueado** pelo tempo configurado em `BLOCK_TIME`.

### 🔁 Resposta

Durante o bloqueio, qualquer requisição retorna:

```http
HTTP 429 Too Many Requests
you have reached the maximum number of requests or actions allowed within a certain time frame
```

---

## 🧪 Testes

### 🔹 Teste 1: Limitação por IP

1. No Postman, crie uma requisição `GET` para:

```
http://localhost:8080/
```

2. **Não adicione** o header `API_KEY`.  
3. Clique em **"Send" rapidamente várias vezes** (ex: 12 vezes se o limite for 10).

**Resultado esperado**:
- As primeiras requisições retornam `200 OK`.
- Após atingir o limite, retorna `429 Too Many Requests`.

📸 Exemplo de print:
```
📷 Limite por IP - 429
```

---

### 🔹 Teste 2: Limitação por Token

1. No Postman, crie uma requisição `GET` para:

```
http://localhost:8080/
```

2. Na aba **Headers**, adicione:

```
Key:    API_KEY
Value:  abc123
```

3. Clique em "Send" várias vezes rapidamente.

**Resultado esperado**:
- As primeiras requisições retornam `200 OK`.
- Após o limite, retorna `429 Too Many Requests`.

📸 Exemplo de print:
```
📷 Limite por Token - 429
```

---

### 🔹 Teste 3: Bloqueio temporário

1. Após receber o `429`, continue enviando requisições.
2. Todas devem retornar `429` até o tempo de bloqueio acabar.
3. Após o tempo de bloqueio (`BLOCK_TIME`), as requisições voltam a ser aceitas.

📸 Exemplo de print:
```
📷 Bloqueio temporário
```

---

## 🧪 Testes Automatizados

*Em desenvolvimento ou futuros testes poderão ser incluídos aqui.*

---

## 🔄 Troca de Mecanismo de Persistência

A lógica de armazenamento é desacoplada via interface.  
Para trocar o Redis por outro mecanismo, basta implementar a interface:

```go
internal/storage.LimiterStorage
```

---

## 📝 Observações

- 🔄 O rate limiter pode ser adaptado facilmente para outros frameworks web em Go.
- ⚙️ Teste sob diferentes condições de carga para garantir robustez.

---

## 👤 Autor

[Eliézer Alves Bernardino](https://www.linkedin.com/in/eliezeralvesbernardino/)

---

## 💡 Dica

📁 Coloque os **prints reais dos testes** na pasta `prints/` do seu projeto e ajuste os caminhos das imagens no README.

Se quiser, posso te ajudar a montar os prints ou revisar novamente depois que você adicionar as imagens!
