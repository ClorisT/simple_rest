package pwd

import (
	"fmt"
	"log"
	"net/http"
	"simple_rest/api/protocol"
	"simple_rest/database"
	"simple_rest/env"

	"github.com/gin-gonic/gin"
)

// ChangePersonInput : Input
type ChangePersonInput struct {
	Account  string `form:"Account" binding:"exists"`
	Password string `form:"Password" binding:"exists"`
}
type isOkResult struct {
	IsOK bool
}

// User : Ouput
type User struct {
	ID       int
	Account  string
	Password string
}

// Change : Update API
func Change(c *gin.Context) {

	res := protocol.Response{}
	var input ChangePersonInput
	res.Result = &input

	var resultCreate isOkResult

	if err := c.Bind(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, res)
		return

	}

	//傳入參數檢查是否存在資料庫
	u, err := SearchChangeData(input.Account)

	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.SomethingWrongRes(err))
		return
	}

	//資料吻合 則 更新資料
	if u != nil {
		fmt.Println("Update now....")
		var user User

		fn := "UpdateData"
		dbS := database.GetConn(env.AccountDB)
		tempID := u.ID
		fmt.Println(tempID)
		sql := " UPDATE account_db.user SET password = ? WHERE id = ? "
		rows, err := dbS.Query(sql, input.Password, tempID)

		if err != nil {
			log.Fatalf("Exec Query Failed. fn:%s , err:%s", fn, err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {

			user = User{}

			if err := rows.Scan(&user.ID, &user.Account, &user.Password); err != nil {

				log.Fatalf("Fatch Data Error. fn:%s , err:%s", fn, err.Error())
				break
			}
		}

		resultCreate.IsOK = true
		res.Result = resultCreate

	} else {

		fmt.Println("No this acc.")
		resultCreate.IsOK = false
		res.Result = resultCreate
	}

	c.JSON(http.StatusOK, res)

	fmt.Println("--------------------------End--------------------------")

	return
}

//SearchChangeData : Check
func SearchChangeData(name string) (user *User, err error) {

	fn := "SearchChangeData"
	dbS := database.GetConn(env.AccountDB)

	sql := " SELECT id,account,password FROM account_db.user where Account = ? "
	rows, err := dbS.Query(sql, name)

	if err != nil {
		log.Fatalf("Exec Query Failed. fn:%s , err:%s", fn, err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {

		user = &User{}

		if err := rows.Scan(&user.ID, &user.Account, &user.Password); err != nil {

			log.Fatalf("Fatch Data Error. fn:%s , err:%s", fn, err.Error())
			break
		}

	}

	return
}
