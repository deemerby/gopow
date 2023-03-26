# gopow
Design and implement “Word of Wisdom” tcp server

## Content

- [Content](#Content)
- [What is this?](#What-is-this)
- [Getting Started](#Getting-Started)
  - [Local Usage](#Local-Usage)
  - [Flags](#Flags)
  - [Versioning](#Versioning)

  ------------------------------------------------------------------------------------------

## What is this?

This TCP server with protection from DDOS attacks with the Prof of Work
(https://en.wikipedia.org/wiki/Proof_of_work).


  ------------------------------------------------------------------------------------------

## Getting Started

These instructions will get you a clone of repository and running on your local machine

  ------------------------------------------------------------------------------------------

### Local Usage

The following steps will clone the `gopow` binary to your `$GOBIN` directory.

```sh
git clone git@github.com:deemerby/gopow.git
cd gopow
make
```

Also you may just pull docker image:

```sh
docker pull dlahuta/go-pow:latest
```

And then run it as a server:
```
docker run --rm dlahuta/go-pow:latest
```

Or run it as a client:
```
docker run --rm dlahuta/go-pow:latest --type.server=false 
```


[Back to content](#Content)

  ------------------------------------------------------------------------------------------

### Flags

Here's the full set of flags for application.

| Flag                      | Description                                                             | Required      | Default Value         |
| ------------------------- | ----------------------------------------------------------------------- | ------------- | --------------------- |
| `server.address`          | address of tcp server                                                   | No            | `"0.0.0.0"`           |
| `server.port`             | port of tcp server                                                      | No            | `"8084"`              |
| `type.server`             | type of application (server = true, client = false)                     | No            | `"true"`              |
| `hashcash.zero.cnt`       | bits number of "partial pre-image" (zero) bits in the hashed code       | No            | `"6"`                 |
| `hashcash.duration`       | maximum period after which the request will be considered expired(sec)  | No            | `"50"`                |
| `hashcash.max.iteration`  | github repository tag (all tags will be processed if skipped)           | No            | `"1000000000"`        |
| `max.random.number`       | github Personal Access Token for communication with gihub.api           | No            | `"100000"`            |
| `logging.level`           | log level of application                                                | No            | `"info"`              |

Possible log levels: `"panic"`, `"fatal"`, `"error"`, `"warn"`, `"warning"`, `"info"`, `"debug"`, `"trace"`

[Back to content](#Content)


### Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com:deemerby/gopow.git/tags).

[Back to content](#Content)

  ------------------------------------------------------------------------------------------