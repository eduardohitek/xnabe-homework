#EBANX - TAKE HOME ASSIGNMENT

Basic API for the EBANX Hiring Process

##How to run:

`docker-compose up`

**URL**: http://localhost:8086

**Endpoints**:

`GET /balance?account_id=`

`POST /event`
 
 `{"type":"deposit", "destination":"100", "amount":10}`
 `{"type":"withdraw", "origin":"200", "amount":10}`
 `{"type":"transfer", "origin":"100", "amount":15, "destination":"300"}`