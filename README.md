# datadog apm dumper

## Usage

```
❯ ./server
2019/11/15-18:00:35 Starman::Server (type Net::Server::PreFork) starting! pid(8630)
Resolved [*]:8126 to [0.0.0.0]:8126, IPv4
Binding to TCP port 8126 on host 0.0.0.0 with IPv4
Setting gid to "20 20 20 12 61 79 80 81 98 702 701 33 100 204 250 395 398 399"
Starman: Accepting connections at http://*:8126/

# Run datadog apm client with datadog agent apm url: http://localhost:8126
# Then you'll get:

\ [
    [0] [
        [0] {
            duration    104000,
            error       0,
            meta        {
                env                "dev",
                host               "MBP-19-002",
                http.client_ip     "::1",
                http.method        "GET",
                http.query         "",
                http.referer       "",
                http.status_code   200,
                http.ua            "curl/7.54.0",
                http.url           "/foo/bar",
                system.pid         11418
            },
            metrics     {
                _dd1.sr.eausr                1,
                _sampling_priority_rate_v1   1,
                _sampling_priority_v1        1
            },
            name        "http.request",
            parent_id   0,
            resource    "GET /foo/bar",
            service     "test",
            span_id     4682563644132643360,
            start       1573808675665981000,
            trace_id    4682563644132643360,
            type        "web"
        }
    ]
]
```

## License

MIT
