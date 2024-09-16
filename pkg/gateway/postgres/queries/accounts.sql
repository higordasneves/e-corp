-- name: InsertAccount :exec
INSERT INTO accounts (id, document_number, name, secret, balance, created_at)
VALUES (@id, @document_number, @name, @secret, @balance, @created_at);