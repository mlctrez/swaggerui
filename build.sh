#!/usr/bin/env bash

rm -rf web
rm -rf dist
mkdir dist
cd dist
git clone https://github.com/swagger-api/swagger-ui.git
cd swagger-ui
npm install
npm run-script build
cd ../..
mv dist/swagger-ui/dist web

cd web
mv index.html _index.html

zip ../dist/main.go.zip -q -r *

cd ..

go build -o dist/swaggerui main.go

cat dist/main.go.zip >> dist/swaggerui

zip -q -A dist/swaggerui

mv dist/swaggerui ~/bin/swaggerui

