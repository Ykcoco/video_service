package db

import (
	_"database/sql"
	mydb "go_dev/video_service/db/mysql"
	"fmt"
)


func UserSignup(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("fail to insert, err." + err.Error())
		return false
	} 

	defer stmt.Close()

	ret, err := stmt.Exec(username,password)
	if err != nil {
		fmt.Println("fail to insert, err", err.Error())
		return false
	} 

	if rowsAffected, err := ret.RowsAffected(); err == nil  && rowsAffected > 0 {
		return true
	}
	return true
}

//判断密码是否一致
func UserSignIn(username string, encpwd string) bool {
 stmt, err :=	mydb.DBConn().Prepare("select * from tbl_user where user_name = ? LIMIT 1")
 if err != nil {
	fmt.Println(err.Error())
	return false
 }
 
 defer stmt.Close()
 rows, err := stmt.Query(username)
 if err != nil {
	fmt.Println(err.Error())
	return false
 }else if rows == nil {
	fmt.Println("username name not found "+username)
	return false
 }

 pRows := mydb.ParseRows(rows)
 if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
	 return true
 }
 return false
}

//存放登录的token

func UpdateToken(username string, token string) bool {
   stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token (user_name, user_token) values(?,?)")
   if err != nil {
	   fmt.Println(err.Error())
	   return false
   }

   defer stmt.Close()
   _, err = stmt.Exec(username, token)
   if err != nil {
	fmt.Println(err.Error())
	return false
}
    return true

}

//数据库查询user对应的token信息
func CheckUserToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user_token where user_name = ? LIMIT 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}


	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}else if rows == nil {
		fmt.Println("username name not found "+username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	
	if len(pRows) > 0 && string(pRows[0]["user_token"].([]byte)) == token {
		return true
	}
	return false
}



type User struct {
	Username string `json:"username"`
	Email string  `json:"email"`
	Phone string `json:"phone"`
	SignupAt string `json:"signupat"`
	LastActiveAt string `json:"lastactiveat"`
	Status int `json:"status"`
}

//查询用户信息
func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare(
		"select user_name,signup_at,last_active from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt, &user.LastActiveAt)
	if err != nil {
		return user, err
	}
	return user, nil
}