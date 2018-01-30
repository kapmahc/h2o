# H2O

A complete open source e-commerce solution.

## For Development

- install go by gvm

  ```bash
  GOPATH=$HOME/go
  PATH=$GOPATH/bin:$PATH
  export GOPATH PATH
  ```

- nvm global package path

```bash
NPM_PACKAGES="${HOME}/.npm-packages"
PATH="$NPM_PACKAGES/bin:$PATH"
export NPM_PACKAGES PATH
```

## Usage

```bash
go get -u github.com/kardianos/govendor
go get -d -u github.com/kapmahc/h2o
cd $GOPATH/src/github.com/kapmahc/h2o
govendor sync
make
```

## Atom plugins

enable autosave

- go-plus
- file-icons
- atom-beautify(enable newline, beautify on save; need python-sqlparse)
- language-babel
- language-ini

## Notes

- Generate a random key

  ```
  openssl rand -base64 32
  ```

- ~/.npmrc

  ```
  prefix=${HOME}/.npm-packages
  ```

- Create database

```
psql -U postgres
CREATE DATABASE db-name WITH ENCODING = 'UTF8';
CREATE USER user-name WITH PASSWORD 'change-me';
GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;
```

- Chrome browser: F12 => Console settings => Log XMLHTTPRequests

- Rabbitmq Management Plugin(<http://localhost:15612>)

  ```
  rabbitmq-plugins enable rabbitmq_management
  rabbitmqctl change_password guest change-me
  rabbitmqctl add_user who-am-i change-me
  rabbitmqctl set_user_tags who-am-i administrator
  rabbitmqctl list_vhosts
  rabbitmqctl add_vhost v-host
  rabbitmqctl set_permissions -p v-host who-am-i ".*" ".*" ".*"
  ```

- "RPC failed; HTTP 301 curl 22 The requested URL returned error: 301"

  ```
  git config --global http.https://gopkg.in.followRedirects true
  ```

- 'Peer authentication failed for user', open file "/etc/postgresql/9.5/main/pg_hba.conf" change line:

  ```
  local   all             all                                     peer  
  TO:
  local   all             all                                     md5
  ```

- Generate openssl certs

  ```
  openssl genrsa -out www.change-me.com.key 2048
  openssl req -new -x509 -key www.change-me.com.key -out www.change-me.com.crt -days 3650 # Common Name:*.change-me.com
  ```

## Documents

- [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)

- [favicon.ico](http://icoconvert.com/)

- [smver](http://semver.org/)

- [banner.txt](http://patorjk.com/software/taag/)

- [Ant Design](https://ant.design/docs/react/introduce)

- [Ant Design Pro](https://pro.ant.design/docs/getting-started)

- [AWS](http://docs.aws.amazon.com/general/latest/gr/rande.html)

- [quill](https://quilljs.com/)
