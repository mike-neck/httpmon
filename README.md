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

supported methods
---

* [x] GET
* [ ] POST
* [ ] PUT
* [ ] DELETE
* [ ] PATCH
