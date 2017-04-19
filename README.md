# Nudity check server

[![Go Report Card](https://goreportcard.com/badge/github.com/tb0hdan/nudity_check_server)](https://goreportcard.com/report/github.com/tb0hdan/nudity_check_server)

HTTP Server that accepts image URL and returns nudity status


```bash
cd nudity_check_server
go get github.com/koyachi/go-nude
go build
./nudity_check_server
```

then

```bash
curl -X GET 'http://localhost:8000/?u=aHR0cDovLzY4Lm1lZGlhLnR1bWJsci5jb20vN2VlNThiOTM2MGU1YzA0MTIxOTQ4ODJiOWI0ZDNmOTYvdHVtYmxyX251bnNxMUljck0xdTI2eDJvbzFfMTI4MC5qcGc='
```

response:

```json
{"isNude": "true", "status": "ok"}
```

# Building query string

WARNING: Server doesn't support image URLs with query strings intentionally to minimize cache poisoning

## Bash

```bash
echo -n 'http://example.com/path/to/image.jpg'|base64
```

or even

```bash
echo 'http://example.com/path/to/image.jpg'|base64
```

## Python

```
>>> import base64
>>> base64.urlsafe_b64encode('http://example.com/other/path/to/image.jpg')
'aHR0cDovL2V4YW1wbGUuY29tL290aGVyL3BhdGgvdG8vaW1hZ2UuanBn'
>>>
```
