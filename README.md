# Bintray Resource

Versions objects in bintray

## Source Configuration

* `subject_name`: *Required.* The Bintray Subject under which to store the package

* `repo_name`: *Required.* The Bintray Repository under which to store the package

* `package_name`: *Required.* The name of the package.

* `username`: *Required.* The username for authentication.

* `api_key`: *Required.* The apikey for authentication. 

## Behavior

### `check`: Fetch versions for the package

Versions will be found for the mentioned package

### `in`: Download the package for the version.

Places the following files in the destination:

* `(package_name)`: The file fetched from bintray.

#### Parameters

*None.*


### `out`: Upload a package to bintray.

Given a path specified by `from`, upload it to bintray. 
The path must identify a single file.

#### Parameters

* `from`: *Required.* A path specifying the file to upload.

* `version_file`: *Required.* A path to a file specifying the version of the file being uploaded.
