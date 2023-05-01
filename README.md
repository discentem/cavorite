## **Disclaimer**

This is not production ready or even a code complete project. 

## What is this? 

Inspired by https://github.com/facebook/IT-CPE/tree/main/pantri. A re-write in go with support for s3, minio and eventually other storage systems. 

## Man page

```shell

% ./pantri_but_go -h
NAME:
   pantri_but_go - pantri: but in go!

USAGE:
   pantri_but_go [global options] command [command options] [arguments...]

COMMANDS:
   init         Initalize pantri.
   upload, u    Upload the specified file
   retrieve, r  Retrieve the specified file
   delete, d    Delete the specified file
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug              Set debug to true for enhanced logging (default: false)
   --help, -h           show help (default: false)
   --source_repo value  path to source repo
```

## Full workflow

### Minio (S3)

1. Run minio server.
   ```shell
   make minio
   ```

1. Create needed directories.
   ```shell
   mkdir ~/some_git_project
   ```

1. Create a bucket in web console: http://127.0.0.1:9001. Let's assume the bucket is called `test`.

1. Build `pantri_but_go`

   ```bash
   make build
   ```

   This will output and copy-and-pastable command to set `PANTRI_BIN`. This will be convenient for the following steps.

1. Set PANTRI_BIN environment variable. This should point to a built copy of pantri_but_go. For example:

   ```bash
   PANTRI_BIN=bazel-out/darwin_arm64-fastbuild/bin/pantri_but_go_/pantri_but_go
   ```

1. Initialize pantri. This assumes default Minio credentials. **You should change these credentials for a production deployment**.

   ```shell
   AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin $PANTRI_BIN --source_repo ~/some_git_project init --pantri_address http://127.0.0.1:9000/test s3
   ```
   If successful you should see:
   ```
   2022/10/17 23:22:14 initializing pantri config at /Users/brandon/some_git_project/.pantri/config
   ```

   Inspect the config:
   ```shell
   cat ~/some_git_project/.pantri/config 
   ```

   ```
   {
      "type": "s3",
      "pantri_address": "http://127.0.0.1:9000/test",
      "options": {
         "metadata_file_extension": ".pfile",
         "remove_from_sourcerepo": false
      }
   }
   ```

1. Upload a binary. 

   `--destination` is the [s3 key](https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingObjects.html#:~:text=of%20the%20following%3A-,Key,-The%20name%20that) that the binary will be uploaded to. 

   The first argument after all the flags (`~/Downloads/googlechromebeta.dmg`) is the path where the object you are uploading can be found on your filesystem.

   ```shell
   AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin $PANTRI_BIN --source_repo ~/some_git_project upload --destination chrome/googlechromebeta.dmg ~/Downloads/googlechromebeta.dmg
   ```

1. Observe that the binary has been uploaded to Minio. Nagivate to http://127.0.0.1/buckets/test/browse to confirm.

1. Confirm pantri metadata has been written.
   ```shell
   brandon@Brandons-MacBook-Pro pantri_but_go % cat ~/some_git_project/Users/brandon/Downloads/googlechromebeta.dmg.pfile
   ```

   ```
   {
      "name": "chrome/googlechromebeta.dmg",
      "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "date_modified": "2022-10-05T10:56:17.051936728-07:00"
   }% 
   ```

1. Retrieve binaries from minio.

   ```shell
   % AWS_ACCESS_KEY_ID=minioadmin AWS_SECRET_ACCESS_KEY=minioadmin $PANTRI_BIN --source_repo ~/some_git_project retrieve Users/brandon/Downloads/googlechromebeta.dmg
   2022/10/18 21:57:53 type "s3" detected in pantri "http://127.0.0.1:9000/test"
   2022/10/18 21:57:53 Retrieving [Users/brandon/Downloads/googlechromebeta.dmg]
   ```

### Local Storage

1. Create directories and initialize pantri
   ```shell
   % mkdir ~/pantri
   % mkdir ~/some_git_project
   % $PANTRI_BIN --source_repo ~/some_git_project init --pantri_address ~/pantri
   2022/07/01 00:47:27 initializing pantri config at /Users/brandon/some_git_project/.pantri/config
   % cat /Users/brandon/some_git_project/.pantri/config
   {
      "type": "local",
      "pantri_address": "/Users/brandon/pantri",
      "options": {
         "remove_from_sourcerepo": false
      }
   }%
   ```
1. Upload a binary

   ```shell
   % cd ~/Downloads 
   % $PANTRI_BIN --source_repo ~/some_git_project upload go1.18.3.darwin-arm64.pkg          
   2022/07/01 00:48:25 type "local" detected in pantri "/Users/brandon/pantri"
   ```
1. Observe binary "uploaded" to pantri and metadata (.pfile) stashed in git repo

   ```shell
   2022/07/01 00:48:25 Uploading [go1.18.3.darwin-arm64.pkg]
   % ls ~/some_git_project 
   go1.18.3.darwin-arm64.pkg.pfile
   % ls ~/pantri
   go1.18.3.darwin-arm64.pkg	go1.18.3.darwin-arm64.pkg.pfile
   ```
1. Retrieve binaries

   ```shell
   % $PANTRI_BIN --source_repo ~/some_git_project retrieve go1.18.3.darwin-arm64.pkg
   2022/07/01 00:49:17 type "local" detected in pantri "/Users/brandon/pantri"
   2022/07/01 00:49:17 Retrieving [go1.18.3.darwin-arm64.pkg]
   % ls ~/some_git_project 
   go1.18.3.darwin-arm64.pkg	go1.18.3.darwin-arm64.pkg.pfile
   ```
