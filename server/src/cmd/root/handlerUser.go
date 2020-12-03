package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"porter/define"
	"porter/wlog"

	"github.com/gin-gonic/gin"
)

func genToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func Login(c *gin.Context) {
	param := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		return
	}

	account := &Account{}
	result := DB.Where("name = ? and password = md5(?)", param.Username, param.Password).First(&account)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.WrongPassword})
		return
	}

	token := genToken()
	result = DB.Model(account).Where(&Account{UID: account.UID}).Update("token", token)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.UpdateAccountErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": define.Success, "data": gin.H{"token": token}})
}

func AccountInfo(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusOK, gin.H{"code": define.IllegalToken})
		return
	}
	account := &Account{}
	result := DB.Model(account).Where("token = ?", token).First(account)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.IllegalToken})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": define.Success, "data": account})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": define.Success})
}
