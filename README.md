# Go API Base

## Initial setup

* Create a Firebase project
* Click database and create it
* Click gear and enable billing
* Go to settings -> Service accounts, click "Generate new private key". This will download a JSON file. 
* Run `base64 -w 0 account.json` to get encoded version of the file (for secrets)
* Go to https://console.cloud.google.com/ , choose your firebase project then:
  * search for "Cloud Build API" and enable it.
  * search for "Cloud Run API" and enable it.
  * In IAM & admin, choose the firebase-adminsdk service account, click the edit (pencil) and add Project Owner role. (see below for reduced scope)

## Setup local environment

Set local env vars. Easiest way is to create a file in `secrets/dev.env` with the following (be sure to .gitignore secrets/), then `source secrets/dev.env`. 


```sh
export G_PROJECT_ID=FIREBASE_ID
export G_SERVICE_NAME=example
export G_KEY=BASE64_ENCODED_STRING_FROM_ABOVE
```


## Code

Copy [main.go](example/main.go) as a starting point.

```go
go mod init
go build
./example
```

## Deploy

Copy this [Dockerfile](example/Dockerfile) as is, no changes required.

Set cgloud project ID:

```sh
gcloud config set project $G_PROJECT_ID
```

Copy the [example Makefile](example/Makefile) and put into your project dir.

Then run:

```sh
make deploy
```

If you need other environment variables, it'll fail here, but go look at the Cloud Run interface and you'll see the service. Click it, then click Deploy New Revision, then at the bottom you'll see "Environment Variables". Add them there. You don't need to add the google ones above. 

ALSO, if the allow-unauthenticated didn't work (I've noticed this happen), click the service, go to permissions and [see this](https://cloud.google.com/run/docs/securing/managing-access?authuser=1&_ga=2.204426711.-650445000.1578069338#making_a_service_public).

## Auto Deploying

Go to https://github.com/treeder/YOUR_REPO/settings/secrets and add all of the above env vars.

Copy the GitHub action in this repo at [.github/workflows/main.yml](.github/worksflows/main.yml) and put
it in the same location in your repo. Commit it and push it then check the Actions tab for progress.

## User interface

### Deploying Static App to Firebase

REDO

### For firebase auth / google sign in

You'll need to whitelist the \*.web.app domains from firebase to use the web.app version. Go to `https://console.cloud.google.com/apis/credentials`, edit the OAuth 2.0 Client ID that says `Web client (auto created by Google Service)` and add the domains there. 

TODO: probably have to do the same on the production domains.


## Reducing firebase-adminsdk scope

Seems the following roles may be all that's needed:

```
Cloud Build Service Account
Firebase Admin SDK Administrator Service Agent
Service Account Token Creator
Cloud Run Admin
Storage Object Admin
```
