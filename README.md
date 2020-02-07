# Go API Base

## Initial setup

* Create a Firebase project
* Click database and create it
* Click gear and enable billing
* Go to https://console.cloud.google.com/ , choose your firebase project then:
  * search for "Cloud Build API" and enable it.
  * search for "Cloud Run API" and enable it.
  * In IAM & admin, choose the firebase-adminsdk service account, click the edit (pencil) and add Project Owner role.
* Go to settings -> Service accounts, click "Generate new private key". This will download a JSON file. 
* Run `base64 -w 0 account.json` to get encoded version of the file

## Setup local environment

Set local env vars:

```sh
export G_PROJECT_ID=FIREBASE_ID
export G_KEY=BASE64_ENCODED_STRING_FROM_ABOVE
export G_SERVICE_NAME=example
```

Easiest way is to create a file with lines like above, eg: `secrets.env` (be sure to .gitignore that), then 
just `source secrets.env`. 

## Code

Copy [main.go](example/main.go) and start from there.

```go
go build
./example
```

## Deploy

Set cgloud project ID:

```sh
gcloud config set project $G_PROJECT_ID
```

Copy the [example Makefile](example/Makefile) and put into your project dir.

Then run:

```sh
make deploy
```

## Auto Deploying

Go to https://github.com/treeder/YOUR_REPO/settings/secrets and add all of the above env vars.

Copy the GitHub action in this repo at [.github/worksflows/main.yml](.github/worksflows/main.yml) and put
it in the same location in your repo. Commit it and push it then check the Actions tab for progress.
