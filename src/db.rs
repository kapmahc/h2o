use r2d2::{Pool, PooledConnection};
use r2d2_postgres::{PostgresConnectionManager, TlsMode};
use postgres;
use chrono::NaiveDateTime;

use super::env;
use super::result::{Error, Result};

pub trait Database {
    fn driver(&self) -> &str;
    fn ping(&self) -> Result<()>;
    fn create(&self) -> Result<()>;
    fn drop(&self) -> Result<()>;
    fn open(&mut self) -> Result<()>;
    fn up(&self, name: String, script: String) -> Result<()>;
    fn down(&self, name: String, script: String) -> Result<()>;
    fn versions(&self) -> Result<Vec<(String, NaiveDateTime)>>;
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

    fn tls_mode(&self) -> TlsMode {
        return match self.config.extra.get("sslmode") {
            _ => TlsMode::None,
        };
    }

    pub fn client(&self) -> Result<PooledConnection<PostgresConnectionManager>> {
        match self.pool {
            Some(ref p) => {
                let c = try!(p.get());
                Ok(c)
            }
            None => Err(Error::NotFound),
        }
    }

    // fn query<F>(&self, f: F) -> Result<()>
    // where
    //     F: Fn(PooledConnection<PostgresConnectionManager>) -> Result<()>,
    // {
    //     match self.pool {
    //         Some(ref pool) => {
    //             let con = try!(pool.get());
    //             return f(con);
    //         }
    //         None => Err(Error::NotFound),
    //     }
    // }

    // fn transaction<F>(&self, f: F) -> Result<()>
    // where
    //     F: Fn(&postgres::transaction::Transaction) -> Result<()>,
    // {
    //     self.query(|c| {
    //         let t = try!(c.transaction());
    //         match f(&t) {
    //             Ok(_) => {
    //                 try!(t.commit());
    //                 Ok(())
    //             }
    //             Err(e) => Err(e),
    //         }
    //     })
    // }

    fn check(&self, t: &postgres::transaction::Transaction) -> Result<()> {
        try!(t.execute(
            "CREATE TABLE IF NOT EXISTS schema_migrations(version VARCHAR(255) PRIMARY KEY, created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW())",
            &[]
        ));
        return Ok(());
    }
}

impl Database for PostgreSQL {
    fn ping(&self) -> Result<()> {
        let c = try!(self.client());
        try!(c.execute("SELECT NOW()", &[]));
        return Ok(());
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
        let c = try!(self.client());
        let t = try!(c.transaction());
        try!(t.batch_execute(&script));
        try!(t.execute(
            "INSERT INTO schema_migrations(version) VALUES($1)",
            &[&name]
        ));
        try!(t.commit());
        Ok(())
    }

    fn down(&self, name: String, script: String) -> Result<()> {
        info!("rollback {}", name);
        let c = try!(self.client());
        let t = try!(c.transaction());
        try!(t.batch_execute(&script));
        try!(t.execute("DELETE FROM schema_migrations WHERE version = $1", &[&name]));
        try!(t.commit());
        Ok(())
    }

    fn driver(&self) -> &str {
        &self.config.driver
    }

    fn versions(&self) -> Result<Vec<(String, NaiveDateTime)>> {
        let mut items = Vec::new();
        let c = try!(self.client());
        let t = try!(c.transaction());
        try!(self.check(&t));
        for row in &t.query(
            "SELECT version, created_at FROM schema_migrations ORDER BY version ASC",
            &[],
        )? {
            let version: String = row.get("version");
            let created_at: NaiveDateTime = row.get("created_at");
            items.push((version, created_at))
        }
        try!(t.commit());
        Ok(items)
    }
}
