package library

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// Db router and DB instance
type Db struct {
	DBConn      *sql.Conn
	DB          *sql.DB
	DBConnSlave *sql.Conn
	DBSlave     *sql.DB
	TX          *sql.Tx
	Query       string
	Dialect     string `default:"mysql"`
	Params      []interface{}
	Result      []interface{}
	Context     context.Context
	// New fields for Postgres compatibility
	PrimaryKey string // Set this to "id" or the PK column name to enable RETURNING
	Ignore     bool   // Set to true for INSERT IGNORE / ON CONFLICT DO NOTHING
}

const DbError = "Got error  preparing a.Query %s a.Params %v error %s "

func (a *Db) StartTransaction() error {

	if a.DBConn != nil {

		tx, err := a.DBConn.BeginTx(a.Context, nil)
		if err != nil {

			log.Printf("error starting transaction %s ", err.Error())
			return err
		}

		a.TX = tx

		return nil
	}

	tx, err := a.DB.BeginTx(a.Context, nil)
	if err != nil {

		log.Printf("error starting transaction %s ", err.Error())
		return err
	}

	a.TX = tx

	return nil

}

func (a *Db) Rollback() error {

	if a.TX == nil {

		return fmt.Errorf("Transaction was not started ")

	}

	return a.TX.Rollback()

}

func (a *Db) Commit() error {

	if a.TX == nil {

		return fmt.Errorf("Transaction was not started ")

	}

	return a.TX.Commit()
}

func (a *Db) InsertQueryWithContext() (lastInsertID int64, err error) {

	if a.Dialect == "postgres" {

		var lastInsertId sql.NullInt64

		err = a.DB.QueryRowContext(a.Context, a.Query, a.Params...).Scan(&lastInsertId)
		if err != nil {

			log.Printf(DbError, a.Query, a.Params, err.Error())
			return 0, err
		}

		return lastInsertId.Int64, nil

	}

	var stmt *sql.Stmt

	if a.DBConn != nil {

		stmt, err = a.DBConn.PrepareContext(a.Context, a.Query)
		if err != nil {

			log.Printf(DbError, a.Query, a.Params, err.Error())
			return 0, err
		}

	} else {

		stmt, err = a.DB.PrepareContext(a.Context, a.Query)
		if err != nil {

			log.Printf(DbError, a.Query, a.Params, err.Error())
			return 0, err
		}

	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	return lastInsertId, nil
}

// Deprecated: Use InsertQueryWithContext
func (a *Db) InsertQuery() (lastInsertID int64, err error) {

	a.Context = context.TODO()
	return a.InsertQueryWithContext()

}

func (a *Db) InsertQueryWithContextTx() (lastInsertID int64, err error) {

	if a.Dialect == "postgres" {

		var lastInsertId sql.NullInt64

		err = a.DB.QueryRowContext(a.Context, a.Query, a.Params...).Scan(&lastInsertId)
		if err != nil {

			log.Printf(DbError, a.Query, a.Params, err.Error())
			return 0, err
		}

		return lastInsertId.Int64, nil

	}

	if a.TX == nil {

		if err = a.StartTransaction(); err != nil {

			return 0, err
		}
	}

	stmt, err := a.TX.PrepareContext(a.Context, a.Query)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	return lastInsertId, nil
}

func (a *Db) UpdateQueryWithContext() (rowsAffected int64, err error) {

	var stmt *sql.Stmt

	if a.DBConn != nil {

		stmt, err = a.DBConn.PrepareContext(a.Context, a.Query)
		if err != nil {

			log.Printf(DbError, a.Query, a.Params, err.Error())
			return 0, err
		}

	} else {

		stmt, err = a.DB.PrepareContext(a.Context, a.Query)
		if err != nil {

			log.Printf(DbError, a.Query, a.Params, err.Error())
			return 0, err
		}

	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	rowsaffected, err := res.RowsAffected()
	if err != nil {
		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	return rowsaffected, nil
}

// Deprecated: Use UpdateQueryWithContext
func (a *Db) UpdateQuery() (rowsAffected int64, err error) {

	a.Context = context.TODO()
	return a.UpdateQueryWithContext()

}

func (a *Db) UpdateQueryWithContextTx() (rowsAffected int64, err error) {

	if a.TX == nil {

		if err = a.StartTransaction(); err != nil {

			return 0, err
		}
	}

	stmt, err := a.TX.PrepareContext(a.Context, a.Query)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	rowsaffected, err := res.RowsAffected()
	if err != nil {
		log.Printf(DbError, a.Query, a.Params, err.Error())
		return 0, err
	}

	return rowsaffected, nil
}

func (a *Db) InsertInTransactionWithContext() (lastInsertID *int64, err error) {

	wasNil := false

	if a.TX == nil {

		wasNil = true

		if a.DBConn != nil {

			a.TX, err = a.DBConn.BeginTx(a.Context, nil)
			if err != nil {

				log.Printf("Got error starting transaction %s ", err.Error())
				return nil, err
			}

		} else {

			a.TX, err = a.DB.BeginTx(a.Context, nil)
			if err != nil {

				log.Printf("Got error starting transaction %s ", err.Error())
				return nil, err
			}
		}
	}

	stmt, err := a.TX.PrepareContext(a.Context, a.Query)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		if wasNil {

			_ = a.TX.Rollback()
		}
		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {

		if wasNil {

			_ = a.TX.Rollback()
		}
		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	if wasNil {

		_ = a.TX.Rollback()
	}

	return &lastInsertId, nil
}

// Deprecated: Use InsertInTransactionWithContext
func (a *Db) InsertInTransaction() (lastInsertID *int64, err error) {

	a.Context = context.TODO()
	return a.InsertInTransactionWithContext()
}

func (a *Db) InsertIgnoreWithContext() (lastInsertID *int64, err error) {

	if a.Dialect == "postgres" {

	}

	var stmt *sql.Stmt

	if a.DBConn != nil {

		stmt, err = a.DBConn.PrepareContext(a.Context, a.Query)
		if err != nil {

			log.Printf(DbError, a.Query, a.Params, err.Error())
			return nil, err
		}

	} else {

		stmt, err = a.DB.PrepareContext(a.Context, a.Query)
		if err != nil {

			log.Printf(DbError, a.Query, a.Params, err.Error())
			return nil, err
		}

	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, nil
	}

	return &lastInsertId, nil
}

// Deprecated: Use InsertIgnoreWithContext
func (a *Db) InsertIgnore() (lastInsertID *int64, err error) {

	a.Context = context.TODO()
	return a.InsertIgnoreWithContext()
}

func (a *Db) InsertIgnoreWithContextTx() (lastInsertID *int64, err error) {

	if a.TX == nil {

		if err = a.StartTransaction(); err != nil {

			return nil, err
		}
	}

	stmt, err := a.TX.PrepareContext(a.Context, a.Query)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, nil
	}

	return &lastInsertId, nil
}

// Deprecated: Use InsertIgnoreWithContextTx
func (a *Db) InsertIgnoreInTransactionWithContext() (lastInsertID *int64, err error) {

	stmt, err := a.TX.PrepareContext(a.Context, a.Query)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, nil
	}

	return &lastInsertId, nil
}

// Deprecated: Use UpdateQueryWithContextTx
func (a *Db) UpdateInTransactionWithContext() (rowsAffected *int64, err error) {

	stmt, err := a.TX.PrepareContext(a.Context, a.Query)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(a.Context, a.Params...)
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	rowsaffected, err := res.RowsAffected()
	if err != nil {

		log.Printf(DbError, a.Query, a.Params, err.Error())
		return nil, err
	}

	return &rowsaffected, nil
}

func (a *Db) FetchOneWithContext() *sql.Row {

	if a.Dialect != "postgres" {

		if a.DBConn != nil {

			_, err := a.DBConn.ExecContext(a.Context, "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")
			if err != nil {

				log.Printf("error disabling ONLY_FULL_GROUP_BY %s", err.Error())
			}

		} else {

			_, err := a.DB.ExecContext(a.Context, "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")
			if err != nil {

				log.Printf("error disabling ONLY_FULL_GROUP_BY %s", err.Error())
			}
		}

	}

	a.removeValidParameters()

	if a.Params == nil || len(a.Params) == 0 {

		if a.DBConn != nil {

			return a.DBConn.QueryRowContext(a.Context, a.Query)

		}

		return a.DB.QueryRowContext(a.Context, a.Query)

	}

	if a.DBConn != nil {

		return a.DBConn.QueryRowContext(a.Context, a.Query, a.Params...)

	}

	return a.DB.QueryRowContext(a.Context, a.Query, a.Params...)
}

// Deprecated: Use FetchOneWithContext
func (a *Db) FetchOne() *sql.Row {

	a.Context = context.TODO()
	return a.FetchOneWithContext()
}

func (a *Db) FetchOneSlaveWithContext() *sql.Row {

	if a.Dialect != "postgres" {

		if a.DBConnSlave != nil {

			_, err := a.DBConnSlave.ExecContext(a.Context, "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")
			if err != nil {

				log.Printf("error disabling ONLY_FULL_GROUP_BY %s", err.Error())
			}

		} else {

			_, err := a.DBSlave.ExecContext(a.Context, "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")
			if err != nil {

				log.Printf("error disabling ONLY_FULL_GROUP_BY %s", err.Error())
			}
		}

	}

	a.removeValidParameters()

	if a.Params == nil || len(a.Params) == 0 {

		if a.DBConnSlave != nil {

			return a.DBConnSlave.QueryRowContext(a.Context, a.Query)

		}

		return a.DBSlave.QueryRowContext(a.Context, a.Query)

	}

	if a.DBConnSlave != nil {

		return a.DBConnSlave.QueryRowContext(a.Context, a.Query, a.Params...)

	}

	return a.DBSlave.QueryRowContext(a.Context, a.Query, a.Params...)
}

func (a *Db) FetchWithContext() (*sql.Rows, error) {

	if a.Dialect != "postgres" {

		if a.DBConn != nil {

			_, err := a.DBConn.ExecContext(a.Context, "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")
			if err != nil {

				log.Printf("error disabling ONLY_FULL_GROUP_BY %s", err.Error())
			}

		} else {

			_, err := a.DB.ExecContext(a.Context, "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")
			if err != nil {

				log.Printf("error disabling ONLY_FULL_GROUP_BY %s", err.Error())
			}
		}

	}

	a.removeValidParameters()

	if a.Params == nil || len(a.Params) == 0 {

		if a.DBConn != nil {

			rows, err := a.DBConn.QueryContext(a.Context, a.Query)
			if err != nil {

				log.Printf("error fetching results from database using query %s | no params |  error %s", a.Query, err.Error())
			}

			return rows, err
		}

		rows, err := a.DB.QueryContext(a.Context, a.Query)
		if err != nil {

			log.Printf("error fetching results from database using query %s | no params |  error %s", a.Query, err.Error())
		}

		return rows, err

	}

	if a.DBConn != nil {

		rows, err := a.DBConn.QueryContext(a.Context, a.Query, a.Params...)
		if err != nil {

			log.Printf("error fetching results from database using query %s | params %v |  error %s", a.Query, a.Params, err.Error())
		}

		return rows, err

	}

	rows, err := a.DB.QueryContext(a.Context, a.Query, a.Params...)
	if err != nil {

		log.Printf("error fetching results from database using query %s | params %v |  error %s", a.Query, a.Params, err.Error())
	}

	return rows, err
}

// Deprecated: Use FetchWithContext
func (a *Db) Fetch() (*sql.Rows, error) {

	a.Context = context.TODO()
	return a.FetchWithContext()
}

func (a *Db) FetchSlaveWithContext() (*sql.Rows, error) {

	if a.Dialect != "postgres" {

		if a.DBConnSlave != nil {

			_, err := a.DBConnSlave.ExecContext(a.Context, "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")
			if err != nil {

				log.Printf("error disabling ONLY_FULL_GROUP_BY %s", err.Error())
			}

		} else {

			_, err := a.DBSlave.ExecContext(a.Context, "SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")
			if err != nil {

				log.Printf("error disabling ONLY_FULL_GROUP_BY %s", err.Error())
			}
		}

	}

	a.removeValidParameters()

	if a.Params == nil || len(a.Params) == 0 {

		if a.DBConnSlave != nil {

			rows, err := a.DBConnSlave.QueryContext(a.Context, a.Query)
			if err != nil {

				log.Printf("error fetching results from database using query %s | no params |  error %s", a.Query, err.Error())
			}

			return rows, err

		}

		rows, err := a.DBSlave.QueryContext(a.Context, a.Query)
		if err != nil {

			log.Printf("error fetching results from database using query %s | no params |  error %s", a.Query, err.Error())
		}

		return rows, err

	}

	if a.DBConnSlave != nil {

		rows, err := a.DBConnSlave.QueryContext(a.Context, a.Query, a.Params...)
		if err != nil {

			log.Printf("error fetching results from database using query %s | params %v |  error %s", a.Query, a.Params, err.Error())
		}

		return rows, err

	}

	rows, err := a.DBSlave.QueryContext(a.Context, a.Query, a.Params...)
	if err != nil {

		log.Printf("error fetching results from database using query %s | params %v |  error %s", a.Query, a.Params, err.Error())
	}

	return rows, err
}

func (a *Db) SetParams(params ...interface{}) {

	a.Params = params
}

func (a *Db) SetQuery(query string) {

	a.Query = query
}

func (a *Db) setResults(result ...interface{}) {

	a.Result = result
}

func (a *Db) InsertWithContext(tableName string, data map[string]interface{}) (int64, error) {

	var placeHoldersParts, columns []string
	var params []interface{}

	x := 0

	for column, param := range data {

		x++
		params = append(params, param)
		columns = append(columns, column)
		if a.Dialect == "postgres" {

			placeHoldersParts = append(placeHoldersParts, fmt.Sprintf("$%d", x))

		} else {

			placeHoldersParts = append(placeHoldersParts, "?")

		}
	}

	sqlQueryParts := fmt.Sprintf("INSERT IGNORE INTO %s (%s) %s (%s) ", tableName, strings.Join(columns, ","), a.getValueKeyword(), strings.Join(placeHoldersParts, ","))

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.InsertQueryWithContext()
}

// Deprecated: Use InsertWithContext
func (a *Db) Insert(tableName string, data map[string]interface{}) (int64, error) {

	a.Context = context.TODO()
	return a.InsertWithContext(tableName, data)
}

func (a *Db) InsertWithContextTx(tableName string, data map[string]interface{}) (int64, error) {

	if a.TX == nil {

		if err := a.StartTransaction(); err != nil {

			return 0, err
		}
	}

	var placeHoldersParts, columns []string
	var params []interface{}

	x := 0

	for column, param := range data {

		x++
		params = append(params, param)
		columns = append(columns, column)
		if a.Dialect == "postgres" {

			placeHoldersParts = append(placeHoldersParts, fmt.Sprintf("$%d", x))

		} else {

			placeHoldersParts = append(placeHoldersParts, "?")

		}
	}

	sqlQueryParts := fmt.Sprintf("INSERT IGNORE INTO %s (%s) %s (%s) ", tableName, strings.Join(columns, ","), a.getValueKeyword(), strings.Join(placeHoldersParts, ","))

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.InsertQueryWithContextTx()
}

func (a *Db) UpsertWithContext(tableName string, data map[string]interface{}, updates []string) (int64, error) {

	var placeHoldersParts, updatesPart, columns []string
	var params []interface{}

	x := 0

	for column, param := range data {

		x++
		params = append(params, param)
		columns = append(columns, column)
		if a.Dialect == "postgres" {

			placeHoldersParts = append(placeHoldersParts, fmt.Sprintf("$%d", x))

		} else {

			placeHoldersParts = append(placeHoldersParts, "?")

		}
	}

	updateString := ""

	if updates != nil {

		for _, f := range updates {

			updatesPart = append(updatesPart, fmt.Sprintf("%s=VALUES(%s)", f, f))
		}

		updateString = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s ", strings.Join(updatesPart, ","))
	}

	sqlQueryParts := fmt.Sprintf("INSERT INTO %s (%s) %s (%s) %s", tableName, strings.Join(columns, ","), a.getValueKeyword(), strings.Join(placeHoldersParts, ","), updateString)

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.InsertQueryWithContext()
}

// Deprecated: Use UpsertWithContext
func (a *Db) Upsert(tableName string, data map[string]interface{}, updates []string) (int64, error) {

	a.Context = context.TODO()
	return a.UpsertWithContext(tableName, data, updates)
}

func (a *Db) UpsertWithContextTx(tableName string, data map[string]interface{}, updates []string) (int64, error) {

	if a.TX == nil {

		if err := a.StartTransaction(); err != nil {

			return 0, err
		}
	}

	var placeHoldersParts, updatesPart, columns []string
	var params []interface{}

	x := 0

	for column, param := range data {

		x++
		params = append(params, param)
		columns = append(columns, column)
		if a.Dialect == "postgres" {

			placeHoldersParts = append(placeHoldersParts, fmt.Sprintf("$%d", x))

		} else {

			placeHoldersParts = append(placeHoldersParts, "?")

		}
	}

	updateString := ""

	if updates != nil {

		for _, f := range updates {

			updatesPart = append(updatesPart, fmt.Sprintf("%s=VALUES(%s)", f, f))
		}

		updateString = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s ", strings.Join(updatesPart, ","))
	}

	sqlQueryParts := fmt.Sprintf("INSERT INTO %s (%s) %s (%s) %s", tableName, strings.Join(columns, ","), a.getValueKeyword(), strings.Join(placeHoldersParts, ","), updateString)

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.InsertQueryWithContextTx()
}

func (a *Db) UpdateWithContext(tableName string, andCondition, data map[string]interface{}) (int64, error) {

	var conditions, columns []string
	var params []interface{}

	x := 0
	for column, param := range data {

		x++
		params = append(params, param)
		if a.Dialect == "postgres" {

			columns = append(columns, fmt.Sprintf("%s = $%d ", column, x))

		} else {

			columns = append(columns, fmt.Sprintf("%s = ? ", column))

		}
	}

	for column, value := range andCondition {

		x++
		if a.Dialect == "postgres" {

			conditions = append(conditions, fmt.Sprintf("%s = $%d ", column, x))

		} else {

			conditions = append(conditions, fmt.Sprintf("%s = ? ", column))

		}

		params = append(params, value)

	}

	sqlQueryParts := fmt.Sprintf("UPDATE  %s SET %s WHERE %s ", tableName, strings.Join(columns, ","), strings.Join(conditions, " AND "))

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.UpdateQueryWithContext()
}

// Deprecated: Use UpdateWithContext
func (a *Db) Update(tableName string, andCondition, data map[string]interface{}) (int64, error) {

	a.Context = context.TODO()
	return a.UpdateWithContext(tableName, andCondition, data)
}

func (a *Db) UpdateWithContextTx(tableName string, andCondition, data map[string]interface{}) (int64, error) {

	if a.TX == nil {

		if err := a.StartTransaction(); err != nil {

			return 0, err
		}
	}

	var conditions, columns []string
	var params []interface{}

	x := 0
	for column, param := range data {

		x++
		params = append(params, param)
		if a.Dialect == "postgres" {

			columns = append(columns, fmt.Sprintf("%s = $%d ", column, x))

		} else {

			columns = append(columns, fmt.Sprintf("%s = ? ", column))

		}
	}

	for column, value := range andCondition {

		x++
		if a.Dialect == "postgres" {

			conditions = append(conditions, fmt.Sprintf("%s = $%d ", column, x))

		} else {

			conditions = append(conditions, fmt.Sprintf("%s = ? ", column))

		}

		params = append(params, value)

	}

	sqlQueryParts := fmt.Sprintf("UPDATE  %s SET %s WHERE %s ", tableName, strings.Join(columns, ","), strings.Join(conditions, " AND "))

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.UpdateQueryWithContextTx()
}

func (a *Db) DeleteWithContext(tableName string, andCondition map[string]interface{}) (int64, error) {

	var conditions []string
	var params []interface{}

	x := 0
	for column, value := range andCondition {

		x++
		if a.Dialect == "postgres" {

			conditions = append(conditions, fmt.Sprintf("%s = $%d ", column, x))

		} else {

			conditions = append(conditions, fmt.Sprintf("%s = ? ", column))

		}
		params = append(params, value)
	}

	sqlQueryParts := fmt.Sprintf("DELETE FROM %s WHERE %s ", tableName, strings.Join(conditions, " AND "))

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.UpdateQueryWithContext()
}

// Deprecated: Use DeleteWithContext
func (a *Db) Delete(tableName string, andCondition map[string]interface{}) (int64, error) {

	a.Context = context.TODO()
	return a.DeleteWithContext(tableName, andCondition)
}

func (a *Db) DeleteWithContextTx(tableName string, andCondition map[string]interface{}) (int64, error) {

	if a.TX == nil {

		if err := a.StartTransaction(); err != nil {

			return 0, err
		}
	}

	var conditions []string
	var params []interface{}

	x := 0
	for column, value := range andCondition {

		x++
		if a.Dialect == "postgres" {

			conditions = append(conditions, fmt.Sprintf("%s = $%d ", column, x))

		} else {

			conditions = append(conditions, fmt.Sprintf("%s = ? ", column))

		}
		params = append(params, value)
	}

	sqlQueryParts := fmt.Sprintf("DELETE FROM %s WHERE %s ", tableName, strings.Join(conditions, " AND "))

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.UpdateQueryWithContextTx()
}

func (a *Db) UpsertDataWithContext(tableName string, primaryKey string, data map[string]interface{}, conflicts, updates []string) (int64, error) {

	var placeHoldersParts, updatesPart, columns []string
	var params []interface{}

	x := 0
	for column, param := range data {

		x++
		params = append(params, param)
		columns = append(columns, column)

		if a.Dialect == "postgres" {

			placeHoldersParts = append(placeHoldersParts, fmt.Sprintf("$%d", x))

		} else {

			placeHoldersParts = append(placeHoldersParts, "?")

		}
	}

	updateString := ""

	if updates != nil {

		for _, f := range updates {

			if a.Dialect == "postgres" {

				//excluded.
				updatesPart = append(updatesPart, fmt.Sprintf("%s=excluded.%s", f, f))

			} else {

				updatesPart = append(updatesPart, fmt.Sprintf("%s=VALUES(%s)", f, f))

			}
		}

		if a.Dialect == "postgres" {

			updateString = fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s ", strings.Join(conflicts, ","), strings.Join(updatesPart, ","))

		} else {

			updateString = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s ", strings.Join(updatesPart, ","))

		}
	}

	var sqlQueryParts = ""

	if a.Dialect == "postgres" {

		if len(primaryKey) > 0 {

			sqlQueryParts = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) %s RETURNING %s", tableName, strings.Join(columns, ","), strings.Join(placeHoldersParts, ","), updateString, primaryKey)

		} else {

			sqlQueryParts = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) %s ", tableName, strings.Join(columns, ","), strings.Join(placeHoldersParts, ","), updateString)

		}

	} else {

		sqlQueryParts = fmt.Sprintf("INSERT INTO %s (%s) VALUE (%s) %s", tableName, strings.Join(columns, ","), strings.Join(placeHoldersParts, ","), updateString)

	}

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.InsertQueryWithContext()
}

// Deprecated: Use UpsertDataWithContext
func (a *Db) UpsertData(tableName string, primaryKey string, data map[string]interface{}, conflicts, updates []string) (int64, error) {

	a.Context = context.TODO()
	return a.UpsertDataWithContext(tableName, primaryKey, data, conflicts, updates)
}

func (a *Db) UpsertDataWithContextTx(tableName string, primaryKey string, data map[string]interface{}, conflicts, updates []string) (int64, error) {

	if a.TX == nil {

		if err := a.StartTransaction(); err != nil {

			return 0, err
		}
	}

	var placeHoldersParts, updatesPart, columns []string
	var params []interface{}

	x := 0
	for column, param := range data {

		x++
		params = append(params, param)
		columns = append(columns, column)

		if a.Dialect == "postgres" {

			placeHoldersParts = append(placeHoldersParts, fmt.Sprintf("$%d", x))

		} else {

			placeHoldersParts = append(placeHoldersParts, "?")

		}
	}

	updateString := ""

	if updates != nil {

		for _, f := range updates {

			if a.Dialect == "postgres" {

				//excluded.
				updatesPart = append(updatesPart, fmt.Sprintf("%s=excluded.%s", f, f))

			} else {

				updatesPart = append(updatesPart, fmt.Sprintf("%s=VALUES(%s)", f, f))

			}
		}

		if a.Dialect == "postgres" {

			updateString = fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s ", strings.Join(conflicts, ","), strings.Join(updatesPart, ","))

		} else {

			updateString = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s ", strings.Join(updatesPart, ","))

		}
	}

	var sqlQueryParts = ""

	if a.Dialect == "postgres" {

		if len(primaryKey) > 0 {

			sqlQueryParts = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) %s RETURNING %s", tableName, strings.Join(columns, ","), strings.Join(placeHoldersParts, ","), updateString, primaryKey)

		} else {

			sqlQueryParts = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) %s ", tableName, strings.Join(columns, ","), strings.Join(placeHoldersParts, ","), updateString)

		}

	} else {

		sqlQueryParts = fmt.Sprintf("INSERT INTO %s (%s) VALUE (%s) %s", tableName, strings.Join(columns, ","), strings.Join(placeHoldersParts, ","), updateString)

	}

	a.SetQuery(sqlQueryParts)
	a.SetParams(params...)
	return a.InsertQueryWithContextTx()
}

func (a *Db) getValueKeyword() string {

	if a.Dialect == "postgres" {

		return "VALUES"
	}

	return "VALUES"
}

func (a *Db) removeValidParameters() {

	var par []interface{}

	for _, p := range a.Params {

		if p != nil {

			par = append(par, p)
		}
	}

	a.Params = par

}

// Save executes the query and returns the LastInsertID
func (a *Db) Save() (int64, error) {

	// Postgres uses QueryRow + Scan for RETURNING id
	if a.Dialect == "postgres" && a.PrimaryKey != "" {

		var lastID int64
		var err error

		if a.TX != nil {

			err = a.TX.QueryRowContext(a.Context, a.Query, a.Params...).Scan(&lastID)

		} else {

			err = a.DB.QueryRowContext(a.Context, a.Query, a.Params...).Scan(&lastID)

		}

		if err != nil && err != sql.ErrNoRows {

			return 0, err

		}

		return lastID, nil

	}

	// Default Exec behavior (MySQL)
	var res sql.Result
	var err error

	if a.TX != nil {
		res, err = a.TX.ExecContext(a.Context, a.Query, a.Params...)
	} else {
		res, err = a.DB.ExecContext(a.Context, a.Query, a.Params...)
	}

	if err != nil {
		return 0, err
	}

	// MySQL supports res.LastInsertId()
	if a.Dialect != "postgres" {

		return res.LastInsertId()
	}

	return 0, nil
}
