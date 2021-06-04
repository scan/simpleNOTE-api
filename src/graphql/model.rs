use chrono::prelude::*;
use juniper::GraphQLObject;

#[derive(Debug, Clone, GraphQLObject, PartialEq, Eq)]
#[graphql(description = "Private account information")]
pub struct Account {
    pub email: String,
    pub created_at: DateTime<Utc>
}
