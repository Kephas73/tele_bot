package repository

import (
    "botTele/infrastructure/logger"
    "botTele/model"
    "context"
    "fmt"
    "github.com/jmoiron/sqlx"
    "strings"
)

type IFileRepository interface {
    Select(ctx context.Context, size, index int) ([]*model.File, error)
    Insert(ctx context.Context, file *model.File) (*model.File, error)
    Update(ctx context.Context, files ...*model.File) ([]*model.File, error)
    SelectFileByState(ctx context.Context, state int) ([]*model.File, error)
    SelectDynamic(ctx context.Context, queryType int) ([]map[string]interface{}, []string, error)
}

type FileRepository struct {
    *sqlx.DB
}

func NewFileRepository(sqlx *sqlx.DB) IFileRepository {
    fileRepo := FileRepository{
        sqlx,
    }
    return &fileRepo
}

func (repository *FileRepository) Insert(ctx context.Context, file *model.File) (*model.File, error) {

    query := `INSERT INTO file_upload(file_path, state, description, created_time, updated_time) VALUES (:file_path, :state, :description, :created_time, :updated_time)`
    rs, err := repository.NamedExecContext(ctx, query, file)
    if err != nil {
        logger.Error("FileRepository::Insert: - Error: %v", err)
        return nil, err
    }

    file.ID, _ = rs.LastInsertId()
    return file, nil
}

func (repository *FileRepository) Update(ctx context.Context, files ...*model.File) ([]*model.File, error) {

    query := make([]string, 0)
    for _, v := range files {
        query = append(query,
            fmt.Sprintf("UPDATE file_upload SET state = %d, description = '%s', updated_time = %d WHERE id = %d ", v.State, v.Description, v.UpdatedTime, v.ID))
    }
    _, err := repository.ExecContext(ctx, strings.Join(query, ";"))
    if err != nil {
        logger.Error("FileRepository::Update: - Error: %v", err)
        return nil, err
    }

    return files, nil
}

// Select func
// Có thể thêm các lựa chọn chọn state, file_path
func (repository *FileRepository) Select(ctx context.Context, size, index int) ([]*model.File, error) {
    query := `SELECT id, file_path, state, description, created_time, updated_time FROM file_upload ORDER BY created_time DESC LIMIT ? OFFSET ?`

    result := make([]*model.File, 0)
    err := repository.SelectContext(ctx, &result, query, size, size*index)
    if err != nil {
        logger.Error("FileRepository::Select: - Error: %v", err)
        return nil, err
    }

    return result, nil
}

func (repository *FileRepository) SelectFileByState(ctx context.Context, state int) ([]*model.File, error) {
    query := `SELECT id, file_path, state, description, created_time, updated_time FROM file_upload WHERE state = ? ORDER BY created_time ASC`

    result := make([]*model.File, 0)
    err := repository.SelectContext(ctx, &result, query, state)
    if err != nil {
        logger.Error("FileRepository::SelectFileProcessing: - Error: %v", err)
        return nil, err
    }

    return result, nil
}

func (repository *FileRepository) SelectDynamic(ctx context.Context, queryType int) ([]map[string]interface{}, []string, error) {
    query := `SELECT id, file_path, state, description, created_time, updated_time FROM file_upload `

    if queryType == 1 {
        query = `SELECT id, file_path, v2, state, description, created_time, updated_time FROM file_upload_v2 `
    }

    rows, err := repository.QueryxContext(ctx, query)
    if err != nil {
        logger.Error("FileRepository::SelectI: -Query: %s, Error: %v", query, err)
        return nil, nil, err
    }

    return SelectScan(rows)
}

func SelectScan(rows *sqlx.Rows) ([]map[string]interface{}, []string, error) {
    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        return nil, nil, err
    }
    numColumns := len(columns)

    values := make([]interface{}, numColumns)
    for i := range values {
        values[i] = new(string)
    }

    var results []map[string]interface{}
    for rows.Next() {
        if err := rows.Scan(values...); err != nil {
            return nil, nil, err
        }

        dest := make(map[string]interface{}, numColumns)
        for i, column := range columns {
            dest[column] = *(values[i].(*string))
        }
        results = append(results, dest)
    }

    if err = rows.Err(); err != nil {
        return nil, nil, err
    }
    return results, columns, nil
}
