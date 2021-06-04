use anyhow::Result;
use chrono::prelude::*;
use diesel::prelude::*;
use diesel::{Insertable, QueryDsl, RunQueryDsl};
use uuid::Uuid;

use super::{schema::account_tokens, Account, AccountToken, RepositoryConnection};

#[derive(Insertable)]
#[table_name = "account_tokens"]
struct NewToken<'a> {
    pub token: Uuid,
    pub account_id: Uuid,
    pub user_agent: Option<&'a str>,
}

pub fn generate_token(
    conn: &RepositoryConnection,
    account: &Account,
    user_agent: Option<&str>,
) -> Result<AccountToken> {
    let new_token = NewToken {
        token: Uuid::new_v4(),
        account_id: account.id,
        user_agent: user_agent,
    };

    let token = diesel::insert_into(account_tokens::table)
        .values(&new_token)
        .get_result(conn)?;

    Ok(token)
}

pub fn get_account_by_token(conn: &RepositoryConnection, token_str: &str) -> Result<Account> {
    use super::schema::account_tokens::dsl::*;
    use super::schema::accounts::dsl::*;

    let token_id = Uuid::parse_str(token_str)?;

    conn.transaction(|| {
        let (account, _) = accounts
            .inner_join(account_tokens)
            .filter(token.eq(token_id))
            .first::<(Account, AccountToken)>(conn)?;

        diesel::update(account_tokens.filter(token.eq(token_id)))
            .set(last_used_at.eq(Utc::now()))
            .execute(conn)?;

        Ok(account)
    })
}
