package orm

import "database/sql"

// connect 连接数据库
func connect(source string) (*sql.DB, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// genPrepare 生成占位符
func genPrepare(v int) []string {
	temp := make([]string, v)
	for i := 0; i < v; i++ {
		temp[i] = "?"
	}
	return temp
}
