dist=dist
pkg=github.com/kapmahc/h2o/web
theme=moon

VERSION=`git rev-parse --short HEAD`
BUILD_TIME=`date -R`
AUTHOR_NAME=`git config --get user.name`
AUTHOR_EMAIL=`git config --get user.email`
COPYRIGHT=`head -n 1 LICENSE`
USAGE=`sed -n '3p' README.md`


build: api www
	cd $(dist) && tar cfJ ../$(dist).tar.xz *

www:
	cd ant-design-pro && git apply ../ant-design-pro.patch && npm run build
	-cp -r ant-design-pro/dist $(dist)/assets

api:
	go build -ldflags "-s -w -X ${pkg}.Version=${VERSION} -X '${pkg}.BuildTime=${BUILD_TIME}' -X '${pkg}.AuthorName=${AUTHOR_NAME}' -X ${pkg}.AuthorEmail=${AUTHOR_EMAIL} -X '${pkg}.Copyright=${COPYRIGHT}' -X '${pkg}.Usage=${USAGE}'" -o ${dist}/h2o main.go
	-cp -r db locales templates README.md $(dist)/

clean:
	-rm -r $(dist) $(dist).tar.xz
	cd ant-design-pro && git checkout -f

download:
	govendor sync
	git clone https://github.com/ant-design/ant-design-pro.git
	cd ant-design-pro && npm install
