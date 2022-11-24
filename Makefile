IMAGE_REPOSITORY	?= docker-sandbox.infra.cloudera.com/apb/k8s-node-watcher
IMAGE_VERSION		?= 1.0.0-SNAPSHOT

clean:
	rm -rf build/

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED="0" go build -o build/k8s-node-watcher

docker-build: clean build
	docker build -t ${IMAGE_REPOSITORY}:${IMAGE_VERSION} .

docker-push: docker-build
	docker push ${IMAGE_REPOSITORY}:${IMAGE_VERSION}

