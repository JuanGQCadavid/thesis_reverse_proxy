
## How to find data?

1. Open wireshark
2. set the next query

```
tcp.port == 9000 || udp.port == 9000 # for questDB 
tcp.port == 8812 || udp.port == 8812 # for questDB postgress connection 
tcp.port == 5432 || udp.port == 5432 # for Postgres 
```

## The data I found so far


```
1652	741.404288	::1	::1	HTTP	770	POST /write?precision=n HTTP/1.1  (text/plain)

HEADERS
authorization: Basic YW5kbWVzaXNlc3RhamE6YW5kbWUxNGVudHJ5OTFodHI=\r\n
accept: */*\r\n

BODY: content-type: text/plain; charset=utf-8\r\n
measurement,measurement_type=pedestrian,series_type=pedestrian-out,device_identity=electricity_counter_mäe_13_2,api_key=admin_key unit="count",value=99i 1738728016500000000\n
measurement,measurement_type=pedestrian,series_type=pedestrian-in,device_identity=electricity_counter_mäe_13_2,api_key=admin_key unit="count",value=77i 1738728016500000000\n

1657	741.527686	::1	::1	HTTP	153	HTTP/1.1 204 OK

1625	741.191964	::1	::1	HTTP/JSON	340	HTTP/1.1 200 OK , JSON (application/json)
Weird, this is for 

1621	741.119455	::1	::1	HTTP	171	GET /settings HTTP/1.1
```

## Events

```wireshark


[…] avc_VehicleEvent,API_KEY=admin_key,device=36a783e5_230f_44e5_8fbf_b2c8f9ef7580,fragment_type=avc_Vehicle text="Vehicle detected",vehicle_class=2i,occupancy=660i,gap=12220i,length=12.9,lane_id=2i,vehicle_id=4706611i,speed=70i,test=1i 


[…] avc_VehicleEvent,API_KEY=admin_key,device=36a783e5_230f_44e5_8fbf_b2c8f9ef7580,fragment_type=avc_Vehicle text="Vehicle detected",vehicle_class=2i,occupancy=660i,gap=12220i,length=12.9,lane_id=2i,vehicle_id=4706611i,speed=70i,test=1i 

```


## HTTP

```http request
# This is the status line
POST /measurement/measurements HTTP/1.1

# This is call fields
Content-Type: application/json
Authorization: Basic YWRtaW5fa2V5OjMxMjk=
# Client programs that initiate the request https://www.rfc-editor.org/rfc/rfc9110.html#name-user-agents
User-Agent: PostmanRuntime/7.51.1
Accept: */*
Postman-Token: 1b3e276b-6ee8-429c-a6bd-7cc870a1b264

# This could be even not needed at all 
# https://www.rfc-editor.org/rfc/rfc9110.html#section-7.1-3 
# Never the less it could be good to be set up as the Origin Server URI-host:Port
Host: localhost:8080

# I should make sure this is crated by the proxy too.
Connection: keep-alive

# He he he, lets add also the Via
# https://www.rfc-editor.org/rfc/rfc9110.html#name-via
# Via: 1.0 fred, 1.1 p.example.net

Accept-Encoding: gzip, deflate, br

# Should I also add encodign?
# https://www.rfc-editor.org/rfc/rfc9110.html#name-content-encoding
# Well maybe: https://github.com/questdb/questdb/pull/6165
# THE FUCK! It was released this week!!!!
# https://github.com/questdb/questdb/releases/tag/9.3.3
# https://github.com/questdb/questdb/releases/tag/9.1.0
# https://stackoverflow.com/questions/59494383/how-can-you-tell-if-a-pr-on-github-thats-merged-has-made-it-into-a-release
# https://github.com/questdb/questdb/commit/7eb19b9439ceab3f34e87535f30675dca7521a57#diff-6fd7374e50bb1b6c72158ad48871a01bb9012bcefd4153836f0e64204d22ac0aR46-R131
# THE FUCKKKKK!!!! I'm so fucking smart bro
# https://stackoverflow.com/questions/23086488/how-to-gzip-compress-an-http-request-in-go
Content-Encoding: gzip

# I should also check this https://www.rfc-editor.org/rfc/rfc9110.html#section-8.8.3.3-8


# https://www.rfc-editor.org/rfc/rfc9110.html#name-content-length
# THis is only for the content 
# Only for POST/PUT
Content-Length: 743 

# There is always an empty line btw the Headers/represnetative headers /fields and the content /body
{
    "measurements": [
        {
            "type": "pedestrian-counter",
            "source": {
                "id": "128498"
            },
            "time": "2025-02-05T04:00:16.500Z",
            "pedestrian": {
                "pedestrianOut": {
                    "unit": "pc",
                    "value": 12
                }
            }
        },
        {
            "type": "pedestrian-counter",
            "source": {
                "id": "128498"
            },
            "time": "2025-02-05T05:00:16.500Z",
            "pedestrian": {
                "pedestrianOut": {
                    "unit": "pc",
                    "value": 33
                }
            }
        }
    ]
}

```

# Tha fucking response

```http request
HTTP/1.1 201 Created
date: Sun, 01 Mar 2026 09:37:51 GMT
server: uvicorn
content-length: 48
content-type: application/json

{"message":"Added measurements to the database"}
```


```http request

GET /settings HTTP/1.1
user-agent: questdb/python/4.1.0
accept: */*
host: localhost:9000

POST /write?precision=n HTTP/1.1
content-length: 332
user-agent: questdb/python/4.1.0
accept: */*
host: localhost:9000
content-type: text/plain; charset=utf-8
authorization: Basic YW5kbWVzaXNlc3RhamE6YW5kbWUxNGVudHJ5OTFodHI=

measurement,measurement_type=pedestrian-counter,series_type=pedestrian:pedestrianOut,device_identity=128498,api_key=admin_key unit="pc",value=12i 1738728016500000000
measurement,measurement_type=pedestrian-counter,series_type=pedestrian:pedestrianOut,device_identity=128498,api_key=admin_key unit="pc",value=33i 1738731616500000000


```


https://questdb.com/docs/ingestion/ilp/overview/


```http request

HTTP/1.1 200 OK
Server: questDB/1.0
Date: Sun, 1 Mar 2026 09:37:52 GMT
Transfer-Encoding: chunked
Content-Type: application/json

73
{"release.type":"OSS","release.version":"8.2.1","acl.enabled":false,"posthog.enabled":false,"posthog.api.key":null}
00

HTTP/1.1 204 OK
# https://www.rfc-editor.org/rfc/rfc9110.html#name-server
Server: questDB/1.0
Date: Sun, 1 Mar 2026 09:37:52 GMT


```