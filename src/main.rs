// #![feature(plugin)]
// #![plugin(rocket_codegen)]
//
// extern crate rocket;
//
// #[get("/")]
// fn index() -> &'static str {
//     "Hello, world!"
// }
//
// fn main() {
//     rocket::ignite().mount("/", routes![index]).launch();
// }

extern crate env_logger;
extern crate h2o;

fn main() {
    env_logger::init();
    match h2o::console::run() {
        Ok(_) => {}
        Err(e) => println!("{}", e),
    }
}
