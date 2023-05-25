# Stores
Stores are the interface by which metadata is managed and files are uploaded/retrieved.

## Upload
When uploading files, the following items should be handled in the implementation:

* Metadata should be generated from the file and written alongside the binary on disk
* Upload the file to the bucket

## Retrieve
When retrieving files, the following items need to be implemented:

* Read the metadata file
* Make sure the local path exists where the file will reside based on the metadata
* Download the file into the given location
* Confirm the hash of the local file matches the hash in the metadata, Error if it does not
* Delete the file if the hashes do not match
