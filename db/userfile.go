package db

import (
	mydb "go_dev/video_service/db/mysql"
	"fmt"
)

//用户文件表
type Userfile struct{
	Username string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

func OnUserFileUploadFinished(username ,filehash, filename string, filesize int64) bool {
   stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
	"`file_size`) values (?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username ,filehash, filename,filesize)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]Userfile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_name,file_size,upload_at," +
	"last_update from tbl_user_file where user_name=? limit ?")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer	stmt.Close()
	
	rows, err := stmt.Query(username,limit)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var userfiles []Userfile
	for rows.Next(){
		ufile := Userfile{}
		err := rows.Scan(&ufile.FileHash, &ufile.FileName,&ufile.FileSize,&ufile.UploadAt,&ufile.LastUpdated)
	    if err != nil {
			fmt.Println(err)
			break
		}

		userfiles = append(userfiles,ufile)
	}
	return userfiles,nil
}