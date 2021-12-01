package main

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

// sqlmock 是一个实现 sql/driver 的mock库。
// 它不需要建立真正的数据库连接就可以在测试中模拟任何 sql 驱动程序的行为。
// 使用它可以很方便的在编写单元测试的时候mock sql语句的执行结果。

// TestShouldUpdateStats sql执行成功的测试用例
func TestShouldUpdateStats(t *testing.T) {
	// mock一个*sql.DB对象，不需连接真实的数据库
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// mock执行指定SQL语句时的返回结果
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 将mock的DB对象传入我们的函数中
	if err = recordStats(db, 2, 3); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// 确保期望的结果都满足
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestShouldRollbackStatUpdatesOnFailure sql执行失败回滚的测试用例
func TestShouldRollbackStatUpdatesOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").
		WithArgs(2, 3).
		WillReturnError(fmt.Errorf("some mysql error"))
	mock.ExpectRollback()

	// now we execute our method
	if err = recordStats(db, 2, 3); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
