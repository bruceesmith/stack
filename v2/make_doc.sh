#!/bin/bash
go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest ./... --output read
cat readme_header.md read readme_footer.md >README.md
rm read
