package tool

import (
	"encoding/base64"
	"fmt"
	"imitate-zhihu/result"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//上传图片，返回素材库图片ID和url
func GraphUpload(path string)(string,string){
	//base64 标准编码
	b,err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	sEnc2 := base64.StdEncoding.EncodeToString(b)
	fnamePos := strings.LastIndex(path,"/")
	fname := path[fnamePos+1:]

	res, err := http.PostForm("http://hn216.api.yesapi.cn/?s=App.CDN.UploadImgByBase64&app_key=3EE1399D7F5DC953E746D33FF9B06E0E&file_name="+fname,
		url.Values{"file":{"data:image/png;base64,"+sEnc2}})
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	str := string(body)
	posUrl := strings.Index(str,"url")
	posId := strings.Index(str,"id")
	id := str[posId+5:posId+6]
	pUrl := str[posUrl+6:posId-3]
	return id,pUrl
}

//根据素材库图片ID删除图片
func GraphDeleteById(id string)  result.Result{
	res, err := http.Post("http://hn216.api.yesapi.cn/?s=App.CDN.DeleteById&app_key=3EE1399D7F5DC953E746D33FF9B06E0E&id="+id,"application/x-www-form-urlencoded",nil)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	str := string(body)
	fmt.Println(str)
	pos := strings.Index(str,"ret")

	ret,_ := StrToInt(str[pos+5:pos+8])
	if ret == 200 {
		return result.Ok
	}
	return result.GraphDeleteErr
}

