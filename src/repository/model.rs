use chrono::{DateTime, Utc};
use diesel::Queryable;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize, Queryable, PartialEq, Eq)]
pub struct Account {
  pub id: Uuid,
  pub email: String,
  pub encrypted_password: String,
  pub created_at: DateTime<Utc>,
  pub updated_at: Option<DateTime<Utc>>,
  pub last_login: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Queryable, PartialEq, Eq, PartialOrd, Ord)]
pub struct AccountToken {
    pub token: Uuid,
    pub account_id: Uuid,
    pub user_agent: Option<String>,
    pub created_at: DateTime<Utc>,
    pub last_used_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Queryable, PartialEq, Eq, PartialOrd, Ord)]
pub struct Note {
    pub id: Uuid,
    pub account_id: Uuid,
    pub title: String,
    pub contents: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: Option<DateTime<Utc>>,
    pub deleted_login: Option<DateTime<Utc>>,
}
