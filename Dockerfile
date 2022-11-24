FROM ubuntu:latest

WORKDIR /app

COPY ./build/k8s-node-watcher ./

COPY ./kubeconfig ./

ENTRYPOINT ["/app/k8s-node-watcher", "--kubeconfig", "/app/kubeconfig"]
