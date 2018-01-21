use std::fs;
use std::os::unix::fs::OpenOptionsExt;
use std::io::{Read, Write};
use std::path::{Path, PathBuf};
use std::process::Command;

use rand::{self, Rng};
use toml;
use chrono::prelude::{DateTime, Utc};
use rocket;
use base64;

use super::result::{Error, Result};
use super::env;
use super::db::{Database, PostgreSQL};

pub struct App {}

impl App {
    pub fn start(&self, daemon: bool) -> Result<()> {
        return Ok(());
    }

    pub fn stop(&self) -> Result<()> {
        return Ok(());
    }

    pub fn generate_nginx(&self, https: bool) -> Result<()> {
        return Ok(());
    }

    pub fn generate_locale(&self, name: String) -> Result<()> {
        if name.is_empty() {
            return Err(Error::NotFound);
        }
        let root = Path::new(self.locales_dir());
        try!(fs::create_dir_all(&root));

        let mut file = root.join(name);
        file.set_extension(self.locales_ext());
        info!("generate file {}", file.display());
        try!(
            fs::OpenOptions::new()
                .write(true)
                .create_new(true)
                .mode(0o600)
                .open(file)
        );

        return Ok(());
    }

    pub fn generate_migration(&self, name: String) -> Result<()> {
        if name.is_empty() {
            return Err(Error::NotFound);
        }

        let cfg = try!(self.parse_config());
        let now: DateTime<Utc> = Utc::now();
        let root = self.migrations_dir(cfg.database.driver).join(format!(
            "{}_{}",
            now.format("%Y%m%d%H%M%S"),
            name
        ));
        try!(fs::create_dir_all(&root));
        let files = vec!["up", "down"];
        for n in files.into_iter() {
            let mut file = root.join(n);
            file.set_extension(self.migrations_ext());
            info!("generate file {}", file.display());
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

    pub fn generate_config(&self) -> Result<()> {
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
                driver: env::POSTGRESQL.to_string(),
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
            // rabbitmq: env::RabbitMQ {
            //     host: "localhost".to_string(),
            //     port: 5672,
            //     user: "guest".to_string(),
            //     password: "guest".to_string(),
            //     _virtual: env::NAME.to_string(),
            // },
        };
        let buf = try!(toml::to_vec(&cfg));

        let name = self.config_file();
        info!("generate file {}", name);
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

    pub fn database_connect(&self) -> Result<()> {
        let cfg = try!(self.parse_config()).database;
        match cfg.driver.as_ref() {
            env::POSTGRESQL => {
                try!(
                    Command::new("psql")
                        .arg("-U")
                        .arg(cfg.user)
                        .arg("-h")
                        .arg(cfg.host)
                        .arg("-p")
                        .arg(cfg.port.to_string())
                        .arg("-d")
                        .arg(cfg.name)
                        .status()
                );
                Ok(())
            }
            _ => Err(Error::NotFound),
        }
    }

    pub fn database_migrate(&self) -> Result<()> {
        self.database(
            |d| {
                let mut files = Vec::new();
                for it in try!(fs::read_dir(self.migrations_dir(d.driver().to_string()))) {
                    files.push(try!(it));
                }
                files.sort_by_key(|it| it.path());

                let versions = try!(d.versions());
                let items: Vec<String> = versions
                    .iter()
                    .map(|it| {
                        let &(ref vr, ref _ts) = it;
                        vr.to_string()
                    })
                    .collect();

                for it in files {
                    if let Some(path) = it.path().file_name() {
                        if let Some(name) = path.to_str() {
                            info!("find migration {}", name);
                            if !items.contains(&name.to_string()) {
                                let mut file = it.path().join("up");
                                file.set_extension(self.migrations_ext());
                                let mut fd = try!(fs::OpenOptions::new().read(true).open(file));
                                let mut buf = String::new();
                                try!(fd.read_to_string(&mut buf));
                                try!(d.up(name.to_string(), buf));
                            }
                        }
                    }
                }
                Ok(())
            },
            true,
        )
    }

    pub fn database_rollback(&self) -> Result<()> {
        self.database(
            |d| match try!(d.versions()).pop() {
                Some(it) => {
                    let (vr, _ts) = it;
                    let name = vr.to_string();
                    let mut file = self.migrations_dir(d.driver().to_string())
                        .join(&name)
                        .join("down");
                    file.set_extension(self.migrations_ext());
                    let mut fd = try!(fs::OpenOptions::new().read(true).open(file));
                    let mut buf = String::new();
                    try!(fd.read_to_string(&mut buf));
                    try!(d.down(name, buf));
                    Ok(())
                }
                None => Err(Error::NotFound),
            },
            true,
        )
    }

    pub fn database_status(&self) -> Result<()> {
        self.database(
            |c| {
                println!("{:<32}\t{}", "VERSION", "CREATED AT");
                for (v, t) in try!(c.versions()) {
                    println!(
                        "{:<32}\t{}",
                        v,
                        DateTime::<Utc>::from_utc(t, Utc).to_rfc2822()
                    );
                }
                Ok(())
            },
            true,
        )
    }

    pub fn database<F>(&self, f: F, open: bool) -> Result<()>
    where
        F: Fn(&Database) -> Result<()>,
    {
        let cfg = try!(self.parse_config()).database;
        match cfg.driver.as_ref() {
            env::POSTGRESQL => {
                let mut con = PostgreSQL::new(cfg);
                if open {
                    try!(con.open());
                }
                return f(&con);
            }
            _ => Err(Error::NotFound),
        }
    }

    pub fn show_version(&self) -> Result<()> {
        info!("{}", env::VERSION);
        return Ok(());
    }

    fn parse_config(&self) -> Result<env::Config> {
        let mut file = try!(fs::File::open(self.config_file()));
        let mut buf = Vec::new();
        try!(file.read_to_end(&mut buf));
        let cfg: env::Config = try!(toml::from_slice(&buf));
        return Ok(cfg);
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
