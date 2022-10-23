package model

import "time"

type File struct {
    ID          int64      `db:"id" json:"id,omitempty"`
    FilePath    string     `db:"file_path" json:"file_path,omitempty"`
    State       int        `db:"state" json:"state,omitempty"`
    Description string     `db:"description" json:"description,omitempty"`
    CreatedTime int32      `db:"created_time" json:"created_time,omitempty"`
    UpdatedTime int32      `db:"updated_time" json:"updated_time,omitempty"`
    CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
    UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type Member struct {
    Name    string `csv:"name"`
    Address string `csv:"address"`
    Age     int    `csv:"age"`
}

type Class struct {
    Name   string `json:"name" csv:"name"`
    Number int    `json:"number" csv:"number"`
}
