.PHONY: doc test check-env deploy logs post-fixture

ENTRY_POINT=Handle
RUNTIME=go116
MEMORY=128M

# Include a .env file if it exists (getting NAME, PROJECT, and REGION).
-include .env

export GO111MODULE=on

doc: README.md

README.md: *.go README.md.template
	go generate

test: *.go
	go test ./...

check-env:
	@test -n "$(NAME)" || (echo "Missing environment variable NAME" ; false)
	@test -n "$(PROJECT)" || (echo "Missing environment variable PROJECT" ; false)
	@test -n "$(REGION)" || (echo "Missing environment variable REGION" ; false)

deploy: test check-env
	gcloud functions deploy $(NAME) --project=$(PROJECT) --region=$(REGION) --entry-point=$(ENTRY_POINT) --runtime=$(RUNTIME) --trigger-http --memory=$(MEMORY)

logs: check-env
	gcloud functions logs read $(NAME) --project=$(PROJECT) --region=$(REGION) --limit=100

post-fixture: check-env
	curl -D - -X POST -H "Content-Type: application/json" -d @test/dnssec.rotation_complete.json $(shell gcloud functions describe $(NAME) --project=$(PROJECT) --region=$(REGION) --format='get(httpsTrigger.url)')?token=$(TOKEN)
