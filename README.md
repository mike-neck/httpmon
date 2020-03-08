httpmon
---

`httpmon` is a synthetic monitoring command line tool, written by golang.

usage
---

#### An example of success case.

```shell-session
$ httpmon --status 200 --timeout 5s https://example.com/api
ok
```

#### An example of failure case.

```shell-session
$ httpmon --status 200 --timeout 5s https://example.com/api
status
  expect: 200
  actual: 503
test failed: 1 failed in 1 cases
```

options
---

- `-method` (`-X` for short) - HTTP Method [supporting methods](https://github.com/mike-neck/httpmon#supported-methods)
- `-request-header` - Request HTTP Header.
    - format: `[header-name]=[header-value]` (ex. `content-type=application/json`)
    - Can be specified multiple times.
- `-timeout`(`-t` for short) - Timeout for http client
    - format: numberUNIT (ex. `5m` means 5 minutes, `30s` means 30 seconds)
- `-status` (`-s` for short) - Expecting HTTP Status.
- `-response-time`(`-r` for short) - Expecting Response time within. Format is the same as timeout.
- `-expect-header` - Expecting HTTP Response header. Format is the same as request header.

supported methods
---

* [x] GET
* [ ] POST
* [ ] PUT
* [ ] DELETE
* [ ] PATCH
