restart:
	docker compose down
	docker compose up -d

stop:
	docker compose down

generate-openapi:
	@oapi-codegen --version > /dev/null 2>&1 || go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0
	
	@mkdir -p internal/users/handler/generated
	oapi-codegen -package generated -generate types api/openapi/users.yaml > internal/users/handler/generated/openapi_types.gen.go

	@mkdir -p internal/trainings/handler/generated
	oapi-codegen -package generated -generate types api/openapi/trainings.yaml > internal/trainings/handler/generated/openapi_types.gen.go

generate-grpc:
	@mkdir -p internal/trainer/handler/generated
	protoc ./api/protobuf/*.proto --go_out=. --go-grpc_out=.
