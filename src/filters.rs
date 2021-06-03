use std::convert::Infallible;
use warp::{Filter, Rejection, Reply};

use crate::handlers;

pub fn all() -> impl Filter<Extract = impl Reply, Error = Infallible> + Clone {
    health().or(not_found())
}

fn health() -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
    warp::path!("health")
        .and(warp::get())
        .and_then(handlers::health)
}

fn not_found() -> impl Filter<Extract = impl Reply, Error = Infallible> + Clone {
    warp::any().and_then(handlers::not_found)
}
