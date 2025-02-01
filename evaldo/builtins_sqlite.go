//go:build !no_sqlite
// +build !no_sqlite

package evaldo

// import "C"

import (
	"database/sql"

	"github.com/refaktor/rye/env"

	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const MODE_SQLITE = 1
const MODE_PSQL = 2

func SQL_EvalBlock(es *env.ProgramState, mode int, values []any) (*env.ProgramState, []any) {
	var bu strings.Builder
	var str string
	for es.Ser.Pos() < es.Ser.Len() {
		es, str, values = SQL_EvalExpression(es, values, mode)
		bu.WriteString(str + " ")
		//fmt.Println(bu.String())
	}
	//fmt.Println(bu.String())
	es.Res = *env.NewString(bu.String())
	return es, values
}

// mode 1 SQLite, 2 postgresql

func SQL_EvalExpression(es *env.ProgramState, vals []any, mode int) (*env.ProgramState, string, []any) {
	object := es.Ser.Pop()

	switch obj := object.(type) {
	case env.Integer:
		return es, strconv.FormatInt(obj.Value, 10), vals
	case env.String:
		return es, "'" + obj.Value + "'", vals
	/*case env.VoidType:
		es.Res = object
	case env.TagwordType:
		es.Res = object
	case env.UriType:
		es.Res = object */
	case env.Word:
		return es, es.Idx.GetWord(obj.Index), vals
	case env.Opword:
		return es, es.Idx.GetWord(obj.Index)[1:], vals
	case env.Pipeword:
		return es, es.Idx.GetWord(obj.Index)[1:], vals
	case env.Block:
		ser := es.Ser
		es.Ser = obj.Series
		es1, vals1 := SQL_EvalBlock(es, mode, vals)
		//trace("VALS")
		// trace(vals1)
		es.Ser = ser
		//vals2 := append(vals, vals1)
		return es1, "( " + es.Res.(env.String).Value + " )", vals1
	/*case env.GenwordType:
		return EvalGenword(es, object.(env.Genword), nil, false)
	case env.SetwordType:
		return EvalSetword(es, object.(env.Setword)) */
	case env.Getword:
		val, _ := es.Ctx.Get(obj.Index)
		vals = append(vals, resultToJS(val))
		var ph string
		switch mode {
		case 1:
			ph = "?"
		case 2:
			ph = "$" + strconv.Itoa(len(vals))
		}
		return es, ph, vals
	case env.Comma:
		return es, ", ", vals
	default:
		fmt.Println("OTHER SQL NODE")
		return es, "Error 123112431", vals
	}
	// return es, "ERROR", vals
}

var Builtins_sqlite = map[string]*env.Builtin{

	"sqlite-schema//open": {
		Argsn: 1,
		Doc:   "Opening sqlite schema.",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			// arg0.Trace("SQLITE OPEN TODO :::::::::")
			switch str := arg0.(type) {
			case env.Uri:
				// fmt.Println(str.Path)
				db, _ := sql.Open("sqlite3", str.GetPath()) // TODO -- we need to make path parser in URI then this will be path
				return *env.NewNative(ps.Idx, db, "Rye-sqlite")
			default:
				return MakeArgError(ps, 1, []env.Type{env.UriType}, "sqlite-schema//open")
			}

		},
	},

	"htmlize": {
		Argsn: 1,
		Doc:   "Converting to html.",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch str := arg0.(type) {
			case env.Table:
				return *env.NewString(str.ToHtml())
			default:
				return MakeArgError(ps, 1, []env.Type{env.TableType}, "htmlize")
			}

		},
	},

	"Rye-sqlite//exec": {
		Argsn: 2,
		Doc:   "Executes SQL over a database.",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			var sqlstr string
			var vals []any
			switch db1 := arg0.(type) {
			case env.Native:
				switch str := arg1.(type) {
				case env.Block:
					//fmt.Println("BLOCK ****** *****")
					ser := ps.Ser
					ps.Ser = str.Series
					values := make([]any, 0)
					_, vals = SQL_EvalBlock(ps, MODE_SQLITE, values)
					sqlstr = ps.Res.(env.String).Value
					ps.Ser = ser
				case env.String:
					sqlstr = str.Value
				default:
					return MakeArgError(ps, 2, []env.Type{env.BlockType, env.StringType}, "Rye-sqlite//exec")
				}
				if sqlstr != "" {
					db2 := db1.Value.(*sql.DB)
					_, err := db2.Exec(sqlstr, vals...)
					if err != nil {
						return MakeBuiltinError(ps, err.Error(), "Rye-sqlite//exec")
					}
					// rows, err := db1.Value.(*sql.DB).Query(sqlstr, vals...)
					return arg0
				} else {
					return MakeBuiltinError(ps, "sql string not found.", "Rye-sqlite//exec")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "Rye-sqlite//exec")
			}
		},
	},

	"Rye-sqlite//query": {
		Argsn: 2,
		Doc:   "Query a SQLite database with SQL.",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			var sqlstr string
			var vals []any
			switch db1 := arg0.(type) {
			case env.Native:
				switch str := arg1.(type) {
				case env.Block:
					//fmt.Println("BLOCK ****** *****")
					ser := ps.Ser
					ps.Ser = str.Series
					values := make([]any, 0)
					_, vals = SQL_EvalBlock(ps, MODE_SQLITE, values)
					sqlstr = ps.Res.(env.String).Value
					ps.Ser = ser
				case env.String:
					sqlstr = str.Value
				default:
					return MakeArgError(ps, 2, []env.Type{env.BlockType, env.StringType}, "Rye-sqlite//query")
				}
				if sqlstr != "" {
					rows, err := db1.Value.(*sql.DB).Query(sqlstr, vals...)
					if err != nil {
						return MakeBuiltinError(ps, err.Error(), "Rye-sqlite//query")
					}
					columns, _ := rows.Columns()
					spr := env.NewTable(columns)
					// result := make([]map[string]any, 0)
					if err != nil {
						fmt.Println(err.Error())
					} else {
						cols, _ := rows.Columns()
						for rows.Next() {
							var sr env.TableRow
							columns := make([]any, len(cols))
							columnPointers := make([]any, len(cols))
							for i := range columns {
								columnPointers[i] = &columns[i]
							}

							// TODO Scan the result into the column pointers...
							if err := rows.Scan(columnPointers...); err != nil {
								return env.NewError(err.Error())
							}

							// Create our map, and retrieve the value for each column from the pointers slice,
							// storing it in the map with the name of the column as the key.
							m := make(map[string]any)
							for i, colName := range cols {
								val := columnPointers[i].(*any)
								// fmt.Println(val)
								m[colName] = *val
								sr.Values = append(sr.Values, *val)
							}
							spr.AddRow(sr)
							// result = append(result, m)
							// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
						}
						rows.Close() //good habit to close
						//fmt.Println("+++++")
						//fmt.Print(result)
						// return *env.NewNative(ps.Idx, *spr, "Rye-table")
						return *spr
					}
					return MakeBuiltinError(ps, "Empty SQL.", "Rye-sqlite//query")
				} else {
					return MakeBuiltinError(ps, "Sql string not found.", "Rye-sqlite//query")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "Rye-sqlite//query")
			}
		},
	},
}
