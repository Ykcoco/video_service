package main

import(
	"fmt"
	"net/http"
   "go_dev/video_service/handle"
   _"go_dev/video_service/db/mysql"

)


func main(){
   
   //静态资源处理
   http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

   http.HandleFunc("/file/upload", handle.HTTPinterceptor(handle.UploadLoadhandle))
   http.HandleFunc("/file/upload/suc", handle.UploadSuchandle)
   http.HandleFunc("/file/meta", handle.GetFileInfo)
   http.HandleFunc("/file/query", handle.FileMetaQueryHandler)
   http.HandleFunc("/file/filefastUpload", handle.HTTPinterceptor(handle.TryfastUploadHandler))//秒传

   http.HandleFunc("/file/download", handle.DownloadFile)
   http.HandleFunc("/file/fileupdate", handle.UpdateFileMetaData)
   http.HandleFunc("/file/filedelete", handle.DeleteFile)

   http.HandleFunc("/user/sigup", handle.UsersignHandle)//用户注册   
   http.HandleFunc("/user/sigin", handle.SigninHandler)//用户登录
   http.HandleFunc("/user/info", handle.HTTPinterceptor(handle.UserinfoHandler))//用户登录
   err := http.ListenAndServe(":1234",nil)
   fmt.Println("端口已监听 1234")  
   if err != nil {
	   fmt.Println("start wrong: ", err)
   }
}