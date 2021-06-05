mod account;
mod model;
mod schema;
mod token;

use anyhow::Result;
use diesel::{
    pg::PgConnection,
    r2d2::{ConnectionManager, PooledConnection},
};
use r2d2::Pool;
use std::env;

pub use account::{find_account, insert_account, load_and_verify_account};
pub use model::{Account, AccountToken, Note};
pub use token::{generate_token, get_account_by_token};

pub type RepositoryPool = Pool<ConnectionManager<PgConnection>>;
pub type RepositoryConnection = PooledConnection<ConnectionManager<PgConnection>>;

fn build_database_url() -> Result<String, env::VarError> {
    let host = env::var("DATABASE_HOST")?;
    let port = env::var("DATABASE_PORT")?;
    let user = env::var("DATABASE_USER")?;
    let pass = env::var("DATABASE_PASSWORD")?;
    let dbname = env::var("DATABASE_NAME")?;

    Ok(format!(
        "postgres://{}:{}@{}:{}/{}?sslmode=require",
        user, pass, host, port, dbname
    ))
}

pub fn establish_connection() -> Result<RepositoryPool> {
    let database_url = env::var("DATABASE_URL").or_else(|_| build_database_url())?;

    let manager = ConnectionManager::<PgConnection>::new(&database_url);
    let pool = Pool::builder().max_size(3).build(manager)?;

    Ok(pool)
}
