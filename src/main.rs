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

extern crate h2o;

fn main() {
    match h2o::app::run() {
        Ok(_) => {}
        Err(e) => println!("{}", e),
    }
}
