# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
# mv GMVivWiki.exe GMVivWiki_amd64.exe

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build