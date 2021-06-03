mod context;

pub use context::Context;

use juniper::{graphql_object, EmptyMutation, EmptySubscription, RootNode};

pub struct Query;

#[graphql_object(
  Context = Context,
)]
impl Query {
    fn apiVersion() -> &str {
        "0.1"
    }
}

type Schema = RootNode<'static, Query, EmptyMutation<Context>, EmptySubscription<Context>>;

pub fn schema() -> Schema {
    Schema::new(
        Query,
        EmptyMutation::<Context>::new(),
        EmptySubscription::<Context>::new(),
    )
}
