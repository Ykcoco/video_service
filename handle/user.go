package handle

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"go_dev/video_service/util"
	db "go_dev/video_service/db"
)


//用户注册接口
func UsersignHandle(w http.ResponseWriter, r *http.Request) {
   
	if r.Method == http.MethodGet {
	   data , err := ioutil.ReadFile("./static/view/signup.html")
	   if err != nil {
		   fmt.Println(err)
		   w.WriteHeader(http.StatusInternalServerError)
		   return 
	   }

	   io.WriteString(w, string(data))
	   return 
	}
	
	r.ParseForm()
	username := r.Form.Get("username")
	pwd := r.Form.Get("password")
	
	if len(username) <= 3 ||  len(pwd) < 5 {
		w.Write([]byte("username must be >3 and pwd must be > 5"))
	    return 
	}
	
	password := util.Sha1([]byte(util.PWD_SALT+pwd))
	
    suc := db.UserSignup(username,password)
	if  suc {
		w.Write([]byte("SUCCESS"))
	}  else {
        w.Write([]byte("fail"))
	}
}

//用户登录接口
func SigninHandler(w http.ResponseWriter, r *http.Request){

	if r.Method == http.MethodGet {
		data , err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}
 
		io.WriteString(w, string(data))
		return 
	 }
	 
	//1 校验用户名及密码
	r.ParseForm()
	username := r.Form.Get("username")
	pwd := r.Form.Get("password")

	password := util.Sha1([]byte(util.PWD_SALT+pwd))
	pwdChecked := db.UserSignIn(username,password)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	//2 生成访问凭证
	token := util.GenToken(username)
	upRes := db.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}
	//3 登录成功后重定向到主页

	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))

	resp := util.RespMsg{
		Code:0,
		Msg:"ok",
		Data:struct{
			Location  string
			Username  string
			Token string
		}{
			Location:"http://" + r.Host + "/static/view/home.html",
			Username:username,
			Token:token,
		},
	}

	w.Write(resp.JSONBytes())
}

//查询用户数据
func UserinfoHandler(w http.ResponseWriter, r *http.Request){
	//1 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")


	//2 验证token是否有效
	//3 查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return 
	}
	//4 组装并响应用户数据
	resp := util.RespMsg{
	  Code:0,
	  Msg:"ok",
	  Data:user,
	}

    w.Write(resp.JSONBytes())
}

