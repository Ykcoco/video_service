package handle

import (
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"os"
	"time"
	"strconv"
	"encoding/json"
	"go_dev/video_service/meta"
	"go_dev/video_service/util"
	dblayer "go_dev/video_service/db"
)



func UploadLoadhandle(w http.ResponseWriter, r *http.Request) { //处理文件上传
    if r.Method == "GET" {
		fmt.Println("get")
		//返回视频上传页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel read error")
			return 
		}
		
		io.WriteString(w, string(data))
        
	} else if r.Method == "POST"{
		r.ParseForm()
		username := r.Form.Get("username")
	   //接受视频文件
	   file, header ,err := r.FormFile("file")
	   if err != nil {
		   fmt.Println("error of get file %s", err)
		   return
	   }
	  
	   
	   defer file.Close()

	//增添视频文件的原信息
    filemetas := filemeta.FileMeta {
       FileName :  header.Filename,
	   Location : "./tmp/"+ header.Filename,
	   UploadLoadAt: time.Now().Format("2003-01-02 15:04:05"),
	}

	 newfile , err := os.Create(filemetas.Location)
	  if err != nil {
		  fmt.Println("file of create newfile error%s \n",err)
		  return 
	  }

	defer newfile.Close()
	

	//获取sha值
	newfile.Seek(0,0)
    filemetas.FileSha1 = util.FileSha1(newfile)	
	fmt.Println("上传的文件元信息",filemetas)
	
	filesize , err := io.Copy(newfile, file)
	if err != nil {
	fmt.Println("error of save data of file %s \n",err)
	   return
    }
   filemetas.FileSize = filesize


   
	// 将信息存入数据库中(文件表)
	_ = filemeta.UpdateFileDatas(filemetas)
	
	//将信息保存到数据库中（用户文件表）

    suc :=	dblayer.OnUserFileUploadFinished(username,filemetas.FileSha1, filemetas.FileName, filemetas.FileSize) 
	 
	if suc {
		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}else{
        w.Write([]byte("upload failed!"))
 	}

	}
}


//上传成功的
func UploadSuchandle(w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "uploadfinish!!")
}

//查询视频元信息接口 
func GetFileInfo(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	file, err := filemeta.GetfiledataDB(filehash)
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
    }

	data, err := json.Marshal(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}


//批量查询文件元信息接口
func FileMetaQueryHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	file, err := dblayer.QueryUserFileMetas(username,limitCnt)

    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
    }

	data, err := json.Marshal(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}


//下载视频文件接口
func DownloadFile(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	file := filemeta.Getfiledata(filehash)

	files , err := os.Open(file.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	fmt.Println("该文件的地址是; ",files)

	defer files.Close()
	
	data, err := ioutil.ReadAll(files)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}


	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("Content-disposition","attachment;filename=\""+file.FileName+"\"")
    w.Write(data)
}

//重命名问价元信息的接口
func UpdateFileMetaData(w http.ResponseWriter, r *http.Request){
	r.ParseForm()


	op := r.Form.Get("op")
	filehash := r.Form.Get("filehash")
	newname := r.Form.Get("newname")

	if op != "0" {
		w.WriteHeader(http.StatusForbidden)
		return 
	}



	curfilemeta := filemeta.Getfiledata(filehash)
	curfilemeta.FileName = newname
	filemeta.UpdateFileData(curfilemeta)

	data, err := json.Marshal(curfilemeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//删除文件信息
func DeleteFile(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	filemetas := filemeta.Getfiledata(filehash)
	os.Remove(filemetas.Location)
	
	filemeta.RemoveFileMeta(filehash)
	
	w.WriteHeader(http.StatusOK)
}

//文件秒传
//秒传TryfastUploadHandler
func TryfastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	filemeta, err := dblayer.GetFileMeta(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	if filemeta == nil {
		resp := util.RespMsg{
			Code:-1,
			Msg:"秒传失败，请访问普通上传接口",
		}

		w.Write(resp.JSONBytes())
		return 
	}
    
    suc := dblayer.OnUserFileUploadFinished(username ,filehash, filename,int64(filesize))

    if suc {
		resp := util.RespMsg{
			Code:0,
			Msg:"秒传成功",
		}

		w.Write(resp.JSONBytes())
		return 
	}else {
		resp := util.RespMsg{
			Code:-2,
			Msg:"秒传失败，请稍后重试",
		}

		w.Write(resp.JSONBytes())
		return 
	}
}