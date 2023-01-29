# Yandex Disk File Upload

This is a simple script to authenticate yourself and upload files to Yandex Disk. It is written in Go and uses the [Yandex Disk REST API][1] to upload files to your Yandex Disk account.

## Getting Started

1. Create a new application in the [Yandex Disk OAuth][2] page. You will need to provide a redirect URL. Since this script is meant to be run locally for now, we use `http://localhost:8080/callback` as the redirect URL.

2. Copy the .env.example file to .env and fill in the values for the variables.

3. Get the dependencies for the script by running the following command:

```bash
go get
```

4. Run the script with the following command:

```bash
go run main.go source/file.txt app:/remote/file.txt
```

## Build Windows executable

```bash
GOOS=windows GOARCH=amd64 go build
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

[1]: https://yandex.com/dev/disk/api/reference/upload.html
[2]: https://oauth.yandex.com/client/new
