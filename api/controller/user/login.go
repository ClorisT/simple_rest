package user

import (
	"fmt"
	"log"
	"net/http"
	"simple_rest/api/protocol"
	"simple_rest/database"
	"simple_rest/env"

	"github.com/gin-gonic/gin"
)

// LoginPersonInput : Input
type LoginPersonInput struct {
	Account  string `form:"Account" binding:"exists"`
	Password string `form:"Password" binding:"exists"`
}

// Login : Verification API
func Login(c *gin.Context) {
	res := protocol.Response{}
	var input LoginPersonInput

	if err := c.Bind(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, res)
		return

	}

	//傳入輸入參數檢查是否存在資料庫，並比對帳號密碼
	u, err := SearchLoginData(input.Account, input.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.SomethingWrongRes(err))
		return
	}

	//不吻合 則 回傳錯誤信息
	if u == nil {

		res.Result = nil
		res.Code = 2
		res.Message = "Login Failed"
		c.JSON(http.StatusBadRequest, res)
		return

	}

	res.Result = nil
	c.JSON(http.StatusOK, res)

	fmt.Println("--------------------------End--------------------------")

	return
}

//SearchLoginData : Check
func SearchLoginData(name string, password string) (user *User, err error) {
	fn := "SearchData"
	dbS := database.GetConn(env.AccountDB)

	sql := " SELECT id,account,password FROM account_db.user where Account = ? AND Password = ? "
	rows, err := dbS.Query(sql, name, password)
	fmt.Println("-----err:", err)
	if err != nil {
		log.Fatalf("Exec Query Failed. fn:%s , err:%s", fn, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {

		user = &User{}
		fmt.Println("user detail:", user)

		if err := rows.Scan(

			&user.ID,
			&user.Account,
			&user.Password,
		); err != nil {
			log.Fatalf("Fatch Data Error. fn:%s , err:%s", fn, err.Error())
			break
		}
	}

	return
}
