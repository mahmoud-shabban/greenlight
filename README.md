# Greenlight Movie API
This is a simple Movies api built in golang, i
uses httprouter as request router.
it is a good project to practice web concepts an
web developments
it applies the following concepts:
- Golang httprouter: instead of the defaul
servmux router
- Rate limiting: return 429 code when limi
exceeded.
- Authentication: token based authenticatio
handeled server side (can be adapted easily t
use JWT)
- Authorization: role based authorization, role
handeled server side.
- Emails for user activation: sends mails wit
activation tokens and instructions to activat
the user.
- Mailpit: uses mailpit for email testing
- Input validation: implemented custom inpu
validation logic
- Error handling: implemented custom erro
handling logic
- Graceful Shutdown 
- Logging: server structure logs.
- Database:
    * utilizing postgresql as backend db
    * using postgresql text search capabilitie to search the movies datbase based on user input.
    * database migrations: implemented databas
ordered and organized db migrations 
- Text search: using postgresql text searc
capability
- Building versioning and code quality contro
gates.
- Metrics: using expvar to expose server metric
and database metrics
- Tracing: manual instrumantation wit
opentelmetry to expose basic tracing spans
- Makefile: to automate building process
- Git pre-commit hooks: to automate the qa
checks before commiting the code to repo
- Dockerfile: to build the docker image of th
api
- Openapi Documentations: documented the ap
using openapi file.

## How to Run
1- Download docker images:
- Mailpit 
- Postgresql
- Jaeger

2- start docker containers<br>
`make docker/up`<br>

3- run database migrations<br>
`make db/migrations/up`

4- start the server<br>
`make api/run`<br><br>
You can use postman or curl to interact with api (see openapi docs at /docs)<br><br>
Access mailpit UI on port **:8025** to see email sent for user activation<br><br>

Access Jaeger UI at **:16686**


5- clean up<br>
stop the server<br>
`make docker/down`<br>
`make clean`<br>


### Read make file for furture information