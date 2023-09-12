# Stores
`Store` is the interface by which files are uploaded/retrieved from a storage backend.

## How to implement

Instead of forking Cavorite to add support for a new storage backend via a new implementation of `Store`, the recommended strategy is to implement a plugin. See [the localstore plugin](../../plugin/localstore/) as an example.

You must implement all the methods of `Store`.

### What should a plugin's Upload/Retrieve functions do?

## Upload
In your implemention of `Upload()`, the following tasks need to be handled:

1. Parsing `.cavorite/config` which is straightfoward:
    ```go
    import (
        "github.com/discentem/cavorite/internal/config"
    )
    // ...
    func (s *LocalStore) GetOptions() (stores.Options, error) {
        return config.LoadOptions(afero.NewOsFs())
    }
    ```

1. Uploading the file to the bucket.

    > This implementation is more complicated than GetOptions() and depends entirely on what artifact store you are using. See [internal/stores/s3](../stores/s3.go) for a detailed example.
    
    ```go
    func (s *LocalStore) Upload(ctx context.Context, objects ...string) error {
	    opts, err := s.GetOptions()
        if err != nil {
            return err
        }
        backendAddress := opts.BackendAddress
        s.logger.Info(fmt.Sprintf("Uploading %v via localstore plugin", objects))
        // call your artifact storage provider's API and pass objects
        return UploadToMagicalStore(backendAddress, ...objects)
    }
    ```

## Retrieve
When retrieving files, the following tasks need to be handled in the implementation:

1. Read the metadata file
1. Make sure the local path exists where the file will reside based on the metadata
1. Download the file into the given location
1. Confirm the hash of the local file matches the hash in the metadata, Error if it does not
1. Delete the file if the hashes do not match
