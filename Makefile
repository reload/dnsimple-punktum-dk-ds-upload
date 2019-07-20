.PHONY: doc test check-env deploy logs

ENTRY_POINT=Handle
RUNTIME=go111
MEMORY=128M

# Include a .env file if it exists (getting NAME, PROJECT, and REGION).
-include .env

export GO111MODULE=on

doc: README.md

README.md: *.go .godocdown.tmpl
	godocdown --output=README.md

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
