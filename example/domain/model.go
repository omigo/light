package domain

import (
	"time"

	"github.com/arstd/light/example/enum"
)

// Model 模型示例
type Model struct {
	Id    int
	Name  string
	Flag  bool
	Score float32

	Map   map[string]interface{}
	Time  time.Time
	Array []int64 `light:"xarray int[]"`
	Slice []string

	Status enum.Status
	State  enum.State

	Pointer     *Model
	StructSlice []*Model
	Uint32      uint32
}

/*
drop table if exists models;
create table models (
	id serial primary key,
	name text not null,
	flag bool not null default false,
	score decimal(3,1) not null default 0.0,

	map jsonb not null default '{}',
	time timestamptz not null default now(),
	xarray text[] not null,
	slice text[] not null,

	status smallint not null default 0,
	state text not null default '',

	pointer jsonb not null default '{}',
	struct_slice jsonb not null default '[]',
	uint32 bigint not null default 0
)
*/
