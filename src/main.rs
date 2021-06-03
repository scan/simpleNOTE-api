#![warn(clippy::all)]

#[macro_use]
extern crate diesel;
#[macro_use]
extern crate diesel_migrations;

mod filters;
mod handlers;

use std::env;
use warp::Filter;

// embed_migrations!();

#[tokio::main]
async fn main() {
    if dotenv::dotenv().is_err() {
        log::warn!("loading environment variabled failed")
    };

    if env::var_os("RUST_LOG").is_none() {
        env::set_var("RUST_LOG", "simple_note_api=info");
    }

    env_logger::init();

    let api = filters::all();

    let routes = api
        .with(warp::log("simple_note_api"))
        .with(
            warp::cors()
                .allow_any_origin()
                .allow_methods(vec!["GET", "POST", "PUT", "DELETE", "OPTIONS"])
                .allow_credentials(true)
                .allow_headers(vec![
                    "Accept",
                    "Authorization",
                    "Content-Type",
                    "X-CSRF-Token",
                    "Accept-Language",
                ])
                .expose_header("Link")
                .max_age(300),
        )
        .with(warp::compression::gzip());

    warp::serve(routes).run(([0, 0, 0, 0], 8080)).await;
}
