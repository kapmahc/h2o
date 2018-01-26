extern crate env_logger;
extern crate h2o;

fn main() {
    env_logger::init();
    match h2o::console::run() {
        Ok(_) => {}
        Err(e) => println!("{}", e),
    }
}
