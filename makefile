.PHONY: protogen
protogen:
	@./scripts/protogen.sh bots api/grpc/gen

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh bots internal/bots/ports/httpport httpport
