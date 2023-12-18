package models

import (
	"fmt"
	"gorm.io/gorm"
	"project/utils"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Salt          string
	LoginTime     *time.Time
	HeartbeatTime *time.Time
	LoginOutTime  *time.Time
	Identity      string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

// FindUserByID 通过ID查找用户
func FindUserByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id=? ", id).First(&user)
	//utils.DB.First(&user, id)
	return user
}

// FindUserByName 通过Name查找
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=? ", name).First(&user)
	return user
}

// FindUserByPhone 通过Phone查找
func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone=?", phone).First(&user)
	return user
}

// FindUserByEmail 通过Email查找
func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email=?", email).First(&user)
}

// FindUserByNameAndPwd 登陆通过Name和Pwd查找
func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name= ? and pass_word=? ", name, password).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}
