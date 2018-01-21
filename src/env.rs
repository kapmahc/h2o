use std::collections::HashMap;

pub const VERSION: &'static str = env!("CARGO_PKG_VERSION");
pub const NAME: &'static str = env!("CARGO_PKG_NAME");
pub const DESCRIPTION: &'static str = env!("CARGO_PKG_DESCRIPTION");
pub const HOMEPAGE: &'static str = env!("CARGO_PKG_HOMEPAGE");
pub const AUTHORS: &'static str = env!("CARGO_PKG_AUTHORS");

#[derive(Serialize, Deserialize, Debug)]
pub struct Config {
    pub secret: String,
    pub env: String,
    pub http: Http,
    pub database: Database,
    pub redis: Redis,
    pub rabbitmq: RabbitMQ,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Http {
    pub name: String,
    pub theme: String,
    pub workers: u8,
    pub port: u16,
    pub limits: u64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Database {
    pub driver: String,
    pub host: String,
    pub port: u16,
    pub user: String,
    pub name: String,
    pub password: String,
    pub extra: HashMap<String, String>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Redis {
    pub host: String,
    pub port: u16,
    pub db: u8,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct RabbitMQ {
    pub host: String,
    pub port: u16,
    pub user: String,
    #[serde(rename = "vrirtual")]
    pub _virtual: String,
    pub password: String,
}
