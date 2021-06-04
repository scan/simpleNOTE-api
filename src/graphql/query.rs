use super::{model::Account, Context};
use anyhow::anyhow;
use juniper::{graphql_object, FieldResult};

pub struct Query;

#[graphql_object(
  Context = Context,
)]
impl Query {
    fn apiVersion() -> &str {
        "0.1"
    }

    fn me(context: &Context) -> FieldResult<Account> {
        let account = context
            .account
            .clone()
            .ok_or(anyhow!("authentication failure"))?;

        Ok(Account {
            email: account.email,
            created_at: account.created_at,
        })
    }
}
