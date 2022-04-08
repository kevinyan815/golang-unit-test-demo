package sqlmock_gorm_demo

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	mock sqlmock.Sqlmock
	err  error
	db   *sql.DB
)

// TestMain是在当前package下，最先运行的一个函数，常用于初始化
func TestMain(m *testing.M) {
	// 创建一个 sqlmock 的数据库连接 和 mock对象，mock对象管理 db 预期要执行的SQL

	// sqlmock 默认使用 sqlmock.QueryMatcherRegex 作为默认的SQL匹配器
	// 该匹配器使用mock.ExpectQuery 和 mock.ExpectExec 的参数作为正则表达式而不是SQL语句字符串
	// 所以指定预期要执行的SQL语句时，会遇到下面的错误
	// ExecQuery: error parsing regexp: invalid or unsupported Perl syntax
	// 有两种办法解决这个问题：
	// 1. 使用regexp.QuoteMeta 把SQL转义成正则表达式 => mock.ExpectQuery(regexp.QuoteMeta(`SELECT ....`))
	// 2. 让sqlmock 使用 QueryMatcherEqual 匹配器，该匹配器把mock.ExpectQuery 和 mock.ExpectExec 的参数作为
	//    预期要执行的SQL语句跟要执行的SQL进行相等比较

	// 这里 让sqlmock 使用 sqlmock.QueryMatcherEqual 作为匹配器匹配器，该匹配器把mock.ExpectQuery
	// 和 mock.ExpectExec 的参数作为 预期要执行的SQL语句跟真正要执行的SQL进行相等比较
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {

		panic(err)
	}
	_DB, err = gorm.Open("mysql", db)

	// m.Run 是调用包下面各个Test函数的入口
	os.Exit(m.Run())
}

func TestCreateUserMock(t *testing.T) {
	user := &User{
		UserName:  "Kevin",
		Secret:    "123456",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mock.ExpectBegin()
	// 因为 sqlmock 使用的是 QueryMatcherEqual 匹配器
	// 所以，下面的sql语句必须精确匹配要执行的SQL（包括符号和空格）
	mock.ExpectExec("INSERT INTO `users` (`username`,`secret`,`created_at`,`updated_at`) VALUES (?,?,?,?)").
		WithArgs(user.UserName, user.Secret, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := CreateUser(user)
	assert.Nil(t, err)

}

func TestGetUserByNameAndPasswordMock(t *testing.T) {
	user := &User{
		Id:        1,
		UserName:  "Kevin",
		Secret:    "123456",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery("SELECT * FROM `users`  WHERE (username = ? AND secret = ?) "+
		"ORDER BY `users`.`id` ASC LIMIT 1").
		WithArgs(user.UserName, user.Secret).
		WillReturnRows(
			// 这里要跟结果集包含的列匹配，因为查询是 SELECT * 所以表的字段都要列出来
			sqlmock.NewRows([]string{"id", "username", "secret", "created_at", "updated_at"}).
				AddRow(1, user.UserName, user.Secret, user.CreatedAt, user.UpdatedAt))
	res, err := GetUserByNameAndPassword(user.UserName, user.Secret)
	assert.Nil(t, err)
	assert.Equal(t, user, res)
}

//func TestUpdateUserNameByIdMock(t *testing.T) {
//	newName := "Kev"
//	var userId int64 = 1
//
//	// GORM 在UPDATE 的时候会自动更新updated_at 字段为当前时间，与这里withArgs传递的 time.Now() 参数不一致，
//	// 目前没有办法Mock测试GORM的UPDATE方法，除非用Exec直接执行更新SQL，不过那就失去使用ORM的意义了
//	// 这个先跳过
//	mock.ExpectBegin()
//	mock.ExpectExec("UPDATE `users` SET `updated_at` = ?, `username` = ?  WHERE (id = ?)").
//		WithArgs(time.Now(), newName, userId).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectCommit()
//
//	err := UpdateUserNameById(newName, userId)
//	assert.Nil(t, err)
//	// 调用确保期望的结果都满足的方法可以通过，但还是会告警updated_at 字段实际执行的值与预期指定的值不一致。
//	//mock.ExpectationsWereMet()
//}
