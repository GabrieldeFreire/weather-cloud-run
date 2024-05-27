# Weather App
Aplicativo de Clima
Este é um aplicativo simples de clima escrito em Go que recebe um código postal brasileiro (CEP), identifica a cidade e retorna o clima atual (temperatura em Celsius, Fahrenheit e Kelvin). O sistema foi projetado para ser implantado no Google Cloud Run.


## Pré-requisitos
- [Go](https://golang.org/dl/) 1.22 ou superior
- [Docker](https://www.docker.com/get-started)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) instalado e configurado

## Configurar variáveis de ambiente
Cria o arquivo `.env` com o comando abaixo
```bash
cp .env.example .env
```
Defina as variáveis no arquivo `.env`


## Execução local
```bash
make run
```
## Execução do Teste
```bash
make test
```

## Endpoints
### GET /?cep=<cep>
Retorna a temperatura em celsius, fahrenheit e kelvin referente ao CEP


### Cloud run 
```bash
curl -X GET "https://weather-app-q6xkm35uia-ue.a.run.app/?cep=01153000"
```
