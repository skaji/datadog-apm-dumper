# datadog apm dumper

## Usage

```
❯ ./datadog-apm-dumper
Accepting connections at http://localhost:8126/

# Run datadog apm client with datadog agent apm url: http://localhost:8126/
# Then you'll get:

[
  [
    {
      "name": "foo",
      "service": "main",
      "resource": "foo",
      "type": "",
      "start": 1573821349240505000,
      "duration": 17000,
      "meta": {
        "system.pid": "1264"
      },
      "metrics": {
        "_sampling_priority_rate_v1": 1,
        "_sampling_priority_v1": 1
      },
      "span_id": 2126654781639715089,
      "trace_id": 2126654781639715089,
      "parent_id": 0,
      "error": 0
    }
  ]
]
```

## License

MIT
