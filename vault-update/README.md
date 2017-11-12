# Vault-update

If you lack vault(Hashicorp) PATCH feature, this wrapper is for you.

To update value just run: `vault-update secret/<your secret path> key=value`.

You even can load value from file: `vault-update secret/<your secret path> key=@file_name`

## Requirements

- vault(cli)
- jq
