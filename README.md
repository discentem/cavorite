## **Disclaimer**

This is not production ready nor feature complete. See [Issues](https://github.com/discentem/cavorite/issues) for upcoming features.

## What is this?

A cli tool that makes it easy to track large, binary files in source control repositories.

Inspired by https://github.com/facebook/IT-CPE/tree/main/pantri, this is a re-write in Go with support for s3, minio, Google Cloud Storage, and possibly other storage systems. See [internal/stores](internal/stores) for information about implementing new storage drivers.

## Man page

``` 

Usage:
   [flags]
   [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize a new cavorite repo
  retrieve    retrieve a file from cavorite
  upload      Upload a file to cavorite

Flags:
      --debug   Run in debug mode
  -h, --help    help for this command
      --vv      Run in verbose logging mode

Use " [command] --help" for more information about a command.
```

## Full testing workflow

### Minio (S3)

> These steps for Minio are also performed automatically by our integration test on each pull request and push: [.github/workflows/integration-test.yaml](.github/workflows/integration-test.yaml)

1. Run minio server.
   ```shell
   make minio
   ```

1. Create needed directories.
   ```shell
   mkdir ~/some_git_project
   ```

1. Create a bucket in web console: http://127.0.0.1:9001. Let's assume the bucket is called `test`.

1. Build `cavorite`

   ```bash
   make build
   ```

   This will output and copy-and-pastable command to set `cavorite_BIN`. This will be convenient for the following steps.

1. Set cavorite_BIN environment variable. This should point to a built copy of cavorite. For example:

   ```bash
   $ cavorite_BIN=bazel-out/darwin_arm64-fastbuild/bin/cavorite_/cavorite
   ```

1. Export the credentials for minio into the env vars

   ```bash
   $ export AWS_ACCESS_KEY_ID=minioadmin
   $ export AWS_SECRET_ACCESS_KEY=minioadmin
   ```

1. Initialize cavorite. This assumes default Minio credentials. **You should change these credentials for a production deployment**.

   ```shell
   $ $cavorite_BIN init ~/some_git_project --backend_address http://127.0.0.1:9000/test  --store_type=s3 --region="us-east-1"
   ```
   If successful you should see:
   ```
   2022/10/17 23:22:14 initializing cavorite config at ~/some_git_project/.cavorite/config
   ```

   Inspect the config:
   ```shell
   $ cd ~/some_git_project
   $ cat .cavorite/config
   ```

   ```
   {
      "store_type": "s3",
      "options": {
         "backend_address": "http://127.0.0.1:9000/test",
         "metadata_file_extension": "cfile",
         "region": "us-east-1"
      }
   }
   ```

1. Upload a binary.

   The first argument (`~/Downloads/googlechromebeta.dmg`) after all the flags is the path where the object you are uploading can be found on your filesystem.

   ```shell
   $ $cavorite_BIN upload ~/some_git_project/googlechromebeta.dmg
   ```

1. Observe that the binary has been uploaded to Minio. Nagivate to http://127.0.0.1/buckets/test/browse to confirm.

1. Confirm cavorite metadata has been written.
   ```shell
   $ cat ~/some_git_project/googlechromebeta.dmg.cfile
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
   $cavorite_BIN retrieve googlechromebeta.dmg.cfile

   2022/10/18 21:57:53 type "s3" detected in cavorite "http://127.0.0.1:9000/test"
   2022/10/18 21:57:53 Retrieving [~/some_git_project/googlechromebeta.dmg]
   ```

### Test a plugin-based backend

> This is not yet tested automatically in Github Actions.

1. Compile `cavorite` with `make` and then compile some plugin. If you want to compile the example plugin, run `make with_localstore_plugin`.

1. Set environment variables for cavorite and the plugin you compiled. For example: 

```bash
CAVORITE_BIN=/Users/bk/cavorite/bazel-out/darwin_arm64-fastbuild/bin/cavorite_/cavorite
CAVORITE_PLUGIN=/Users/bk/cavorite/bazel-out/darwin_arm64-fastbuild/bin/plugin/localstore/localstore_/localstore
```

1. Change directory to a git repository (or just a new folder):

```bash
cd ~/some_git_repo
```

1. Initialize cavorite in `~/some_git_repo` with the plugin: 

```bash
$CAVORITE_BIN init --store_type plugin --backend_address $CAVORITE_PLUGIN .
```

This should result in a configuration file that looks something like this. Your backend_address will be slightly different.

```json
{
 "store_type": "plugin",
 "options": {
  "backend_address": "/Users/bk/cavorite/bazel-out/darwin_arm64-fastbuild/bin/plugin/localstore/localstore_/localstore",
  "metadata_file_extension": "cfile",
  "region": "us-east-1"
 }
}
```

1. 