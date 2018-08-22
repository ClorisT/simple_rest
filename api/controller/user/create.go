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

// User : User Database
type User struct {
	ID       int
	Account  string
	Password string
}

type isOkResult struct {
	IsOK bool
}

// CreatePersonInput : Input
type CreatePersonInput struct {
	Account  string `form:"Account" binding:"required"`
	Password string `form:"Password" binding:"required"`
}

// CreatePersonOutput : Output
type CreatePersonOutput struct {
	Persion CreatePersonInput
}

// Creating : Create API
func Creating(c *gin.Context) {
	res := protocol.Response{}
	var person CreatePersonInput
	var resultCreate isOkResult
	res.Result = &person

	if err := c.Bind(&person); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, res)
		return

	}
	//傳入參數檢查是否存在資料庫
	u, err := SearchData(person.Account, person.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.SomethingWrongRes(err))
		return
	}

	//資料庫無吻合資料 則 新增會員
	if u == nil {
		fmt.Println("Insert now....")
		var user User

		fn := "InsertData"
		dbS := database.GetConn(env.AccountDB)

		sql := " INSERT INTO account_db.user( account,password ) VALUES ( ? , ? ) "
		rows, err := dbS.Query(sql, person.Account, person.Password)

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
		fmt.Println("Have same acc, please change your account name.")

		resultCreate.IsOK = false
		res.Result = resultCreate
	}
	c.JSON(http.StatusOK, res)
	fmt.Println("--------------------------End--------------------------")

	return
}

//SearchData : Check
func SearchData(name string, password string) (user *User, err error) {
	fn := "SearchData"
	dbS := database.GetConn(env.AccountDB)

	sql := " SELECT id,account,password FROM account_db.user where Account = ? "
	rows, err := dbS.Query(sql, name)
	fmt.Println("-----err:", err)
	if err != nil {
		log.Fatalf("Exec Query Failed. fn:%s , err:%s", fn, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {

		//sql := " INSERT INTO account_db.user( account,password ) VALUES ( ? , ? )"
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
