
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