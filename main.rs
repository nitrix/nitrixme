#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use]
extern crate rocket;

use chrono::DurationRound;
use rocket::fs::FileServer;
use rocket_dyn_templates::Template;
use std::collections::HashMap;

#[get("/")]
fn homepage() -> Template {
    let mut context = HashMap::new();
    context.insert("email_prefix", gen_email_prefix());
    Template::render("homepage", &context)
}

#[launch]
fn rocket() -> _ {
    rocket::build()
        .attach(Template::fairing())
        .mount("/", routes![homepage])
        .mount("/static", FileServer::from("static"))
}

fn gen_email_prefix() -> String {
    let now = chrono::Utc::now()
        .duration_trunc(chrono::Duration::days(1))
        .unwrap()
        .to_string();

    let prefix = format!("{:x}", md5::compute(now.as_bytes()));
    return prefix[0..16].to_string();
}
