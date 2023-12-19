package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"project/models"
	"project/utils"
	"strconv"
	"time"
)

// GetUserList
// @Summary 用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "用户列表！",
		"data":    data,
	})
}

// GetUserByName
// @Summary 根据用户名查询用户
// @Tags 用户模块
// @param name query string false "用户名"
// @Success 200 {string} json{"code","message","data"}
// @Router /user/getUserByName [get]
func GetUserByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "用户名不能为空",
			"data":    nil,
		})
		return
	}

	// 调用查询函数获取用户信息
	user := models.FindUserByName(name)

	if user.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户查询成功",
		"data":    user,
	})
}

// GetUserByPhone
// @Summary 根据手机号查询用户
// @Tags 用户模块
// @param phone query string false "手机号"
// @Success 200 {string} json{"code","message","data"}
// @Router /user/getUserByPhone [get]
func GetUserByPhone(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "手机号不能为空",
			"data":    nil,
		})
		return
	}

	// 调用查询函数获取用户信息
	user := models.FindUserByPhone(phone)

	if user.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户查询成功",
		"data":    user,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}

	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())
	//log.Printf("Creating user: %+v\n", user)
	// 检查用户名、密码、确认密码是否为空
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "用户名、密码或确认密码不能为空！",
			"data":    nil,
		})
		return
	}

	// 检查用户名是否已注册
	if data := models.FindUserByName(user.Name); data.Name != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "用户名已注册！",
			"data":    nil,
		})
		return
	}

	// 检查密码和确认密码是否一致
	if password != repassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
			"data":    nil,
		})
		return
	}

	// 使用密码和盐值生成哈希密码
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt

	// 创建用户
	models.CreateUser(user)
	//log.Println("User created:", user)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "新增用户成功",
		"data":    user,
	})
}

// FindUserByNameAndPwd
// @Summary 登陆
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}

	// name := c.Query("name")
	// password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")

	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失败
			"message": "该用户不存在！",
			"data":    data,
		})
		return
	}

	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失败
			"message": "密码不正确",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	// 密码验证成功，更新登录时间
	currentTime := time.Now()
	utils.DB.Model(&user).Update("LoginTime", &currentTime)
	c.JSON(200, gin.H{
		"code":    0, //  0成功  -1失败
		"message": "登陆成功！",
		"data":    data,
	})

}

// DeleteUser
// @Summary 软删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //  0成功  -1失败
		"message": "删除成功",
		"data":    user,
	})
}

// HardDeleteUser
// @Summary 硬删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/hardDeleteUser [get]
func HardDeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	user.ID = uint(id)

	// 执行硬删除
	result := models.HardDeleteUser(user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "删除失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	//user.PassWord = c.PostForm("password")

	existingUser := models.FindUserByID(user.ID)
	if existingUser.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}
	// 获取新密码，如果有提供的话
	newPassword := c.PostForm("password")

	// 如果新密码非空，说明用户想要更新密码
	if newPassword != "" {
		// 获取用户的盐值，可以从数据库中查询，这里简化为示例
		salt := "some_salt"

		// 使用新密码和盐值生成哈希密码
		user.PassWord = utils.MakePassword(newPassword, salt)
	}
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	//fmt.Println("update:", user)

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, //  0成功  -1失败
			"message": "修改参数不匹配",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, //  0成功  -1失败
			"message": "修改成功",
			"data":    user,
		})
	}

}
