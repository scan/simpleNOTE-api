table! {
    accounts (id) {
        id -> Uuid,
        email -> Varchar,
        encrypted_password -> Varchar,
        created_at -> Timestamptz,
        updated_at -> Nullable<Timestamptz>,
        last_login -> Nullable<Timestamptz>,
    }
}

table! {
    account_tokens (token) {
        token -> Uuid,
        account_id -> Uuid,
        user_agent -> Nullable<Varchar>,
        created_at -> Timestamptz,
        last_used_at -> Timestamptz,
    }
}

table! {
    notes (id) {
        id -> Uuid,
        account_id -> Uuid,
        title -> Varchar,
        contents -> Nullable<Text>,
        created_at -> Timestamptz,
        updated_at -> Nullable<Timestamptz>,
        deleted_at -> Nullable<Timestamptz>,
    }
}

joinable!(account_tokens -> accounts (account_id));
joinable!(notes -> accounts (account_id));

allow_tables_to_appear_in_same_query!(
    accounts,
    account_tokens,
    notes,
);
