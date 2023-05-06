## **Disclaimer**

This is not production ready or even a code complete project.

## What is this?

Inspired by https://github.com/facebook/IT-CPE/tree/main/pantri. A re-write in go with support for s3, minio and eventually other storage systems.

## Man page

```
A source control friendly binary storage system

Usage:
   [flags]
   [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize a new Pantri repo
  retrieve    retrieve a file from pantri
  upload      Upload a file to pantri

Flags:
      --debug   Run in debug mode
  -h, --help    help for this command
      --vv      Run in verbose logging mode

Use " [command] --help" for more information about a command.
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
   $ PANTRI_BIN=bazel-out/darwin_arm64-fastbuild/bin/pantri_but_go_/pantri_but_go
   ```

1. Export the credentials for minio into the env vars

   ```bash
   $ export AWS_ACCESS_KEY_ID=minioadmin
   $ export AWS_SECRET_ACCESS_KEY=minioadmin
   ```

1. Initialize pantri. This assumes default Minio credentials. **You should change these credentials for a production deployment**.

   ```shell
   $ $PANTRI_BIN init ~/some_git_project --backend_address http://127.0.0.1:9000/test  --store_type=s3 --region="us-east-1"
   ```
   If successful you should see:
   ```
   2022/10/17 23:22:14 initializing pantri config at ~/some_git_project/.pantri/config
   ```

   Inspect the config:
   ```shell
   $ cd ~/some_git_project
   $ cat .pantri/config
   ```

   ```
   {
      "store_type": "s3",
      "options": {
         "pantri_address": "http://127.0.0.1:9000/test",
         "metadata_file_extension": "pfile",
         "region": "us-east-1"
      }
   }
   ```

1. Upload a binary.

   The first argument after all the flags (`~/Downloads/googlechromebeta.dmg`) is the path where the object you are uploading can be found on your filesystem.

   ```shell
   $ $PANTRI_BIN upload ~/some_git_project/googlechromebeta.dmg
   ```

1. Observe that the binary has been uploaded to Minio. Nagivate to http://127.0.0.1/buckets/test/browse to confirm.

1. Confirm pantri metadata has been written.
   ```shell
   $ cat ~/some_git_project/googlechromebeta.dmg.pfile
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
   # Delete the file that we just uploaded
   rm ~/some_git_project/googlechromebeta.dmg

   # Then retrieve it
   $PANTRI_BIN retrieve googlechromebeta.dmg.pfile

   2022/10/18 21:57:53 type "s3" detected in pantri "http://127.0.0.1:9000/test"
   2022/10/18 21:57:53 Retrieving [~/some_git_project/googlechromebeta.dmg]
   ```

