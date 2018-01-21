use std::fs;
use std::os::unix::fs::OpenOptionsExt;
use std::io::Write;
use std::path::{Path, PathBuf};

use rand::{self, Rng};
use docopt::Docopt;
use toml;
use time;
use rocket;
use base64;

use super::result::{Error, Result};
use super::env;

#[derive(Debug, Deserialize)]
struct Args {
    flag_version: bool,
    flag_https: bool,
    flag_name: String,
    flag_daemon: bool,

    cmd_generate: bool,
    cmd_locale: bool,
    cmd_migration: bool,
    cmd_config: bool,
    cmd_nginx: bool,

    cmd_database: bool,
    cmd_create: bool,
    cmd_connect: bool,
    cmd_migrate: bool,
    cmd_rollback: bool,
    cmd_drop: bool,
    cmd_status: bool,

    cmd_start: bool,
    cmd_stop: bool,
}

pub fn run() -> Result<()> {
    let usage = format!(
        "
{name} - {description}.

VERSION: {version}
AUTHORS: {authors}
HOMEPAGE: {homepage}

USAGE:
  {name} generate config
  {name} generate (locale|migration) [--name=<fn>]
  {name} generate nginx [--https]
  {name} database (create|connect|migrate|rollback|status|drop)
  {name} start [--daemon]
  {name} stop
  {name} (-h | --help)
  {name} --version

OPTIONS:
  -h --help     Show this screen.
  --version     Show version.
  --name=<fn>   File's name.
  --https       Using https?
  --daemon      Run as daemon mode?
    ",
        version = env::VERSION,
        name = env::NAME,
        description = env::DESCRIPTION,
        homepage = env::HOMEPAGE,
        authors = env::AUTHORS,
    );
    let args: Args = try!(try!(Docopt::new(usage)).deserialize());
    // println!("{:?}", args);
    let app = App {};

    if args.flag_version {
        return app.show_version();
    }
    if args.cmd_start {
        return app.start(args.flag_daemon);
    }
    if args.cmd_stop {
        return app.stop();
    }
    if args.cmd_generate {
        if args.cmd_config {
            return app.generate_config();
        }
        if args.cmd_nginx {
            return app.generate_nginx(args.flag_https);
        }
        if args.cmd_migration {
            return app.generate_migration(args.flag_name);
        }
        if args.cmd_locale {
            return app.generate_locale(args.flag_name);
        }
    }
    if args.cmd_database {
        if args.cmd_create {
            return app.database_create();
        }
        if args.cmd_connect {
            return app.database_connect();
        }
        if args.cmd_migrate {
            return app.database_migrate();
        }
        if args.cmd_rollback {
            return app.database_rollback();
        }
        if args.cmd_status {
            return app.database_status();
        }
        if args.cmd_drop {
            return app.database_drop();
        }
    }

    return Ok(());
}

struct App {}

impl App {
    fn start(&self, daemon: bool) -> Result<()> {
        return Ok(());
    }

    fn stop(&self) -> Result<()> {
        return Ok(());
    }

    fn generate_nginx(&self, https: bool) -> Result<()> {
        return Ok(());
    }

    fn generate_locale(&self, name: String) -> Result<()> {
        if name.is_empty() {
            return Err(Error::NotFound);
        }
        let root = Path::new(self.locales_dir());
        try!(fs::create_dir_all(&root));

        let mut file = root.join(name);
        file.set_extension(self.locales_ext());
        println!("generate file {}", file.display());
        try!(
            fs::OpenOptions::new()
                .write(true)
                .create_new(true)
                .mode(0o600)
                .open(file)
        );

        return Ok(());
    }

    fn generate_migration(&self, name: String) -> Result<()> {
        if name.is_empty() {
            return Err(Error::NotFound);
        }

        let root = self.migrations_dir("postgres".to_string())
            .join(try!(time::strftime("%Y%m%d%H%M%S", &time::now_utc())))
            .join(name);
        try!(fs::create_dir_all(&root));
        let files = vec!["up", "down"];
        for n in files.into_iter() {
            let mut file = root.join(n);
            file.set_extension(self.migrations_ext());
            println!("generate file {}", file.display());
            try!(
                fs::OpenOptions::new()
                    .write(true)
                    .create_new(true)
                    .mode(0o600)
                    .open(file)
            );
        }

        return Ok(());
    }

    fn generate_config(&self) -> Result<()> {
        let mut secret: Vec<u8> = (0..32).collect();
        rand::thread_rng().shuffle(&mut secret);

        let cfg = env::Config {
            secret: base64::encode(&secret),
            env: rocket::config::Environment::Development.to_string(),
            http: env::Http {
                name: "www.change-me.com".to_string(),
                limits: 1 << 15,
                port: 8080,
                workers: 4,
                theme: "moon".to_string(),
            },
            database: env::Database {
                driver: "postgres".to_string(),
                host: "localhost".to_string(),
                port: 5432,
                user: "postgres".to_string(),
                name: env::NAME.to_string(),
                password: "".to_string(),
                extra: [("sslmode".to_string(), "disable".to_string())]
                    .iter()
                    .cloned()
                    .collect(),
            },
            redis: env::Redis {
                host: "localhost".to_string(),
                port: 6379,
                db: 0,
            },
            rabbitmq: env::RabbitMQ {
                host: "localhost".to_string(),
                port: 5672,
                user: "guest".to_string(),
                password: "guest".to_string(),
                _virtual: env::NAME.to_string(),
            },
        };
        let buf = try!(toml::to_vec(&cfg));

        let name = self.config_file();
        println!("generate file {}", name);
        let mut file = try!(
            fs::OpenOptions::new()
                .write(true)
                .create_new(true)
                .mode(0o600)
                .open(name)
        );
        try!(file.write_all(&buf));
        return Ok(());
    }

    fn database_create(&self) -> Result<()> {
        return Ok(());
    }

    fn database_connect(&self) -> Result<()> {
        return Ok(());
    }

    fn database_migrate(&self) -> Result<()> {
        return Ok(());
    }

    fn database_rollback(&self) -> Result<()> {
        return Ok(());
    }

    fn database_status(&self) -> Result<()> {
        return Ok(());
    }

    fn database_drop(&self) -> Result<()> {
        return Ok(());
    }

    fn show_version(&self) -> Result<()> {
        println!("{}", env::VERSION);
        return Ok(());
    }

    fn config_file(&self) -> &'static str {
        return "config.toml";
    }
    fn locales_dir(&self) -> &'static str {
        return "locales";
    }
    fn locales_ext(&self) -> &'static str {
        return "ini";
    }
    fn migrations_dir(&self, driver: String) -> PathBuf {
        return Path::new("db").join(driver).join("migrations");
    }
    fn migrations_ext(&self) -> &'static str {
        return "sql";
    }
}
