package posts

import (
	"math/rand"
	"reflect"
	"strconv"
)

// 随机访问中需要什么来保证随机?
 type serverList struct{
	ipList []string
}

func NewserverList(str ... string)*serverList{
	return &serverList{
		ipLIst :append([]string{},str...)
	}
}

func (s *serverList) AddIP(str ... string){
	append(s.ipList,str...)
}

func Random(str ... string)(string ,error){
   serverList:= NewserverList(str)
   r:=rand.Int(10)
   end:=strconv.Atio((r%len(serverList))) 
   for k,v := range serverList.IPList,{
	if v[10]==end{
		retrun  v,nil
	}
   }
}
