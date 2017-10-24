# go-irleak

IRLeak is a project from the Mobile and Pervasive Sensor Systems Lab at University of Maryland, Baltimore County, to create a low-cost device for continuous home energy auditing. `go-irleak` is a server for collecting data associated with the project and storing it in a database. It's meant to be run behind a reverse proxy, e.g. Nginx, which should implement TLS. `go-irleak` has an authentication mechanism and stores data for multiple users from several sources.

## Dependencies

1. `go get github.com/spf13/cobra` Cobra Commander is helpful for running multiple commands from the same Go binary.
2. `go get github.com/spf13/viper` Viper is an excellent tool for parsing config files.
3. `go get github.com/go-sql-driver/mysql` The lab is standardized on MariaDB, so that is the preferred database back-end.
4. `go get github.com/mattn/go-sqlite3` This application also supports SQLite as an alternative back-end, mainly for testing.
5. `go get github.com/bgentry/speakeasy` Speakeasy is a handy library for password prompts if you don't want to put a password in a config file or leave it on your terminal screen.
6. `go get github.com/elithrar/simple-scrypt` Don't store passwords. Store salted hashes.

## Usage

`go-irleak server` starts up the server with whatever configuration is in `config.yaml`.

`go-irleak useradd` is a script for adding a user to the database so he can start uploading data.
