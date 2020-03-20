package db

import (
	"database/sql"
	mydb "go_dev/video_service/db/mysql"
	"fmt"
)

// onfileuploadfinished:文件上传完成，保存meta
func Onfileuploadfinished(filehash string,filename string,filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
			"`file_addr`,`status`) values (?,?,?,?,1)")
	
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}

	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Println("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

//tablefile: 文件表结构体
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}


//getfilemeta：从mysql获取文件元信息
func GetFileMeta(filehash string) (*TableFile , error ) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_addr,file_name,file_size from tbl_file where file_sha1 = ? and status = 1 limit 1")
	
	if err != nil {
		fmt.Println(err.Error())
		return nil , err 
	}

	defer stmt.Close()

	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			//查不到对应记录，返回参数以及错误均为nil
			fmt.Println(err)
			return nil, nil 
		}else {
			fmt.Println(err.Error())
			return nil, err 
		}
	}

	return &tfile, nil 
}

//从mysql 中批量获取文件元信息
func GetFileMetaList(limit int) ([]TableFile, error){
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file " +
	"where  status=1 limit ?")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	 cloumns, _ := rows.Columns()
	 values := make([]sql.RawBytes, len(cloumns))

	 var tfiles []TableFile

     for i := 0; i < len(values) && rows.Next(); i ++ {
		 tfile := TableFile{}
		 err = rows.Scan(&tfile.FileHash, &tfile.FileAddr,
			&tfile.FileName, &tfile.FileSize)

			if err != nil {
				fmt.Println(err.Error())
				break
			}
			tfiles = append(tfiles,tfile)
	 }
	 fmt.Println(len(tfiles))
	 return tfiles, nil
}

//updatelocation:更新文件存储地址
func UpdateFileLocation(filehash string, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare("update tbl_file set`file_addr`=? where  `file_sha1`=? limit 1")
    if err != nil {
		fmt.Println("预编译sql失败： "+ err.Error())
        return false
	}

	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("更新文件location失败, filehash:%s", filehash)
		}

		return true
	}

	return false

}



