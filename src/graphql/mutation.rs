use juniper::{graphql_object, FieldResult};

use super::{context::Context, model::Token};
use crate::repository::{generate_token, insert_account, load_and_verify_account};

#[derive(Debug, Copy, Clone)]
pub struct Mutation;

#[graphql_object(Context = Context, )]
impl Mutation {
    fn registerAccount(context: &Context, email: String, password: String) -> FieldResult<bool> {
        let conn = context.pool.get()?;

        insert_account(&conn, email.as_str(), password.as_str())?;

        Ok(true)
    }

    fn login(context: &Context, email: String, password: String) -> FieldResult<Token> {
        let conn = context.pool.get()?;
        let account = load_and_verify_account(&conn, email.as_str(), password.as_str())?;
        let token = generate_token(&conn, &account, None)?;

        Ok(Token(token.token.to_string()))
    }
}
