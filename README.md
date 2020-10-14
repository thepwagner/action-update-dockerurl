# action-update-dockerurl

## This is not endorsed by or associated with GitHub, Dependabot, etc.

This action checks for available updates to binaries fetched in a Dockerfile, from GitHub releases.

* Updates checksums if included in release (e.g. `SHASUMS.txt`)
* Calculates new checksim if not included in release
* All the features common to [action-update](https://github.com/thepwagner/action-update) actions

Expects a pattern like:
```
# These will be updated
ARG BUF_VERSION=v0.26.0
ARG BUF_CHECKSUM=6bab7d8de7c39558a955c6dc2d19331fb5d630eff9b0f5a4de5e77548db20331

# How they are used is up to you:
RUN curl -o /usr/local/bin/buf -L https://github.com/bufbuild/buf/releases/download/${BUF_VERSION}/buf-Linux-x86_64
RUN echo "$BUF_CHECKSUM  /usr/local/bin/buf" | sha256sum -c
```

## Simplest setup

```
- uses: actions/checkout@v2
  # If you use Actions "push" for CI too, a Personal Access Token is required for update PRs to trigger
  with:
    token: ${{ secrets.MY_GITHUB_PAT }}
- uses: actions/setup-go@v2
  with:
    go-version: '1.15.0'
- uses: thepwagner/action-update-dockerurl@main
  with:
    # If you use Actions "pull_request" for CI too, a Personal Access Token is required for update PRs to trigger
    token: ${{ secrets.MY_GITHUB_PAT }}
    groups: |
      - name: buf
        pattern: github.com/bufbuild/buf
        post-script: script/protoc
```
