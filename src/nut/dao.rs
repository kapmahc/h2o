use super::super::db::PostgreSQL;

pub trait Dao {
    fn method(&self) -> String;
}

impl Dao for PostgreSQL {
    fn method(&self) -> String {
        return "do".to_string();
    }
}
