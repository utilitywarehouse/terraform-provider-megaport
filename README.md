# terraform-provider-megaport

*This project is a work in progress*

## Utilities

To grab a token for the megaport api, you can use the helper tool:

```
$ cd util/megaport_token
$ go run .
```

To revoke a token (and get a new one) you can pass the `--reset` flag to the
tool.
