# envdiff

Compare `.env` files across environments and highlight missing or mismatched keys with colorized output.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base-env> <compare-env> [additional-envs...]
```

### Example

```bash
envdiff .env.example .env.production
```

**Sample output:**

```
✔  DB_HOST          present in both
✔  APP_PORT         present in both
✘  SECRET_KEY       missing in .env.production
⚠  LOG_LEVEL        mismatch  (debug vs info)
✘  REDIS_URL        missing in .env.example
```

### Flags

| Flag | Description |
|------|-------------|
| `--no-color` | Disable colorized output |
| `--keys-only` | Show only key names, omit values |
| `--missing` | Report only missing keys |

---

## How It Works

`envdiff` parses each provided `.env` file, builds a union of all keys, and compares values across files. Missing keys and value mismatches are highlighted using ANSI colors — red for missing, yellow for mismatches, and green for matches.

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)