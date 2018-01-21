use std::ops::Deref;

use r2d2::Pool;
use r2d2_redis::RedisConnectionManager;
use _redis::{cmd, ConnectionAddr, ConnectionInfo};

use super::result::Result;

pub struct Redis {
    pool: Pool<RedisConnectionManager>,
}

impl Redis {
    pub fn open(host: String, port: u16, db: u8) -> Result<Redis> {
        info!("open redis://{}:{}/{}", host, port, db);
        let manager = try!(RedisConnectionManager::new(ConnectionInfo {
            addr: Box::new(ConnectionAddr::Tcp(host, port)),
            db: db as i64,
            passwd: Some("".to_string()),
        }));
        let pool = try!(Pool::builder().build(manager));
        return Ok(Redis { pool: pool });
    }

    pub fn ping(&self) -> Result<()> {
        let con = try!(self.pool.get());
        let _ = try!(cmd("PING").query::<String>(con.deref()));
        return Ok(());
    }

    pub fn set(&self) -> Result<()> {
        return Ok(());
    }
}
