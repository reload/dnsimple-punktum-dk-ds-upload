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

action "Deploy filter: master branch" {
  needs = "Build and test"
  uses = "actions/bin/filter@master"
  args = ["branch master", "not deleted"]
}

action "Authenticate to Google Cloud" {
  needs = ["Deploy filter: master branch"]
  uses = "actions/gcloud/auth@master"
  secrets = ["GCLOUD_AUTH"]
}

action "Deploy to Google Cloud Functions" {
  needs = ["Authenticate to Google Cloud"]
  uses = "actions/gcloud/cli@master"
  secrets = ["CF_NAME", "CF_PROJECT", "CF_REGION"]
  args = "functions deploy ${CF_NAME} --project=${CF_PROJECT} --region=${CF_REGION} --entry-point=${ENTRY_POINT} --runtime=${RUNTIME} --trigger-http --memory=${MEMORY} --format='yaml(status,updateTime,versionId)'"
  env = {
    ENTRY_POINT = "Handle"
    RUNTIME = "go111"
    MEMORY = "128M"
  }
}
