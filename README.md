# Block Path

Forked from https://github.com/traefik/plugin-blockpath

[![Build Status](https://github.com/traefik/plugin-blockpath/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/plugin-blockpath/actions)

Block Path is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which sends an defined HTTP  
response when the requested HTTP path and method matches one the configured [regular expressions](https://github.com/google/re2/wiki/Syntax).

## Configuration

## Static

```toml
[pilot]
    token="xxx"

[experimental.plugins.blockpath]
    modulename = "github.com/ONPIER-OSS/traefik-blockpath"
    version = "v0.2.4"
```

## Dynamic

To configure the `Block Path` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates
and uses the `blockpath` middleware plugin to block all HTTP requests with a path starting with `/foo`. 

```toml
[http.routers]
  [http.routers.my-router]
    rule = "Host(`localhost`)"
    middlewares = ["block-foo"]
    service = "my-service"

# Block all paths starting with /foo
[http.middlewares]
  [http.middlewares.block-foo.plugin.blockpath]
    regex = ["^/foo(.*)"]
    statuscode = 403
    methods = ["GET"]

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```

To configure the `Block Path` plugin via Kubernetes Custom Resources, you should create a [middleware](https://doc.traefik.io/traefik/reference/routing-configuration/kubernetes/crd/http/middleware/) CR. The following example creates the `blockpath` middleware Custom Resource to block all HTTP requests with method `GET` and a path starting with `/foo`.

```yaml
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: block-foo
  namespace: default
spec:
  plugin:
    blockpath:
      methods:
      - GET
      regex:
      - /foo(.*)
      statuscode: 403
```

**Notes:**
* To block all HTTP methods, omit the `methods` field from the configuration.
* If no status code is specified, the plugin returns `403 Forbidden` by default.



To apply this middleware to an Ingress resource, add the following annotation:

```yaml
traefik.ingress.kubernetes.io/router.middlewares: <middleware-namespace>-<middleware-name>@kubernetescrd
```

Replace:
* `<middleware-namespace>` with the namespace where the middleware is deployed
* `<middleware-name>` with the name of your middleware resource
