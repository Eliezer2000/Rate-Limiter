# Rate-Limiter

Rate Limiter em Go
Objetivo
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.<br>
Descrição<br>
Este projeto implementa um rate limiter como middleware para servidores web em Go, permitindo controlar o tráfego de requisições por:<br>
Endereço IP: Limita o número de requisições de um mesmo IP em um intervalo de tempo.<br>
Token de Acesso: Limita requisições baseadas em um token único enviado no header API_KEY. O limite do token, se existir, sempre se sobrepõe ao limite por IP.<br>
As informações de controle são armazenadas no Redis, facilitando a troca futura por outro mecanismo de persistência.<br>

Requisitos<br>
Middleware injetável no servidor web<br>
Configuração de limites via variáveis de ambiente ou arquivo .env<br>
Limitação por IP e/ou por token de acesso<br>
Resposta HTTP 429 com mensagem específica ao exceder o limite<br>
Armazenamento das informações no Redis (via Docker Compose)<br>
Estratégia de persistência desacoplada (pode trocar Redis por outro)<br>
Lógica do limiter separada do middleware<br>
Servidor web na porta 8080<br>

Configuração<br>
1. Clone o repositório<br>
git clone https://github.com/Eliezer2000/Rate-Limiter.git <br>
cd rate_limiter

2. Configure o arquivo .env<br>
RATE_LIMIT_IP=10<br>
RATE_LIMIT_TOKEN=100<br>
BLOCK_TIME=300<br>
REDIS_ADDR=localhost:6379<br>
SERVER_PORT=8080<br>

3. Suba o Redis com Docker Compose<br>
docker-compose up -d<br>


4. Instale as dependências Go<br>
go mod tidy<br>

5. Rode o servidor<br>
go run main.go<br>


Como funciona<br>
Por IP:<br> Se não houver header API_KEY, o limite é aplicado ao IP do cliente.<br>
Por Token:<br> Se o header API_KEY estiver presente, o limite do token é usado (se existir), senão o limite padrão de token.<br>
Bloqueio:<br> Ao exceder o limite, o IP/token é bloqueado pelo tempo configurado em BLOCK_TIME.<br>
Resposta:<br> Durante o bloqueio, qualquer requisição retorna HTTP 429 com a mensagem:<br>

-   you have reached the maximum number of requests or actions allowed within a certain time frame

Testes<br>
Teste 1: Limitação por IP<br>
No Postman, crie uma requisição GET para http://localhost:8080/.<br>
Não adicione o header API_KEY.<br>
Clique em "Send" rapidamente várias vezes (ex: 12 vezes se o limite for 10).<br>
Resultado esperado:<br>
As primeiras requisições retornam 200 OK.<br>
Após atingir o limite, retorna 429 Too Many Requests com a mensagem configurada.<br>
Exemplo de print:<br>
-!Limite por IP - 429<br>

Teste 2: Limitação por Token<br>
No Postman, crie uma requisição GET para http://localhost:8080/.<br>
Na aba Headers, adicione:<br>
Key: API_KEY<br>
Value: abc123<br>
Clique em "Send" rapidamente várias vezes.<br>
Resultado esperado:<br>
As primeiras requisições retornam 200 OK.<br>
Após atingir o limite do token, retorna 429 Too Many Requests.<br>
Exemplo de print:<br>
!Limite por Token - 429<br>

Teste 3: Bloqueio temporário<br>
Após receber o 429, continue enviando requisições.<br>
Todas devem retornar 429 até o tempo de bloqueio acabar.<br>
Após o tempo de bloqueio (BLOCK_TIME), as requisições voltam a ser aceitas.<br>
Exemplo de print:<br>
!Bloqueio temporário<br>
Testes automatizados<br>

Troca de mecanismo de persistência<br>
A lógica de armazenamento é desacoplada via interface. Para trocar o Redis por outro mecanismo, basta implementar a interface LimiterStorage em internal/storage.


Observações<br>
O rate limiter pode ser facilmente adaptado para outros frameworks web em Go.<br>
Teste sob diferentes condições de carga para garantir robustez.<br>
Autor<br>
Seu nome - [Seu LinkedIn/GitHub]<br>
Dica:<br>
Coloque os prints reais dos testes na pasta prints/ do seu projeto e ajuste os caminhos das imagens no README.<br>
Se quiser, posso te ajudar a montar os prints ou revisar o README depois que você adicionar as imagens!<br>
7. 
8. 
