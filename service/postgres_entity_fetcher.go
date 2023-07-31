package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util"
	"crudly/util/result"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type postgresEntityFetcher struct {
	postgres *sql.DB
}

func NewPostgresEntityFetcher(postgres *sql.DB) postgresEntityFetcher {
	return postgresEntityFetcher{
		postgres,
	}
}

func (p *postgresEntityFetcher) FetchEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	tableSchema model.TableSchema,
	id model.EntityId,
) result.R[model.Entity] {
	query := getPostgresEntityQuery(
		projectId,
		tableName,
		id,
	)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Err[model.Entity](fmt.Errorf("error querying postgres: %w", err))
	}

	defer rows.Close()

	if !rows.Next() {
		return result.Err[model.Entity](errs.EntityNotFoundError{})
	}

	return parseEntityFromSqlRow(rows, tableSchema)
}

func (p *postgresEntityFetcher) FetchEntities(
	projectId model.ProjectId,
	tableName model.TableName,
	tableSchema model.TableSchema,
	entityFilter model.EntityFilter,
	entityOrders model.EntityOrders,
	paginationParams model.PaginationParams,
) result.R[model.Entities] {
	query := getPostgresEntitiesQuery(
		projectId,
		tableName,
		entityFilter,
		entityOrders,
		paginationParams,
	)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Err[model.Entities](fmt.Errorf("error querying postgres: %w", err))
	}

	defer rows.Close()

	columnTypes, _ := rows.ColumnTypes()
	columns := make([]any, len(columnTypes))

	entities := model.Entities{}

	for i := range columns {
		columns[i] = new(sql.NullString)
	}

	for rows.Next() {
		entityResult := parseEntityFromSqlRow(rows, tableSchema)

		if entityResult.IsErr() {
			return result.Errf[model.Entities]("error parsing entity: %w", entityResult.UnwrapErr())
		}

		entities = append(entities, entityResult.Unwrap())
	}

	return result.Ok(entities)
}

func parseEntityFromSqlRow(rows *sql.Rows, tableSchema model.TableSchema) result.R[model.Entity] {
	entity := model.Entity{}

	columnTypes, _ := rows.ColumnTypes()
	columns := make([]any, len(columnTypes))

	for i := range columns {
		columns[i] = new(sql.NullString)
	}

	err := rows.Scan(columns...)

	for i, column := range columns {
		nullStr := *(column.(*sql.NullString))

		if !nullStr.Valid {
			continue
		}

		str := nullStr.String

		columnType := *columnTypes[i]

		fieldName := model.FieldName(columnType.Name())

		if fieldName.String() == "id" {
			entity[fieldName] = parsePostgresFieldString(str, model.FieldTypeId)
			continue
		}

		fieldDefinition, ok := tableSchema[fieldName]

		if !ok {
			return result.Errf[model.Entity]("field: %s is not in table schema", fieldName)
		}

		entity[fieldName] = parsePostgresFieldString(str, fieldDefinition.Type)
	}

	if err != nil {
		return result.Errf[model.Entity]("error scaning postgres rows: %w", err)
	}

	return result.Ok(entity)
}

func parsePostgresFieldString(str string, fieldType model.FieldType) any {
	switch fieldType {
	case model.FieldTypeId:
		uuid, err := uuid.Parse(str)

		if err != nil {
			panic(fmt.Sprintf(
				"couldn't parse uuid in sql response, error: %s, val: %s",
				err.Error(),
				str,
			))
		}

		return uuid
	case model.FieldTypeInteger:
		integer, err := strconv.Atoi(str)

		if err != nil {
			panic("couldn't parse integer in sql response")
		}

		return integer
	case model.FieldTypeBoolean:
		if strings.ToLower(str) == "false" {
			return false
		}
		return true
	case model.FieldTypeString:
		return str
	case model.FieldTypeTime:
		time, err := time.Parse(PostgresTimeFormat, str)

		if err != nil {
			panic(fmt.Sprintf("unexpected time format: %s", str))
		}

		return time
	case model.FieldTypeEnum:
		return str
	}

	return nil
}

func getPostgresEntityQuery(projectId model.ProjectId, tableName model.TableName, id model.EntityId) string {
	return "SELECT * FROM \"" + getPostgresTableName(projectId, tableName) + "\" WHERE id = '" + id.String() + "'"
}

func getPostgresFilterString(entityFilter model.EntityFilter) *string {
	if len(entityFilter) == 0 {
		return nil
	}

	filters := ""

	for k, v := range entityFilter {
		filters += fmt.Sprintf(
			"\"%s\" %s %s AND ",
			k,
			getPostgresComparator(v.Type),
			getPostgresFieldValue(v.Comparator).Unwrap(),
		)
	}

	return util.Ptr(strings.TrimSuffix(filters, " AND "))
}

func getPostgresOrderString(orders model.EntityOrders) *string {
	if len(orders) == 0 {
		return nil
	}

	ordersString := ""

	for _, entityOrder := range orders {
		ordersString += fmt.Sprintf(
			"\"%s\" %s,",
			entityOrder.FieldName.String(),
			getPostgresOrder(entityOrder.Type),
		)
	}

	return util.Ptr(strings.TrimSuffix(ordersString, ","))
}

func getPostgresEntitiesQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	entityFilter model.EntityFilter,
	entityOrders model.EntityOrders,
	paginationParams model.PaginationParams,
) string {
	query := fmt.Sprintf("SELECT * FROM \"%s\"", getPostgresTableName(projectId, tableName))

	filtersString := getPostgresFilterString(entityFilter)

	if filtersString != nil {
		query += " WHERE " + *filtersString
	}

	ordersString := getPostgresOrderString(entityOrders)

	if ordersString != nil {
		query += " ORDER BY " + *ordersString
	}

	query += " LIMIT " + paginationParams.Limit.String() +
		" OFFSET " + paginationParams.Offset.String()

	fmt.Println(query)

	return query
}

func getPostgresComparator(fieldFilterType model.FieldFilterType) string {
	switch fieldFilterType {
	case model.FieldFilterTypeEquals:
		return "="
	case model.FieldFilterTypeGreaterThan:
		return ">"
	case model.FieldFilterTypeGreaterThanEq:
		return ">="
	case model.FieldFilterTypeLessThan:
		return "<"
	case model.FieldFilterTypeLessThanEq:
		return "<="
	}
	panic(fmt.Sprintf("invalid field filter type has entered the system: %v", fieldFilterType))
}

func getPostgresOrder(orderType model.FieldOrderType) string {
	switch orderType {
	case model.FieldOrderTypeAscending:
		return "ASC"
	case model.FieldOrderTypeDescending:
		return "DESC"
	}
	panic(fmt.Sprintf("invalid field order type has entered the system: %v", orderType))
}
