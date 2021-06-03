use crate::repository::{Account, RepositoryPool};
use enclose::enclose;

#[derive(Clone)]
pub struct Context {
    pub pool: RepositoryPool,
    pub account: Option<Account>,
}

impl Context {
    pub fn new(pool: RepositoryPool, bearer_token: Option<String>) -> Self {
        let account: Option<Account> = bearer_token
            .as_deref()
            .and_then(parse_bearer_token)
            .and_then(enclose!(
                (pool) | token | {
                    let conn = pool.get().ok()?;
                    // TODO: Load account by token

                    None
                }
            ));

        Context { pool, account }
    }
}

impl juniper::Context for Context {}

fn parse_bearer_token(token: &str) -> Option<&str> {
    if !token.starts_with("Bearer ") {
        return None;
    }

    Some((&token[7..]).trim())
}
