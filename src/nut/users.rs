use std::collections::HashMap;

use rocket_contrib::Template;

#[get("/users/sign-in")]
pub fn get_sign_in() -> Template {
    let mut ctx = HashMap::new();
    ctx.insert("title", "sign in");
    ctx.insert("parent", "layouts/application/index");

    return Template::render("nut/users/sign-in", ctx);
}

#[get("/users/sign-up")]
pub fn get_sign_up() -> &'static str {
    "sign up"
}
