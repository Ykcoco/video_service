package filemeta

import (
   mydb "go_dev/video_service/db"
)

type FileMeta struct {
   FileSha1 string `json:"filesha1"`
   FileName string `json:"filename"`
   FileSize int64 `json:"filesize"`
   Location string `json:"location"`
   UploadLoadAt string `json:"uploadat"`
}

var Filemetas map[string]FileMeta


func init() {
   Filemetas = make(map[string]FileMeta)
}

//新增文件原信息到内存
func UpdateFileData(filemeta FileMeta) {
	Filemetas[filemeta.FileSha1] = filemeta
}

//新增文件信息到数据库
func UpdateFileDatas(filemeta FileMeta) bool {
   return  mydb.Onfileuploadfinished(filemeta.FileSha1, filemeta.FileName, filemeta.FileSize, filemeta.Location)
}

//拿到信息(初始，内存信息)
func Getfiledata(filesha string) FileMeta {
   return Filemetas[filesha]
}


//拿到信息（数据库内）
func GetfiledataDB(filesha string) (FileMeta,error) {
   file, err := mydb.GetFileMeta(filesha)
   
   if err != nil {
       return  FileMeta{} ,nil 
   }

   meta := FileMeta{
      FileSha1 : file.FileHash,
      FileName : file.FileName.String,
      FileSize  : file.FileSize.Int64,
      Location  : file.FileAddr.String,
      UploadLoadAt  : file.FileName.String,
   }

   return meta, nil 
}

//删除文件元信息
func RemoveFileMeta(filehash string) {
   delete(Filemetas,filehash)
}