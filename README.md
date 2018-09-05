# CST

[![Build Status](https://travis-ci.org/tsuru/cst.svg?branch=master)](https://travis-ci.org/tsuru/cst)
[![codecov](https://codecov.io/gh/tsuru/cst/branch/master/graph/badge.svg)](https://codecov.io/gh/tsuru/cst)

CST, stands for Container Security Testing, is a project to provide container security scans against many security engines (currently including only [CoreOS Clair][Clair Website]).

This project was designed to allow security scans out of the box. You would use it since at project's pipeline through where your imagination would go.

## Running CST

This section outlines the required steps to run CST anywhere. The easiest way to
deploy the CST it's using Docker Compose. Thus, you should install 
[Docker][Docker Install] and [Docker Compose][Docker Compose Install] before to
follow the instructions.

### Configuring CoreOS Clair

All configurations about CoreOS Clair is on `etc/clair.cfg.yml` file. That's a
self-explaned config file came from [Clair's repository][Clair Repository],
see more details there.

**WARNING**:
Unfortunately, Clair doesn't handle the database connection string via environment
variables yet. So, we hardcoded the database credentials on its config file.
We appreciate you should change those credentials on: `etc/clair.cfg.yml`
(line 23); and `docker-compose.yml` (envs `POSTGRES_USER`, `POSTGRES_DB` and
`POSTGRES_PASSWORD`).

### Certificate

For start the CST's web server, you will need a certificate and its private key.
Those files must be named `cert.pem` and `key.pem`, respectively, residing in
the `.certs` dir.

In a local env, you can generate a self-signed certificate running the command
below, for instance.

```bash
$ make generate-self-signed-certificate
```

### Run Docker Compose

Now, it's time to run the Docker Compose and deploy the CST's stack. Do that by
running the command below.

```bash
$ docker-compose up -d
```

Finally, you would be able to test the CST web API firing the command:

```
$ curl https://localhost:8443/health
WORKING
```

If everything is OK, you will see the "WORKING" message response.

[Clair Website]: https://coreos.com/clair/
[Clair Repository]: https://github.com/coreos/clair

[Docker Install]:  https://docs.docker.com/install/
[Docker Compose Install]: https://docs.docker.com/compose/install/
