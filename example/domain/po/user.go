package po

import (
	"github.com/ericnts/cherry/base"
	"time"
)

type User struct {
	base.RecordPO

	Areas []Area `gorm:"many2many:user_area;" cherry:"auto"`

	CategoryID string    `gorm:"column:category_id;comment:所属类别" cherry:"verify"`
	Username   string    `gorm:"column:username;not null;comment:用户名" cherry:"index"`
	No         int       `gorm:"column:no;comment:编号" cherry:"index:#CategoryID_No"`
	Password   string    `gorm:"column:password;not null;comment:密码"`
	Age        int       `gorm:"column:age;not null;comment:年龄"`
	Birthday   time.Time `gorm:"column:birthday;size:32;not null;comment:生日"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) SetUsername(username string) {
	u.Username = username
	u.Update("username", username)
}

func (u *User) SetNo(no int) {
	u.No = no
	u.Update("no", no)
}

func (u *User) SetPassword(password string) {
	u.Password = password
	u.Update("password", password)
}

func (u *User) SetAge(age int) {
	u.Age = age
	u.Update("age", age)
}

func (u *User) SetBirthday(birthday time.Time) {
	u.Birthday = birthday
	u.Update("birthday", birthday)
}
