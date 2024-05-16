# metrics-stack

- [metrics-stack](#metrics-stack)
  - [Environment Setup](#environment-setup)
    - [Dependencies](#dependencies)
    - [Environment Variables](#environment-variables)
  - [Tasks](#tasks)
    - [`build`/`build.snapshot`](#buildbuildsnapshot)
    - [`release.snapshot`](#releasesnapshot)
    - [`release`](#release)


This repo contains the lambda-metrics-stack that will be installed on Guest VMs. Inside are configuration files, tools, scripts, and other CI/CD bits to generate `.deb` packages that will be uploaded to Artifactory.

## Environment Setup

### Dependencies

This project requires Go. 

```
apt install go
```

or

```
brew install go
```

Then, download all dependencies with:

```
go mod download -x
```

### Environment Variables

The secrets in this repo are stored in `.env`. You must create this file and populate it with the environment variables that Task will tell you that you need. For example:

```
$ task release         
task: Task "release" cancelled because it is missing required variables: ARTIFACTORY_STAGING_SECRET, GITHUB_TOKEN
```

Different tasks need different environment variables, please refer to the output above to determine what's needed.

## Tasks

[Task](https://taskfile.dev/) is used instead of Makefile. Why? Task is just better. Trust me, you will see :D 

```
$ task --list
task: Available tasks for this project:
* build:                  Build locally, but fail if git state is dirty.
* build.snapshot:         Build locally and ignore a dirty git state.
* release:                Build and submit a new release.
* release.snapshot:       Build and submit a new release, but ignore dirty git state.
```

### `build`/`build.snapshot`

The `build` command only builds the executables and metadata. It does not create `.deb` files.

You can build a local `.deb` file using `task build.snapshot`. The `.snapshot` variant is used if your `HEAD` is not on the latest tagged commit (for example, if you've made changes to Git but have not tagged it, you need `.snapshot`). Otherwise, `task build` should be used.

```
$ task build.snapshot
task: [build.snapshot] goreleaser build --clean --snapshot
  • starting build...
  • loading                                          path=.goreleaser.yaml
  • skipping validate...
  • loading environment variables
    • using token from  $GITHUB_TOKEN 
  • getting and validating git state
    • git state                                      commit=4c89c82c6e5530a5cbb85913367d817ecb23a249 branch=main current_tag=v0.0.17 previous_tag=v0.0.16 dirty=true
    • pipe skipped                                   reason=disabled during snapshot mode
  • parsing tag
  • setting defaults
  • snapshotting
    • building snapshot...                           version=0.0.17-SNAPSHOT-4c89c82
  • checking distribution directory
    • cleaning dist
  • setting up metadata
  • storing release metadata
    • writing                                        file=dist/metadata.json
  • loading go mod information
  • build prerequisites
  • writing effective config file
    • writing                                        config=dist/config.yaml
  • building binaries
    • building                                       binary=dist/default_linux_amd64_v1/lambda-support
    • building                                       binary=dist/default_linux_arm64/lambda-support
    • took: 3s
  • storing artifacts metadata
    • writing                                        file=dist/artifacts.json
  • build succeeded after 4s
  • thanks for using goreleaser!
```

You'll see the built artifacts in `dist/`

```
$ ls -lah dist
total 32
drwxr-xr-x   7 landon  staff   224B May 16 11:56 .
drwxr-xr-x  17 landon  staff   544B May 16 11:56 ..
-rw-r--r--   1 landon  staff   511B May 16 11:56 artifacts.json
-rw-r--r--   1 landon  staff   4.8K May 16 11:56 config.yaml
drwxr-xr-x   3 landon  staff    96B May 16 11:56 default_linux_amd64_v1
drwxr-xr-x   3 landon  staff    96B May 16 11:56 default_linux_arm64
-rw-r--r--   1 landon  staff   255B May 16 11:56 metadata.json
```


### `release.snapshot`

The `release` command will build the artifacts and create `.deb` files. The `release.snapshot` version will build all artifacts, but will not publish the artifacts to Artifactory.

### `release`

This command is the full, end-to-end build-publish-release flow for uploading `.deb` files. This requires the Git state to be clean, meaning that `HEAD` is on a tagged version.
