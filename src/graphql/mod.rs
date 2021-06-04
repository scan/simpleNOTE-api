mod context;
mod query;

pub use context::Context;

use juniper::{EmptyMutation, EmptySubscription, RootNode};

type Schema = RootNode<'static, query::Query, EmptyMutation<Context>, EmptySubscription<Context>>;

pub fn schema() -> Schema {
    Schema::new(
        query::Query,
        EmptyMutation::<Context>::new(),
        EmptySubscription::<Context>::new(),
    )
}
