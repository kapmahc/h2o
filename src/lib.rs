#![feature(plugin)]
#![plugin(rocket_codegen)]

extern crate base64;
extern crate chrono;
extern crate docopt;
extern crate handlebars;
#[macro_use]
extern crate log;
extern crate postgres;
extern crate r2d2;
extern crate r2d2_postgres;
extern crate r2d2_redis;
extern crate rand;
extern crate redis as _redis;
extern crate rocket;
#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate serde_json;
extern crate toml;

pub mod nut;
pub mod forum;
pub mod survey;
pub mod reading;
pub mod erp;
pub mod mall;
pub mod pos;
pub mod ops;

pub mod console;
pub mod env;
pub mod result;
pub mod i18n;
pub mod cache;
pub mod db;
pub mod redis;
pub mod amqp;
pub mod app;
