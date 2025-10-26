# WC GO

This is a simple word counter CLI app written in Go that replicates (with a few twitches) the functionality of the `wc` command in linux.

### Display Options

- `--help`: Display help information.
- `-l`: Count the number of lines in the input.
- `-w`: Count the number of words in the input.
- `-c`: Count the number of characters (bytes) in the input.
- `-header`: Display a top level header for each column

> The flags can be used in any given order. If no flag is passed then all values are shown.

## Examples

### Single file

```bash
wc-go words.txt
```

### Multiple Files

```bash
wc-go words.txt example.txt
```

### With flags

```bash
wc-go -w words.txt
```

### No files (Stdin)

```bash
wc-go
```

```bash
echo 'foo bar baz' | wc-go
```

```bash
wc-go < words.txt
```
