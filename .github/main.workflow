workflow "Build and Deploy" {
  on = "push"
  resolves = [
    "Deploy Google Cloud Function",
  ]
}

action "Deploy branch filter" {
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "Google Cloud Authenticate" {
  needs = ["Deploy branch filter"]
  uses = "actions/gcloud/auth@master"
  secrets = ["GCLOUD_AUTH"]
}

action "Deploy Google Cloud Function" {
  needs = ["Google Cloud Authenticate"]
  uses = "actions/gcloud/cli@master"
  secrets = ["ENV_247964_DOMAIN", "ENV_247964_PASSWORD", "ENV_247964_USERID", "ENV_TOKEN", "ENV_DNSIMPLE_TOKEN"]
  args = "functions deploy ${NAME} --project=${PROJECT} --entry-point=${ENTRY_POINT} --runtime=${RUNTIME} --trigger-http --memory=${MEMORY} --region=${REGION} --set-env-vars=247964_DOMAIN=${ENV_247964_DOMAIN},247964_PASSWORD=${ENV_247964_PASSWORD},247964_USERID=${ENV_247964_USERID},TOKEN=${ENV_TOKEN},DNSIMPLE_TOKEN=${ENV_DNSIMPLE_TOKEN}"
  env = {
    NAME = "dnsimple-dk-hostmaster-ds-upload"
    PROJECT = "reload-internal-alpha"
    ENTRY_POINT = "Handle"
    RUNTIME = "go111"
    REGION = "europe-west1"
    MEMORY = "128M"
  }
}
