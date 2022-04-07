package sqlmock_gorm_demo

import (
	"github.com/jinzhu/gorm"
)

var _DB *gorm.DB

func DB() *gorm.DB {
	return _DB
}

func CreateUser(user *User) (err error) {
	// 打印出要执行的SQL语句 err = DB().Debug().Create(user).Error
	err = DB().Create(user).Error

	return
}

func GetUserById(userId int64) (user *User, err error) {
	user = new(User)
	err = DB().Where("id = ?", userId).First(user).Error

	return
}

func GetUserByNameAndPassword(name, password string) (user *User, err error) {
	user = new(User)
	err = DB().Where("username = ? AND secret = ?", name, password).
		First(&user).Error

	return
}

func UpdateUserNameById(userName string, userId int64) (err error) {
	user := new(User)
	updated := map[string]interface{}{
		"username": userName,
	}
	err = DB().Model(user).Where("id = ?", userId).Updates(updated).Error
	return
}
