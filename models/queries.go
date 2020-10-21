package models

var createMetaTableQuery = `
	create table if not exists %s (
		id serial primary key,
		ts timestamptz not null
	)
`

var getMigrationsListQuery = `
	select ts from %s order by ts asc
`
