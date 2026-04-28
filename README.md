# cli

## Process Managers

To allow a process manager to control `projdocs serve`'s lifecycle, use

```shell
projdocs serve --force --listen
```

- `--force` will ensure a clean start, even if there was not a clean exit during the previous run
- `--listen` holds-open the process (instead of pure daemonization) and will shutdown all services when `SIGTERM` or `SIGINT` are received