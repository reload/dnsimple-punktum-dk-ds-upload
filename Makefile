.PHONY: test deploy logs clean doc

NAME=dnsimple-dk-hostmaster-ds-upload
ENTRY_POINT=Handle
REGION=europe-west1
RUNTIME=go111

export GO111MODULE=on

doc: README.md

README.md: *.go .godocdown.tmpl
	godocdown --output=README.md

test: *.go
	go test ./...

deploy: test
	gcloud functions deploy $(NAME) --entry-point=$(ENTRY_POINT) --runtime=$(RUNTIME) --trigger-http --memory=128M --region=$(REGION)

logs:
	gcloud functions logs read $(NAME) --region=$(REGION) --limit=100
