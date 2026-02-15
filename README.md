# CliDiff

CLI interface breaking change detector. Catches removed flags and subcommands **before** your users' scripts explode.

Single static binary. Zero dependencies. CI-ready.

## Install

```bash
go install github.com/openkickstart/clidiff@latest
# or
git clone https://github.com/openkickstart/clidiff && cd clidiff && go build -o clidiff .
```

## Usage

### 1. Take a snapshot of your CLI

```bash
clidiff snapshot mycli -o baseline.json
```

This runs `mycli --help`, extracts all `--flags` and subcommands, and saves a JSON snapshot.

### 2. After changes, take another snapshot

```bash
clidiff snapshot mycli -o current.json
```

### 3. Diff them

```bash
clidiff diff baseline.json current.json
```

Output:
```
❌ BREAKING: flag removed: --output
❌ BREAKING: subcommand removed: deploy
✅ Added flag: --json

⚠️  Breaking changes detected! Exit code 1.
```

**Exit code 1** on breaking changes — plug it straight into CI.

## CI Integration (GitHub Actions)

```yaml
- name: Check CLI compatibility
  run: |
    clidiff snapshot ./myapp -o new.json
    clidiff diff baseline.json new.json
```

## How It Works

1. Parses `--help` output with regex to extract `--flag` patterns
2. Detects subcommand sections ("Commands:", "Available Commands:", etc.)
3. Compares old vs new: removed = BREAKING, added = safe

## License

MIT
