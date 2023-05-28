# Stores
Stores are the interface by which metadata is managed and files are uploaded/retrieved.

## Upload
When uploading files, the following items should be handled in the implementation:

* Loop through the objects you want to upload.
* Open a file handle, `f`, for each file and call `WriteMetadataToFsys()`
* Call `f.Seek(0, io.SeekStart)` on the file handle to ensure that all the bytes are captured during upload phase
    * If this fails, call `cleanup()` which is returned by `WriteMetadataToFsys()`
* Upload the file to the bucket
    * If this fails, call `cleanup()` which is returned by `WriteMetadataToFsys()`


## Retrieve
When retrieving files, the following items need to be implemented:

* Read the metadata file
* Make sure the local path exists where the file will reside based on the metadata
* Download the file into the given location
* Confirm the hash of the local file matches the hash in the metadata, Error if it does not
* Delete the file if the hashes do not match
