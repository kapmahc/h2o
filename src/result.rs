use std::{error, fmt, io, result};

use docopt;
use toml;
use time;

pub type Result<T> = result::Result<T, Error>;

#[derive(Debug)]
pub enum Error {
    Io(io::Error),
    Docopt(docopt::Error),
    TomlSer(toml::ser::Error),
    TomlDe(toml::de::Error),
    TimeParse(time::ParseError),
    NotFound,
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            Error::Io(ref err) => err.fmt(f),
            Error::Docopt(ref err) => err.fmt(f),
            Error::TomlSer(ref err) => err.fmt(f),
            Error::TomlDe(ref err) => err.fmt(f),
            Error::TimeParse(ref err) => err.fmt(f),
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
            Error::TimeParse(ref err) => err.description(),
            Error::NotFound => "not found",
        }
    }

    fn cause(&self) -> Option<&error::Error> {
        match *self {
            Error::Io(ref err) => Some(err),
            Error::Docopt(ref err) => Some(err),
            Error::TomlSer(ref err) => Some(err),
            Error::TomlDe(ref err) => Some(err),
            Error::TimeParse(ref err) => Some(err),
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

impl From<time::ParseError> for Error {
    fn from(err: time::ParseError) -> Error {
        Error::TimeParse(err)
    }
}
