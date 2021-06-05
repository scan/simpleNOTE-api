use anyhow::{bail, Result};
use bcrypt::{hash, verify, DEFAULT_COST};
use chrono::Utc;
use diesel::prelude::*;
use uuid::Uuid;

use super::{model::Account, schema::accounts, RepositoryConnection};

#[derive(Debug, Insertable)]
#[table_name = "accounts"]
struct NewAccount<'a> {
    pub id: Uuid,
    pub email: &'a str,
    pub encrypted_password: &'a str,
}

pub fn insert_account<'a>(
    conn: &RepositoryConnection,
    email: &'a str,
    password: &'a str,
) -> Result<Account> {
    conn.transaction(|| {
        let id = Uuid::new_v4();
        let encrypted_password = hash(password, DEFAULT_COST)?;

        let new_account = NewAccount {
            id,
            email,
            encrypted_password: encrypted_password.as_str(),
        };

        diesel::insert_into(accounts::table)
            .values(&new_account)
            .execute(conn)?;

        accounts::table
            .find(id)
            .first(conn)
            .map_err(anyhow::Error::from)
    })
}

pub fn load_and_verify_account(
    conn: &RepositoryConnection,
    email: &str,
    password: &str,
) -> Result<Account> {
    conn.transaction(|| {
        let mut account: Account = match accounts::table
            .filter(accounts::email.eq(email))
            .first(conn)
        {
            Ok(account) => account,
            Err(diesel::result::Error::NotFound) => bail!("email or password wrong"),
            Err(err) => bail!(err),
        };

        if !verify(password, &account.encrypted_password)? {
            bail!("email or password wrong")
        }

        let last_login = Utc::now();
        diesel::update(accounts::table.filter(accounts::id.eq(account.id)))
            .set(accounts::last_login.eq(last_login))
            .execute(conn)?;

        account.last_login = Some(last_login);

        Ok(account)
    })
}

#[allow(dead_code)]
pub fn find_account(conn: &RepositoryConnection, id: Uuid) -> Result<Account> {
    accounts::table
        .find(&id)
        .first(conn)
        .map_err(anyhow::Error::from)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::repository::establish_connection;

    #[test]
    fn test_insert_account() {
        let pool = establish_connection().expect("error creating connection");
        let conn = pool.get().expect("error getting connection");

        conn.begin_test_transaction()
            .expect("could not start transaction");

        let account = insert_account(&conn, "someone@example.com", "password").unwrap();

        assert_eq!(account.email, "someone@example.com");
        assert_ne!(account.encrypted_password, "password");
        assert_eq!(account.last_login, None);

        insert_account(&conn, "someone@example.com", "password")
            .expect_err("second account was created");
    }

    #[test]
    fn test_load_and_verify_account() {
        let pool = establish_connection().expect("error creating connection");
        let conn = pool.get().expect("error fetching connection");

        conn.begin_test_transaction()
            .expect("could not start transaction");

        insert_account(&conn, "someone@example.com", "password").expect("error creating account");

        let account = load_and_verify_account(&conn, "someone@example.com", "password").unwrap();

        assert_eq!(account.email, "someone@example.com");
        assert_ne!(account.last_login, None);

        let res = load_and_verify_account(&conn, "someone@example.com", "password2");
        assert!(res.is_err());
    }

    #[test]
    fn test_find_account() {
        let pool = establish_connection().expect("error creating connection");
        let conn = pool.get().expect("error fetching connection");

        conn.begin_test_transaction()
            .expect("could not start transaction");

        let account = insert_account(&conn, "someone@example.com", "password")
            .expect("error creating account");

        let account = find_account(&conn, account.id).unwrap();

        assert_eq!(account.email, "someone@example.com");

        let res = find_account(&conn, Uuid::new_v4());

        assert!(res.is_err());
    }
}
