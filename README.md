# H2O

A complete open source e-commerce solution by Rust.

## Install

- install rust

  ```
  curl https://sh.rustup.rs -sSf | sh
  rustup default nightly
  cargo install rustfmt-nightly
  cargo install racer
  rustup component add rust-src
  ```

- add to your .zshrc

  ```
  export PATH="$HOME/.cargo/bin:$PATH"
  export RUST_SRC_PATH="$(rustc --print sysroot)/lib/rustlib/src/rust/src"
  ```

- test racer

  ```
  racer complete std::io::B
  ```

- test run

  ```
  cargo run -- --version
  ```

## Atom plugins

enable autosave

- language-rust
- racer
- file-icons
- atom-beautify(enable newline, beautify on save; need python-sqlparse)
- language-babel
- language-ini

## Notes

- Generate a random key ````

openssl rand -base64 32

```

- ~/.npmrc
```

prefix=${HOME}/.npm-packages

```

- Create database
```

psql -U postgres CREATE DATABASE db-name WITH ENCODING = 'UTF8'; CREATE USER user-name WITH PASSWORD 'change-me'; GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;

```

- ueditor
```

cd node_modules/ueditor npm install grunt-cli -g npm install grunt

```

- Chrome browser: F12 => Console settings => Log XMLHTTPRequests

- Rabbitmq Management Plugin(<http://localhost:15612>)
```

rabbitmq-plugins enable rabbitmq_management rabbitmqctl change_password guest change-me rabbitmqctl add_user who-am-i change-me rabbitmqctl set_user_tags who-am-i administrator rabbitmqctl list_vhosts rabbitmqctl add_vhost v-host rabbitmqctl set_permissions -p v-host who-am-i "._" "._" ".*"

```

- 'Peer authentication failed for user', open file "/etc/postgresql/9.5/main/pg_hba.conf" change line:
```

local all all peer<br>
TO: local all all md5

```

- Generate openssl certs
```

openssl genrsa -out www.change-me.com.key 2048 openssl req -new -x509 -key www.change-me.com.key -out www.change-me.com.crt -days 3650 # Common Name:*.change-me.com ```

## Documents

- [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)

- [favicon.ico](http://icoconvert.com/)

- [smver](http://semver.org/)

- [banner.txt](http://patorjk.com/software/taag/)

- [The Rust Programming Language](https://doc.rust-lang.org/book/second-edition/)

- [Rocket](https://github.com/SergioBenitez/Rocket)

- [Bootstrap](http://getbootstrap.com/docs/4.0/getting-started/introduction/)

- [AdminLTE](https://github.com/almasaeed2010/AdminLTE)
