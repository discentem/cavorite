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

## Full workflow (for local storage)

1. Create directories and initialize pantri
   ```shell
   % mkdir ~/pantri
   % mkdir ~/some_git_project
   % ./pantri_but_go --source_repo ~/some_git_project init --pantri_address ~/pantri
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
   % ~/pantri_but_go/pantri_but_go --source_repo ~/some_git_project upload go1.18.3.darwin-arm64.pkg          
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
   % ~/pantri_but_go/pantri_but_go --source_repo ~/some_git_project retrieve go1.18.3.darwin-arm64.pkg
   2022/07/01 00:49:17 type "local" detected in pantri "/Users/brandon/pantri"
   2022/07/01 00:49:17 Retrieving [go1.18.3.darwin-arm64.pkg]
   % ls ~/some_git_project 
   go1.18.3.darwin-arm64.pkg	go1.18.3.darwin-arm64.pkg.pfile
   ```

