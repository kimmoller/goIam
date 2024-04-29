package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (pg *postgres) getGroupPermissions(ctx context.Context, groupId string) ([]string, error) {
	query := "select system_id from permission join permission_to_group pmt on permission.id = pmt.permission_id" +
		" join permission_group pg on pmt.group_id = pg.id where pg.id = @groupId"

	args := pgx.NamedArgs{
		"groupId": groupId,
	}

	rows, err := pg.db.Query(ctx, query, args)

	if err != nil {
		return nil, fmt.Errorf("error while fetching permissions for group %s, %w", groupId, err)
	}

	var systemIds []string
	for rows.Next() {
		var systemId string
		err := rows.Scan(&systemId)

		if err != nil {
			return nil, fmt.Errorf("error while scanning permissions for group %s, %w", groupId, err)
		}

		systemIds = append(systemIds, systemId)
	}

	return systemIds, nil
}
