[![test](https://github.com/kkentzo/rv/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kkentzo/rv/actions/workflows/test.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/kkentzo/rv)

# rv

`rv` is a command-line tool that manages multiple local release
versions of a software bundle under a user-specified workspace
directory. `rv` does not manage the deployment process, i.e. the
bundle (in the form of a compressed archive such as zip) must already
exist in the localhost.

`rv` also manages a `current` link under the user-specified workspace
directory. The `current` link points to the "active" release
version. `rv` can also be instructed to keep the N most recent
releases (default: 3).

`rv` is agnostic as to the type of the software bundle; it will just
decompress its contents to the appropriate release directory under the
workspace.

## Usage

All commands and options can be inspected by using `rv help`.

### Release new version

`rv` can be instructed to "release" a new version of the software
bundle in the form of an archive file (e.g. zip) under a
user-specified workpace directory. `rv` will create a new release
directory under the workspace to which the archive contents will be
extracted. It will also update the `current` link to point to the new
release version. Clients can use the `$WORKSPACE/current` path in
order to obtain access to the "active" release.

For example, given a software bundle located at `/tmp/bundle.zip` and
the workspace directory `/opt/workspace`, the release will unzip the
contents of the zip file to a release directory under `/opt/workspace`
and will update the current link:

```bash
$ rv release -w /opt/workspace -a /tmp/bundle.zip
[info] workspace=/opt/workspace
[info] release=20240313151207.365
[release] unpacking bundle=/tmp/bundle.zip to /opt/workspace/20240313151207.365
[release] update current to 20240313151207.365
[success] active version is 20240313151207.365

$ ls -l /opt/workspace
drwxr-xr-x 3 user group 4096 Mar 13 15:12 20240313151207.365
lrwxrwxrwx 1 user group   18 Mar 13 15:12 current -> 20240313151207.365
```

If we release a new version, the old version will be kept (subject to
the `--keep` flag) and the `current` link will be updated to point to
the new release:

```bash
$ rv release -w /opt/workspace -a /tmp/bundle.zip
[info] workspace=/opt/workspace
[info] release=20240313151323.508
[release] unpacking bundle=/tmp/bundle.zip to /opt/workspace/20240313151323.508
[release] update current to 20240313151323.508
[success] active version is 20240313151323.508

$ ls -l /opt/workspace
drwxr-xr-x 3 user group 4096 Mar 13 15:12 20240313151207.365
drwxr-xr-x 3 user group 4096 Mar 13 15:13 20240313151323.508
lrwxrwxrwx 1 user group   18 Mar 13 15:13 current -> 20240313151323.508
```

## List all available release versions

`rv` can display all installed versions under a workspace with the
following command:

```bash
$ rv list -w /opt/workspace
20240313151323.508 <== current
20240313151207.365
```

## Rewind to an existing version

`rv` can be instructed to perform a rewind operation from the latest
release to a previous target release. If no specific target is
specified by the user, then the `current` link is set to the release
that precedes the current one.

The rewind operation will **delete** all releases that were performed
after the target release so as to maintain the integrity of the
cleanup operations.

So, given:

```bash
$ rv list -w /opt/workspace
20240313151323.508 <== current
20240313151207.365
```

the rollback operation to the previous version can be achieved by:

```bash
$ rv rollback -w /opt/workspace
[info] current=20240313151323.508
[rewind] setting current to 20240313151207.365
[cleanup] deleting 20240313151323.508
[success] active version is 20240313151207.365
```
