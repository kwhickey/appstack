use anyhow::{Context, Result};
use axum::{http::StatusCode, Json, Router, routing::get};
use axum::extract::{Path, State};
use serde::{Deserialize, Serialize};
use sqlx::{Error, sqlite::SqlitePool, sqlite::SqlitePoolOptions};
use thiserror::Error;

use crate::AppError::ItemNotFound;

// Database Models
#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
struct Item {
    id: i64,
    name: String,
    description: String,
}

// Custom Errors
#[derive(Error, Debug)]
enum AppError {
    #[error("Database error: {0}")]
    DatabaseError(#[from] Error),

    #[error("Item not found")]
    ItemNotFound,

    #[error("Internal Error: {0}")]
    Internal(String), // For generic errors or anyhow conversions
}

// This enables using `?` on functions that return `Result<_, anyhow::Error>` to turn them into
// `Result<_, AppError>`. That way you don't need to do that manually.
// impl<E> From<E> for AppError
// where
//     E: Into<anyhow::Error> {
//     fn from(err: E) -> Self {
//         Self(err.into())
//     }
// }

impl From<anyhow::Error> for AppError {
    fn from(err: anyhow::Error) -> Self {
        AppError::Internal(err.to_string())
    }
}

impl axum::response::IntoResponse for AppError {
    fn into_response(self) -> axum::response::Response {
        let (status, error_message) = match self {
            AppError::DatabaseError(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Database error"),
            AppError::Internal(internal_err) => {
                (StatusCode::INTERNAL_SERVER_ERROR, &*internal_err.clone())
            }
            AppError::ItemNotFound => (StatusCode::NOT_FOUND, "Item not found"),
        };

        let body = Json(serde_json::json!({ "error": error_message }));
        (status, body).into_response()
    }
}

// API Routes
async fn get_items(State(pool): State<SqlitePool>) -> Result<Json<Vec<Item>>, AppError> {
    let items = sqlx::query_as::<_, Item>("SELECT * FROM items")
        .fetch_all(&pool)
        .await
        .context("Failed to fetch items")?;

    Ok(Json(items))
}

async fn create_item(
    State(pool): State<SqlitePool>,
    Json(item): Json<Item>,
) -> Result<Json<Item>, AppError> {
    let item = sqlx::query_as::<_, Item>(
        "INSERT INTO items (name, description) VALUES (?, ?) RETURNING *",
    )
    .bind(item.name)
    .bind(item.description)
    .fetch_one(&pool)
    .await
    .context("Failed to insert item")?;

    Ok(Json(item))
}

async fn get_item(
    Path(id): Path<i32>,
    State(pool): State<SqlitePool>,
) -> Result<Json<Item>, AppError> {
    let item = sqlx::query_as::<_, Item>("SELECT * FROM items WHERE id = ?")
        .bind(id)
        .fetch_one(&pool)
        .await
        .context(ItemNotFound)?;

    Ok(Json(item))
}

#[tokio::main]
async fn main() -> Result<()> {
    // initialize tracing
    tracing_subscriber::fmt::init();

    // Setup Database Connection
    let pool = SqlitePoolOptions::new()
        .max_connections(5)
        .connect("sqlite:../items.db")
        .await
        .context("Failed to connect to SQLite database")?;

    sqlx::migrate!()
        .run(&pool)
        .await
        .context("Failed to run database migrations")?;

    // Build Router
    let app = Router::new()
        .route("/items", get(get_items).post(create_item))
        .route("/item/:id", get(get_item))
        .with_state(pool);

    // Start Server
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await?;
    axum::serve(listener, app).await?;

    Ok(())
}
