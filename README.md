# terraform-provider-megaport

*This project is a work in progress*

## Utilities

To retrieve a new token for the megaport api:
```
$ eval $(make reset-token)
```

Alternatively, you can use the helper tool directly:
```
$ cd util/megaport_token
$ go run .
```
To revoke the current token (and get a new one) you can pass the `--reset` flag.
