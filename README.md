[![test](https://github.com/kkentzo/rv/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kkentzo/rv/actions/workflows/test.yml)

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
[info] bundle=/tmp/bundle.zip
[release] current=20240313151207.365
[info] finished 20240313151207.365

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
[info] bundle=/tmp/bundle.zip
[release] current=20240313151323.508
[info] finished 20240313151323.508

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

## Rollback to an existing version

TBD
