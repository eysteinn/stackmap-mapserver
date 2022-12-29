BUILDVERSION:=latest
BIN:=msrv-config-runner

build:
	go build -o bin/msrv-config-runner main.go

run:
	go run bin/msrv-cinfig-runner


build-docker: build
	docker build . -t $(BIN):$(BUILDVERSION)

run-docker: build-docker
	docker run $(BIN):$(BUILDVERSION)
	#docker run -p 3000:3000 api-rest:$(BUILDVERSION)

kind-load: build-docker
	kind load docker-image $(BIN):$(BUILDVERSION)

kind-deploy: build-docker
	kind load docker-image $(BIN):$(BUILDVERSION) && kubectl apply -f deployment.yaml 

#swagger-build:
#	swagger generate spec -i ./swagger/swagger_base.yaml -o ./swagger.yaml

swagger-serve:
	cd swagger && swagger serve --flatten --port=9009 -F=swagger main.yaml