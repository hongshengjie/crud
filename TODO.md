# TODO

- [x] (Support Context Query Exec Timeout)
- [x] (Master and Slave)
- [ ] (Trace)
- [ ] (Monitor)
- [ ] (Windows Platform Test)
- [ ] (grpc-web,front end js html) 

### mark 
> protoc -I. -I/usr/local/include --js_out=import_style=commonjs:./web --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./web proto/user.api.proto

> grpcwebproxy  --backend_addr=localhost:9000  --run_tls_server=false --allow_all_origins