use enclose::enclose;
use juniper_warp::make_graphql_filter;
use std::convert::Infallible;
use warp::{header, Filter, Rejection, Reply};

use crate::{
    graphql::{schema, Context},
    handlers,
    repository::RepositoryPool,
};

pub fn all(
    db_pool: RepositoryPool,
) -> impl Filter<Extract = impl Reply, Error = Infallible> + Clone {
    health().or(graphql(db_pool)).or(not_found())
}

fn health() -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
    warp::path!("health")
        .and(warp::get())
        .and_then(handlers::health)
}

fn not_found() -> impl Filter<Extract = impl Reply, Error = Infallible> + Clone {
    warp::any().and_then(handlers::not_found)
}

fn graphql(
    db_pool: RepositoryPool,
) -> impl Filter<Extract = impl Reply, Error = Rejection> + Clone {
    let state = warp::any()
        .and(header::optional::<String>("authorization"))
        .map(enclose!((db_pool) move |bearer_token|
      Context::new(db_pool.clone(), bearer_token)));

    let graphql_filter = make_graphql_filter(schema(), state.boxed());

    let post_filter = warp::post()
        .and(warp::body::content_length_limit(1024 * 16))
        .and(graphql_filter.clone());

    let get_filter = warp::get().and(graphql_filter);

    warp::path("graphql").and(get_filter.or(post_filter))
}
