use sqlx::{sqlite::SqlitePoolOptions, Error, SqlitePool};

pub async fn initialize_db_pool(uri: String) -> Result<SqlitePool, Error> {
    let pool = SqlitePoolOptions::new()
        .max_connections(1)
        .connect(&uri)
        .await?;

    create_db_tables(&pool).await?;
    Ok(pool)
}

async fn create_db_tables(pool: &SqlitePool) -> Result<(), Error> {
    sqlx::query(
        r#"
        CREATE TABLE IF NOT EXISTS auth_tokens (
            id INTEGER PRIMARY KEY,
            token TEXT NOT NULL
        );
        "#,
    )
    .execute(pool)
    .await?;

    Ok(())
}

pub async fn get_auth_token(pool: &SqlitePool) -> Result<Option<String>, Error> {
    let row: Option<(String,)> = sqlx::query_as("SELECT token FROM auth_tokens LIMIT 1")
        .fetch_optional(pool)
        .await?;

    Ok(row.map(|r| r.0))
}

pub async fn add_auth_token(pool: &SqlitePool) -> Result<(), Error> {
    let token = "68d33db988de6541f6f3".to_string(); // In a real application, generate a secure random token
    sqlx::query("INSERT INTO auth_tokens (token) VALUES (?)")
        .bind(token)
        .execute(pool)
        .await?;
    Ok(())
}
