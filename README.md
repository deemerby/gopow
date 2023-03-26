# gopow
Design and implement “Word of Wisdom” tcp server

## Content

- [Content](#Content)
- [What is this?](#What-is-this)
- [Getting Started](#Getting-Started)
  - [Challenge-response Protocol](#Challenge-response-Protocol)
  - [Proof of Work](#Proof-of-Work)
    - [Selection of PoW Algorithm](#Selection-of-PoW-Algorithm)
  - [Local Docker Usage](#Local-Docker-Usage)
  - [Local Usage](#Local-Usage)
  - [Flags](#Flags)
  - [Versioning](#Versioning)

  ------------------------------------------------------------------------------------------

## What is this?

This TCP server with protection from DDOS attacks with the Prof of Work
(https://en.wikipedia.org/wiki/Proof_of_work).

The POW algorithm is Hashcash.
(https://en.wikipedia.org/wiki/Hashcash).

  ------------------------------------------------------------------------------------------

### Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.
+ TCP server should be protected from DDOS attacks with the Prof of Work
(https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
+ The choice of the POW algorithm should be explained.
+ After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
+ Docker file should be provided both for the server and for the client that solves the POW challenge.

  ------------------------------------------------------------------------------------------

## Getting Started

These instructions will get you a clone of repository and running on your local machine.
It will also explain why the hashcash algorithm was chosen to implement this service.

  ------------------------------------------------------------------------------------------

### Challenge-response Protocol

The application protocol consists of the following types:

+ MsgRequest     - client sends request to server to get a new challenge;
+ MsgChallenge   - server sends response with hashcash data to client;
+ MsgProofOfWork - client sends calculated hashcash data as a solved challenge to server;
+ MsgQuote       - server sends random quote from Word of Wisdom slice to client - if the challenge is confirmed and not expired.

### Proof of Work

The Proof of Work (PoW) algorithm is a consensus mechanism used in cryptocurrency mining. Simply put, it requires miners to solve complex mathematical problems in order to validate and record transactions on the blockchain. Each block of transactions contains a unique task, and the miner who solves it first gets the opportunity to add the block to the general chain and receive a reward in the form of new coins. The complexity of the task increases as more miners participate in the network.

#### Selection of PoW Algorithm

Hashcash has been selected, because it used in various other protocols and applications including for combating blog spam. And you may have heard of bitcoin - hashcash is the underlying mechanism used as the work function in bitcoin. The main goal of Hashcash was to minimize the receipt of large amounts of unwanted emails, using hash collision to do so. Initially, its creation was intended to combat email spam and DDoS attacks. However, in more recent times, the system became popular with Bitcoin and other cryptocurrencies, as an essential piece of the mining algorithm. Before bitcoin, SpamAssasin and Microsoft used hashcash in hotmail exchange, outlook, etc.. It is suitable for use in this task.

Hashcash has the following advantages:

+ planty of documentations and articles with description and examples;
+ hashcash is fairly easy to implement in email accounts, servers and spam filters and no central server is needed;
+ It is invisible to users;
+ It is 100% effective against spambots.

but the main disadvantage is the premise of this method is that the algorithm used will bring an objective computational cost to the CPU. If the mathematical problem that the algorithm relies on is broken, or if the computing power of the computer is qualitatively improved, the previous second-level operation becomes milliseconds or even microseconds, and all of this will become meaningless.

[Back to content](#Content)

### Local Docker Usage

The following steps will clone the `gopow` binary to your `$GOBIN` directory.

```sh
git clone git@github.com:deemerby/gopow.git
cd gopow
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


### Local Usage

```sh
git clone git@github.com:deemerby/gopow.git
cd gopow
```

You can run it as a server:
```
make server
or 
go run ./cmd/app/
```

Or run it as a client:
```
make client
or
go run ./cmd/app/ --type.server=false
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
| `hashcash.zero.cnt`       | bits number of "partial pre-image" (zero) bits in the hashed code       | No            | `"5"`                 |
| `hashcash.duration`       | maximum period after which the request will be considered expired(sec)  | No            | `"20"`                |
| `hashcash.max.iteration`  | the maximum iteration to calulate haskcash                              | No            | `"1000000000"`        |
| `max.random.number`       | the maximum random number in the given range                            | No            | `"100000"`            |
| `logging.level`           | log level of application                                                | No            | `"info"`              |

Possible log levels: `"panic"`, `"fatal"`, `"error"`, `"warn"`, `"warning"`, `"info"`, `"debug"`, `"trace"`

[Back to content](#Content)


### Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com:deemerby/gopow.git/tags).

[Back to content](#Content)

  ------------------------------------------------------------------------------------------