# datadog apm/statsd dumper

## Install

Download appropriate binaries from https://github.com/skaji/datadog-apm-dumper/releases/latest.

## Usage

```
❯ datadog-apm-dumper
start statsd server at udp://localhost:8125
start apm server at http://localhost:8126
stat| --- 127.0.0.1:56656
stat| datadog.dogstatsd.client.packets_sent:2|c|#client:go,transport:udp
stat| datadog.dogstatsd.client.bytes_sent:2212|c|#client:go,transport:udp
apm | --- 127.0.0.1:52315
apm | [
apm |   [
apm |     {
apm |       "name": "foo",
apm |       "service": "main",
apm |       "resource": "foo",
apm |       "type": "",
apm |       "start": 1573821349240505000,
apm |       "duration": 17000,
apm |       "meta": {
apm |         "system.pid": "1264"
apm |       },
apm |       "metrics": {
apm |         "_sampling_priority_rate_v1": 1,
apm |         "_sampling_priority_v1": 1
apm |       },
apm |       "span_id": 2126654781639715089,
apm |       "trace_id": 2126654781639715089,
apm |       "parent_id": 0,
apm |       "error": 0
apm |     }
apm |   ]
apm | ]
```

## License

MIT
