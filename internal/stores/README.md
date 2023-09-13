# Stores
Stores are the interface by which files are uploaded/retrieved from a storage backend.

## How to implement

Instead of forking Cavorite to add support for a new storage backend, the recommended strategy is to implement a plugin. See [the localstore plugin](../../plugin/localstore/) as an example.

### What should a plugin's Upload/Retrieve functions do?

## Upload
When uploading files, the following items need to be handled in the implementation:
* Upload the file to the bucket

## Retrieve
When retrieving files, the following items need to be handled in the implementation:

* Read the metadata file
* Make sure the local path exists where the file will reside based on the metadata
* Download the file into the given location
* Confirm the hash of the local file matches the hash in the metadata, Error if it does not
* Delete the file if the hashes do not match
