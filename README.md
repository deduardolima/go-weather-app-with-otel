# Go-Weather-App-With-Otel

## Descrição

Este projeto consiste em dois serviços (A e B) que trabalham juntos para receber um CEP, identificar a cidade e retornar o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin) juntamente com o nome da cidade. O sistema implementa OpenTelemetry (OTEL) e Zipkin para tracing distribuído.

## Visão Geral

- **Serviço A**:
  - Recebe um CEP via POST.
  - Valida o CEP.
  - Encaminha a solicitação para o Serviço B.

- **Serviço B**:
  - Recebe um CEP válido.
  - Identifica a cidade correspondente.
  - Retorna a temperatura atual em graus Celsius, Fahrenheit e Kelvin juntamente com o nome da cidade.

## Funcionalidades

- **OpenTelemetry**: Implementação de tracing distribuído para monitoramento das transações entre serviços.
- **Zipkin**: Utilizado para visualização dos traces coletados pelo OpenTelemetry.




## Estrutura do Projeto


- go-weather-with-otel/
  - .docker/
   - otel-collector-config.yaml
  - service-a/
    - cmd/
      - main.go
      - .env
    - internal/
      - handler.go
      - handler_test.go
    - Dockerfile    
    - go.mod
    - go.sum
  - service-b/
    - cmd/
      - main.go
      - .env    
    - internal/
      - location/
        - client.go
        - service.go
    - weather/ 
      - client.go 
      - handler_test.go
      - handler.go
    - Dockerfile    
    - go.mod
    - go.sum  
  - docker-compose.yml
  - README.md


## Pré-requisitos

- Docker
- Docker Compose

## Configuração
```
git clone https://github.com/deduardolima/go-weather-app-with-otel.git
cd go-weather-app

```

No arquivo `.env` inclua suas credenciais para a API WeatherAPI. Se você ainda não tiver uma chave de API, crie uma conta para obter acesso em:
[WEATHER API](https://www.weatherapi.com/)

```
WEATHER_API_KEY=SUA-API-KEY-AQUI
PORT=8080

```

## Instalação e Execução com Docker
Construa e inicie os containers:
```
docker-compose up --build -d
```

isso irá construir a imagem do aplicativo e iniciar o serviço definido no docker-compose.yml



## Execução 

Faça uma requisição POST para o serviço A:

```sh
curl -X POST http://localhost:8080/input -H "Content-Type: application/json" -d '{"cep":"80010100"}'
```


Visualize os spans no Zipkin:
[http://localhost:9411](http://localhost:9411/)



## Exemplo de Resposta

### Em caso de sucesso:

```json
{
  "temp_C": 19,
  "temp_F": 66.2,
  "temp_K": 292.15
}
```

### Em caso de falha, quando o CEP não é válido (com formato correto):

```json
{
  "error": "invalid zipcode"
}
```

### Em caso de falha, quando o CEP não é encontrado:

```json
{
  "error": "can not find zipcode"
}
```

## Referências

- [Go](https://golang.org/doc/)
- [Docker](https://docs.docker.com/)
- [OpenTelemetry](https://opentelemetry.io/docs/)
- [Zipkin](https://zipkin.io/)

## Créditos

Este projeto foi criado por [Diego Eduardo](http://github.com/deduardolima)







