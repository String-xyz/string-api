include .env.deploy

export
AWS_DEFAULT_PROFILE=${env}-string
API=string-api
ECS_CLUSTER=string-core
SERVICE_TAG=${tag}
ECR=${${env}_AWS_ACCT}.dkr.ecr.us-west-2.amazonaws.com
ECS_API_REPO=${ECR}/${API}
INTERNAL_REPO=${ECR}/admin

all: build push deploy
all-internal: build-internal push-internal deploy-internal

test-envvars:
	@[ "${env}" ] || ( echo "env var is not set"; exit 1 )
	@[ "${tag}" ] || ( echo "env tag is not set"; exit 1 )

build: test-envvars
	GOOS=linux GOARCH=amd64 go build -o ./cmd/app/main ./cmd/app/main.go
	docker build --platform linux/amd64 -t $(ECS_API_REPO):${SERVICE_TAG} cmd/app/
	rm cmd/app/main

push: test-envvars
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR)
	docker push $(ECS_API_REPO):${SERVICE_TAG}

deploy: test-envvars
	aws ecs --region $(AWS_REGION) update-service --cluster $(ECS_CLUSTER) --service ${API} --force-new-deployment

build-internal:test-envvars
	GOOS=linux GOARCH=amd64 go build -o ./cmd/internal/main ./cmd/internal/main.go
	docker build --platform linux/amd64 -t $(INTERNAL_REPO):${SERVICE_TAG} cmd/internal/
	rm cmd/internal/main

push-internal:test-envvars
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(INTERNAL_REPO)
	docker push $(INTERNAL_REPO):${SERVICE_TAG}
	
deploy-internal: test-envvars
	aws ecs --region $(AWS_REGION) update-service --cluster admin --service admin --force-new-deployment