use std::ops::Deref;

use r2d2::Pool;
use r2d2_postgres::{PostgresConnectionManager, TlsMode};

use super::result::Result;

pub struct PostgreSQL {
    pool: Pool<PostgresConnectionManager>,
}

impl PostgreSQL {
    pub fn open(host: String) -> Result<PostgreSQL> {
        let manager = try!(PostgresConnectionManager::new(
            "postgres://postgres@localhost",
            TlsMode::None
        ));
        let pool = try!(Pool::new(manager));
        return Ok(PostgreSQL { pool: pool });
    }

    pub fn ping(&self) -> Result<()> {
        let con = try!(self.pool.get());
        let _ = try!(con.execute("SELECT NOW()", &[]));
        return Ok(());
    }
}
