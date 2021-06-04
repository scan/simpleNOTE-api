use crate::repository::{get_account_by_token, Account, RepositoryPool};
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
                    let account = get_account_by_token(&conn, token).ok()?;

                    Some(account)
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
