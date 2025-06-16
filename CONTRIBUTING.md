# Contributing

## Tools

### Conventional Commit

- install git cz tool global

```sh
npm install -g commitizen
npm install -g cz-conventional-changelog
npm install -g conventional-changelog-cli
echo '{ "path": "cz-conventional-changelog" }' > ~/.czrc
```

### Pre-commit

- install pre-commit in any way you like

```sh
pre-commit autoupdate
pre-commit install
```

```sh
pre-commit run --all-files
```

## Modify CHANGELOG

- git-chglog

```sh
brew tap git-chglog/git-chglog
brew install git-chglog
```

- new tag in default branch

```sh
git checkout main
git pull
```

```sh
VERSION=1.2.1
git tag -a v$VERSION -m "Release v$VERSION"
git push -u origin --tags
git push -u origin --all
```

## Find ignored files

```sh
find . -type f  | git check-ignore --stdin
```
