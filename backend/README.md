# Backend

Backend is a backend of kakeibo.

## `mf` Command

Manage data of Money Forward ME.

```
Usage:
  mf [command]

Available Commands:
  csv         Download csv files from Money Forward ME
  db          Import files to database
  drive       Upload and download files to Google Drive
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.mf.yaml)
  -h, --help            help for mf
  -t, --toggle          Help message for toggle
```

### Setting the environment variables

Commands require the environment variables.

#### `csv` Command

- MONEY_FORWARD_EMAIL
  - The email address of your Money Forward ME account.
- MONEY_FORWARD_PASSWORD
  - The password of your Money Forward ME account.

#### `db` Command

- DB_DRIVER_NAME
  - `sqlite`
- DB_DSN
  - The path to `sqlite` file.
- TEST_DB_DRIVER_NAME
  - `sqlite`
- TEST_DB_DSN
  - The path to `sqlite` file which is used in unit test.

#### `drive` Command

- GOOGLE_APPLICATION_CREDENTIALS
  - An authentication credentials of your service account to use Google Cloud API.
  - See also [Google Cloud Authentication documentation](https://cloud.google.com/docs/authentication).
- GOOGLE_DRIVE_FOLDER_ID
  - The ID of the drive folder to upload and download files.
  - e.g. `https://drive.google.com/drive/folders/<foler id>`

## API

TBD
