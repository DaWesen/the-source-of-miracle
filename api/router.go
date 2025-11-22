package api

import (
	"md6/dao"
	"md6/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "用户名和密码不能为空",
			"error":   "username and password are required",
		})
		return
	}
	flag := dao.SelectUser(username)
	if flag {
		c.JSON(http.StatusConflict, gin.H{
			"status":  409,
			"message": "用户已存在",
			"error":   "user already exists",
		})
		return
	}
	success := dao.AddUser(username, password)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "创建用户失败",
			"error":   "failed to create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "用户注册成功",
		"data": gin.H{
			"username": username,
		},
	})
}
func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证输入
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "用户名和密码不能为空",
			"error":   "username and password are required",
		})
		return
	}
	flag := dao.SelectUser(username)
	if !flag {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "用户不存在",
			"error":   "user doesn't exist",
		})
		return
	}
	selectPassword := dao.SelectPasswordFromUsername(username)
	if selectPassword != password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "密码错误",
			"error":   "wrong password",
		})
		return
	}
	token, err := utils.GenerateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "生成令牌失败",
			"error":   "failed to generate token",
		})
		return
	}
	c.SetCookie("gin_demo_cookie", username, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "登录成功",
		"data": gin.H{
			"user":  username,
			"token": token,
		},
	})
}

func changePassword(c *gin.Context) {
	username := c.PostForm("username")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	if username == "" || oldPassword == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "所有字段都是必需的",
			"error":   "all fields are required",
		})
		return
	}
	if !dao.SelectUser(username) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "用户不存在",
			"error":   "user doesn't exist",
		})
		return
	}
	currentPassword := dao.SelectPasswordFromUsername(username)
	if currentPassword != oldPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "旧密码错误",
			"error":   "wrong old password",
		})
		return
	}
	success := dao.UpdatePassword(username, newPassword)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "更新密码失败",
			"error":   "failed to update password",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "密码更新成功",
		"data": gin.H{
			"username": username,
		},
	})
}
func getUsers(c *gin.Context) {
	users := dao.GetAllUsers()
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"users": users,
		},
	})
}
func verifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "缺少Authorization头部",
			"error":   "Authorization header is missing",
		})
		return
	}
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "无效或过期的令牌",
			"error":   "Invalid or expired token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "令牌有效",
		"data": gin.H{
			"username": claims.Username,
			"valid":    true,
		},
	})
}
func getUserProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "无法获取用户信息",
			"error":   "failed to get user information",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data": gin.H{
			"username": username,
			"profile": gin.H{
				"name":   username,
				"role":   "user",
				"status": "active",
			},
		},
	})
}
