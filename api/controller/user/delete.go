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

// DeletePersonInput : Input
type DeletePersonInput struct {
	Account string `form:"Account" binding:"exists"`
}

// DeletePersonOutput : Output
type DeletePersonOutput struct {
	Persion DeletePersonInput
}

// Delete : Delete API
func Delete(c *gin.Context) {
	res := protocol.Response{}
	var person DeletePersonInput
	var resultCreate isOkResult
	res.Result = &person

	if err := c.Bind(&person); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, res)
		return

	}
	//傳入參數檢查是否存在資料庫
	u, err := SearchDeleteData(person.Account)

	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.SomethingWrongRes(err))
		return
	}

	//資料庫無吻合資料 則 刪除會員
	if u != nil {

		fmt.Println("Delete now....")
		var user User

		fn := "DeleteData"
		dbS := database.GetConn(env.AccountDB)

		sql := " DELETE FROM account_db.user WHERE Account = ? "
		rows, err := dbS.Query(sql, person.Account)

		if err != nil {
			log.Fatalf("Exec Query Failed. fn:%s , err:%s", fn, err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {

			user = User{}

			if err := rows.Scan(

				&user.ID,
				&user.Account,
				&user.Password,
			); err != nil {
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

//SearchDeleteData : Check
func SearchDeleteData(name string) (user *User, err error) {

	fn := "SearchDeleteData"
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
