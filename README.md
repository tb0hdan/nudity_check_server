# nudity_check_server
HTTP Server that accepts image URL and returns nudity status

```bash
cd nudity_check_server
go build
./nudity_check_server
```

then

```bash
curl -X GET 'http://localhost:8000/?u=aHR0cDovLzY4Lm1lZGlhLnR1bWJsci5jb20vN2VlNThiOTM2MGU1YzA0MTIxOTQ4ODJiOWI0ZDNmOTYvdHVtYmxyX251bnNxMUljck0xdTI2eDJvbzFfMTI4MC5qcGc='
```

response:

```json
{"isNude": "true"}
```
