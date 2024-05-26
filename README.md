# Cavorite

<img src="images/cavorite_logo.png" alt="drawing" width="80" background-color="transparent" align="left"/> A cli tool that makes it easy to track large, binary files in source control repositories by swapping the binary files with json metadata. Cavorite is compatible with _any_ SCM system because the binaries are tracked by json metadata files.

<br/>

## **Disclaimer**

This is not production ready nor feature complete. See [Issues](https://github.com/discentem/cavorite/issues) for future features.

Inspired by https://github.com/facebook/IT-CPE/tree/main/pantri, Cavorite is a re-write in Go with support for s3, Minio, Google Cloud Storage, and other storage systems through plugins. See [stores](stores) for information about implementing new storage drivers.

## Using Cavorite

### Minio (S3) backend

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

1. Initialize cavorite.

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

   ```json
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

### Plugin backend (arbitrary storage backends at runtime!)

> This is not yet tested automatically in Github Actions.

> Special thanks to [korylprince](https://github.com/korylprince) for implementing go-plugin support in Cavorite!

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
1. _If you are testing the localstorage plugin_: 

   Create a folder to "upload" your objects to: 

   ```bash
   mkdir ~/fake_artifact_storage
   ```

1. Create a thing you want to upload

   ```bash
   echo "i'm a blob" >> ~/some_git_repo/blob.txt
   ```

1. Initialize cavorite in `~/some_git_repo` with the plugin: 

   ```bash
   $CAVORITE_BIN init --store_type plugin --backend_address ~/fake_artifact_storage --plugin_address $CAVORITE_PLUGIN .
   ```

   This should result in a configuration file that looks something like this. Your backend_address may be slightly different.

   ```json
   {
      "store_type": "plugin",
      "options": {
         "backend_address": "/Users/bk/fake_artifact_storage",
         "plugin_address": "/Users/bk/cavorite/bazel-out/darwin_arm64-fastbuild/bin/plugin/localstore/localstore_/localstore",
         "metadata_file_extension": "cfile",
         "region": "us-east-1"
      }
   }
   ```

1. `$CAVORITE_BIN upload blob.txt`

   What happens after this depends on how the plugin is implemented but generally you should expect to see log messages that an upload was successful.

   For the example plugin in `plugin/localstore`, you'll see

   ```
   2023-09-04T16:22:35.349-0700 [INFO]  plugin.localstore: Uploading [blob.txt] via localstore plugin
   ```

1. Check that metadata was generated: `cat blob.txt.cfile`

   ```json
   {
      "name": "blob.txt",
      "checksum": "17ead08bfb84cd914e1ec5700e1d9e8a7f5e89c8517a32164a6f4cb8fcdb1901",
      "date_modified": "2023-09-25T22:25:41.783231729-07:00"
   }
   ```

1. Check that the file was uploaded.

   If you are testing the localstore plugin:

   ```bash
   cat ~/fake_artifact_storage/blob.txt 
   i'm a blob
   ```

1. `$CAVORITE_BIN retrieve blob.txt.cfile`

## Development

### Prerequisites 

Install [bazelisk](https://github.com/bazelbuild/bazelisk)

### How to build

#### with Bazel

`make`

#### with go build

`make go_build`

### Linting

`make lint`

### Unit Tests

`make test`
