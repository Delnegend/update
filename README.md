# Update

Update checker for Windows applications that lack auto-update.

## How to use

1. Clone the repository, install [Go](https://go.dev/dl/), and execute `go build .`.
2. Create a `update.txt` file in the same directory as the executable. This file stores the executable paths of the applications, either for cases where they are not in the PATH or to override existing PATH entries. Refer to `update.example.txt` for details.
3. Run the executable.

## License

MIT