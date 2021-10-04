# go-surepetapi

## Status

This was a weekend project / proof of concept. This code was working at the
time of commit but is not production quality and is not actively maintained.

## About

This repo contains three interesting tools:

* `github.com/dansimau/go-surepetapi`: Client interface for Surepet cloud API,
  partial implementation.
* `github.com/dansimau/go-surepetapi/cmd/surepet`: Command line interface for
  Surepet cloud API. Partial implementation.
* `github.com/dansimau/go-surepetapi/cmd/hksurepet`: Homekit bridge for Surepet
  cat flap. Creates a switch to manage lock status.

## Configuration

Create a config file called `surepet.yaml` and place it in your current (or a
parent) directory:

```lang=yaml
api:
  authToken:
  emailAddress:
  password:
```

You can fill either `emailAddress` and `password` OR run `surepet token` to
generate an auth token and fill out `authToken` instead.
