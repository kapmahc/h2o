use super::redis::Redis;

pub trait Cache {
    fn method(&self) -> String;
}

impl Cache for Redis {
    fn method(&self) -> String {
        return "".to_string();
    }
}
