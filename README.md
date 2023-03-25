# Go Boilerplate

> This boilerplate is based on
> [techschool/simplebank](https://github.com/techschool/simplebank). The goal is
> to provide a well-structured way to start a Go-backed backend infrastructure.

## Schema Generation

Create the db schema in [dbdiagram.io](https://dbdiagram.io/home) to **decouple
the design from a specific database**. Then, click "Export" to export to the
language you want. In this boilerplate, we'll choose Postgres.

## Local Development

### Install Visual Studio Code extensions

```bash
code --install-extension esbenp.prettier-vscode
code --install-extension foxundermoon.shell-format
code --install-extension mtxr.sqltools
code --install-extension golang.go
```
