# XDS Test Client

Client allows emulating xDS connections without actual running of Envoy proxies. 

## Run
Run Dubbo CP without Dataplane tokens, debug endpoint probably also will be useful:

```shell script
DUBBO_DP_SERVER_AUTH_TYPE=none DUBBO_DIAGNOSTICS_DEBUG_ENDPOINTS=true ./build/artifacts-darwin-amd64/dubbo-cp/dubbo-cp run
```

Run XDS Test Client:

```shell script
make run/xds-client
```

## Env
- `NUM_OF_DATAPLANES` - total number of Dataplanes to emulate
- `NUM_OF_SERVICES` - total number of services to emulate
- `DUBBO_CP_ADDRESS` - address of Dubbo CP 
