# `purgo`

Clean up your filesystem! It's like a spell, you cast `purgo` and you are **magically able to delete large files with ease**.

Install purgo prebuilt binary in [Releases](https://github.com/kauefraga/purgo/releases).

Or clone this repository and build it yourself.

```sh
git clone https://github.com/kauefraga/purgo.git
cd purgo

go build -o purgo cmd/main.go

# This will generate a more lightweight binary
# CGO_ENABLED=0 go build -o purgo -ldflags='-w -s' cmd/main.go
```

## Usage

From your preferred terminal, just run

```sh
./purgo
```

Or specify the flag `--min-size int` to show only files/directories with more than this minimum size (in bytes, default is 1MB)

```sh
./purgo # shows all files/dirs with size > 1MB

./purgo --min-size 1 # shows all files/dirs with size > 1B
./purgo --min-size 1024 # shows all files/dirs with size > 1KB
```

Select a file using <kbd>SPACE</kbd> or <kbd>RETURN</kbd> to delete it, a modal will pop asking you to confirm the deletion.

Select a folder using the same keys to expand/collapse it.

**Notice**: after deleting a file you might need to move the cursor for the tree update.

## License

This project is licensed under the MIT License - See the [LICENSE](https://github.com/kauefraga/pavus/blob/main/LICENSE) for more information.
