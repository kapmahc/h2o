use std::{error, fmt, io, result};

use docopt;
use toml;
use _redis;
use r2d2;
use postgres;
use handlebars;

pub type Result<T> = result::Result<T, Error>;

#[derive(Debug)]
pub enum Error {
    Io(io::Error),
    Docopt(docopt::Error),
    TomlSer(toml::ser::Error),
    TomlDe(toml::de::Error),
    Redis(_redis::RedisError),
    R2d2(r2d2::Error),
    Postgres(postgres::Error),
    HandlebarsTemplateRender(handlebars::TemplateRenderError),
    NotFound,
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            Error::Io(ref err) => err.fmt(f),
            Error::Docopt(ref err) => err.fmt(f),
            Error::TomlSer(ref err) => err.fmt(f),
            Error::TomlDe(ref err) => err.fmt(f),
            Error::Redis(ref err) => err.fmt(f),
            Error::R2d2(ref err) => err.fmt(f),
            Error::Postgres(ref err) => err.fmt(f),
            Error::HandlebarsTemplateRender(ref err) => err.fmt(f),
            Error::NotFound => write!(f, "Not found."),
        }
    }
}

impl error::Error for Error {
    fn description(&self) -> &str {
        match *self {
            Error::Io(ref err) => err.description(),
            Error::Docopt(ref err) => err.description(),
            Error::TomlSer(ref err) => err.description(),
            Error::TomlDe(ref err) => err.description(),
            Error::Redis(ref err) => err.description(),
            Error::R2d2(ref err) => err.description(),
            Error::Postgres(ref err) => err.description(),
            Error::HandlebarsTemplateRender(ref err) => err.description(),
            Error::NotFound => "not found",
        }
    }

    fn cause(&self) -> Option<&error::Error> {
        match *self {
            Error::Io(ref err) => Some(err),
            Error::Docopt(ref err) => Some(err),
            Error::TomlSer(ref err) => Some(err),
            Error::TomlDe(ref err) => Some(err),
            Error::Redis(ref err) => Some(err),
            Error::R2d2(ref err) => Some(err),
            Error::Postgres(ref err) => Some(err),
            Error::HandlebarsTemplateRender(ref err) => Some(err),
            Error::NotFound => None,
        }
    }
}

impl From<io::Error> for Error {
    fn from(err: io::Error) -> Error {
        Error::Io(err)
    }
}

impl From<docopt::Error> for Error {
    fn from(err: docopt::Error) -> Error {
        Error::Docopt(err)
    }
}

impl From<toml::ser::Error> for Error {
    fn from(err: toml::ser::Error) -> Error {
        Error::TomlSer(err)
    }
}

impl From<toml::de::Error> for Error {
    fn from(err: toml::de::Error) -> Error {
        Error::TomlDe(err)
    }
}

impl From<_redis::RedisError> for Error {
    fn from(err: _redis::RedisError) -> Error {
        Error::Redis(err)
    }
}

impl From<r2d2::Error> for Error {
    fn from(err: r2d2::Error) -> Error {
        Error::R2d2(err)
    }
}

impl From<postgres::Error> for Error {
    fn from(err: postgres::Error) -> Error {
        Error::Postgres(err)
    }
}

impl From<handlebars::TemplateRenderError> for Error {
    fn from(err: handlebars::TemplateRenderError) -> Error {
        Error::HandlebarsTemplateRender(err)
    }
}
