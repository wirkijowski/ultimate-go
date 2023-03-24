SHELL := /bin/bash

# ==============================================================================
# Testing running system

# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

admin:
	go run app/tooling/admin/main.go
	
VERSION := 1.0

all: sales-api

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := ardan-starter-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.24.0@sha256:0866296e693efe1fed79d5e6c7af8df71fc73ae45e3679af05342239cdc5bc8e \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/sales-api-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/sales-api-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-sales:
	kubectl get pods -o wide --watch

kind-logs:
	kubectl logs -l app=sales-api --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment sales-api-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=sales-api

# ====================
# Modules support

tidy:
	go mod tidy
	go mod vendor