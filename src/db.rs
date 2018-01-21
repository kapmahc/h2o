use std::ops::Deref;
use std::collections::HashMap;

use r2d2::{Pool, PooledConnection};
use r2d2_postgres::{PostgresConnectionManager, TlsMode};

use super::env;
use super::result::{Error, Result};

pub trait Database {
    fn ping(&self) -> Result<()>;
    fn create(&self) -> Result<()>;
    fn drop(&self) -> Result<()>;
    fn open(&mut self) -> Result<()>;
    fn up(&self, name: String, script: String) -> Result<()>;
    fn down(&self, name: String, script: String) -> Result<()>;
    fn versions(&self) -> Result<Vec<String>>;
}

pub struct PostgreSQL {
    config: env::Database,
    pool: Option<Box<Pool<PostgresConnectionManager>>>,
}

impl PostgreSQL {
    pub fn new(cfg: env::Database) -> PostgreSQL {
        return PostgreSQL {
            config: cfg,
            pool: None,
        };
    }

    fn root(&self) -> Result<Pool<PostgresConnectionManager>> {
        let manager = try!(PostgresConnectionManager::new(
            format!(
                "postgres://{}:{}@{}:{}",
                self.config.user, self.config.password, self.config.host, self.config.port
            ),
            self.tls_mode(),
        ));
        let pool = try!(Pool::new(manager));
        return Ok(pool);
    }

    fn check(&self) -> Result<()> {
        return Ok(());
    }

    fn tls_mode(&self) -> TlsMode {
        return match self.config.extra.get("sslmode") {
            _ => TlsMode::None,
        };
    }

    fn query<DB>(&self, f: DB) -> Result<()>
    where
        DB: Fn(PooledConnection<PostgresConnectionManager>) -> Result<()>,
    {
        match self.pool {
            Some(ref pool) => {
                let con = try!(pool.get());
                return f(con);
            }
            None => Err(Error::NotFound),
        }
    }
    fn transaction(&self) -> Result<()> {
        return Ok(());
    }
}

impl Database for PostgreSQL {
    fn ping(&self) -> Result<()> {
        self.query(|c| {
            try!(c.execute("SELECT NOW()", &[]));
            return Ok(());
        })
    }
    fn open(&mut self) -> Result<()> {
        self.pool = Some(Box::new(try!(Pool::new(try!(
            PostgresConnectionManager::new(
                format!(
                    "postgres://{}:{}@{}:{}/{}",
                    self.config.user,
                    self.config.password,
                    self.config.host,
                    self.config.port,
                    self.config.name,
                ),
                self.tls_mode(),
            )
        )))));
        return Ok(());
    }

    fn create(&self) -> Result<()> {
        info!("create database {}", self.config.name);
        let con = try!(try!(self.root()).get());
        try!(con.execute(
            &format!(
                "CREATE DATABASE {} WITH ENCODING = 'UTF8'",
                self.config.name
            ),
            &[]
        ));
        return Ok(());
    }

    fn drop(&self) -> Result<()> {
        info!("drop database {}", self.config.name);
        let con = try!(try!(self.root()).get());
        try!(con.execute(&format!("DROP DATABASE {}", self.config.name), &[]));
        return Ok(());
    }

    fn up(&self, name: String, script: String) -> Result<()> {
        info!("migrate {}", name);
        // try!(con.execute(script, &[]));
        return Ok(());
    }

    fn down(&self, name: String, script: String) -> Result<()> {
        info!("rollback {}", name);
        return Ok(());
    }

    fn versions(&self) -> Result<Vec<String>> {
        let items = vec!["".to_string()];
        return Ok(items);
    }
}
