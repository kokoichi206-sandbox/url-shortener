# url-shortener

## Usage

### Access to shortened URL

``` sh
$ curl http://localhost:8080/google
```

### Generate shortened URL

``` sh
$ curl -X POST http://localhost:8080/api/v1/urls -H 'Content-Type: application/json' -d '{"original_url":"https://github.com/kokoichi206"}'
{"short_url":"mRJ"}

$ curl -v http://localhost:8080/mRJ
```
