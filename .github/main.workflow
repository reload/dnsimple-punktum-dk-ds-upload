workflow "Build and Deploy" {
  on = "push"
  resolves = [
    "Deploy to Google Cloud Functions",
  ]
}

action "Build and test" {
  uses="cedrickring/golang-action/go1.11@master"
  args = "go test -verbose -race -cover -covermode=atomic ./..."
  env = {
    GO111MODULE = "on"
  }
}

action "Deploy filter: not a deleted branch" {
  needs = "Build and test"
  uses = "actions/bin/filter@master"
  args = "not deleted"
}

action "Deploy filter: master branch" {
  needs = "Deploy filter: not a deleted branch"
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "Authenticate to Google Cloud" {
  needs = ["Deploy filter: master branch"]
  uses = "actions/gcloud/auth@master"
  secrets = ["GCLOUD_AUTH"]
}

action "Deploy to Google Cloud Functions" {
  needs = ["Authenticate to Google Cloud"]
  uses = "actions/gcloud/cli@master"
  secrets = ["ENV_247964_DOMAIN", "ENV_247964_PASSWORD", "ENV_247964_USERID", "ENV_TOKEN", "ENV_DNSIMPLE_TOKEN"]
  args = "functions deploy ${NAME} --project=${PROJECT} --entry-point=${ENTRY_POINT} --runtime=${RUNTIME} --trigger-http --memory=${MEMORY} --region=${REGION} --set-env-vars=247964_DOMAIN=${ENV_247964_DOMAIN},247964_PASSWORD=${ENV_247964_PASSWORD},247964_USERID=${ENV_247964_USERID},TOKEN=${ENV_TOKEN},DNSIMPLE_TOKEN=${ENV_DNSIMPLE_TOKEN} --format=disable"
  env = {
    NAME = "dnsimple-dk-hostmaster-ds-upload"
    PROJECT = "reload-internal-alpha"
    ENTRY_POINT = "Handle"
    RUNTIME = "go111"
    REGION = "europe-west1"
    MEMORY = "128M"
  }
}
