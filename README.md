# Go API Base

## Initial setup

* Create a project Firebase
* Click database and create it
* Go to settings -> Service accounts, click "Generate new private key". This will download a JSON file. 
* Run `base64 -w 0 account.json` to get encoded version of the file
* Click gear and enable billing
* Go to https://console.cloud.google.com/ , choose your firebase project then:
  * search for "Cloud Build API" and enable it.
  * search for "Cloud Run API" and enable it.
  * In IAM & admin, choose the firebase-adminsdk service account, click the edit (pencil) and add Project Owner role.

## Set env vars

```
export G_PROJECT_ID=FIREBASE_ID
export G_KEY=BASE64_ENCODED_STRING_FROM_ABOVE
```

Easiest way is to create a file with lines like above, eg: secrets.env (be sure to .gitignore that). 

## Code

Copy [main.go](example/main.go) and start from there.

## Auto Deploying

Follow [README here](https://github.com/GoogleCloudPlatform/github-actions/tree/master/example-workflows) then
setup [workflow here](https://github.com/GoogleCloudPlatform/github-actions/tree/master/example-workflows/cloud-run).
