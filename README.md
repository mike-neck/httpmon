httpmon
---

`httpmon` is a synthetic monitoring command line tool, written by golang.

usage
---

#### An example of success case.

```shell-session
$ httpmon https://example.com/api --expect-status 200 --timeout 5s
ok
```

#### An example of failure case.

```shell-session
$ httpmon https://example.com/api --expect-status 200 --timeout 5s
failure
  expect status: 200
     got status: 503
```

supported methods
---

* [ ] GET
* [ ] POST
* [ ] PUT
* [ ] DELETE
* [ ] PATCH
