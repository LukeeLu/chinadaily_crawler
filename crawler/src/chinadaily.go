package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func httpget1(url string)(result string,err error){
	resp,err1:=http.Get(url)
	if err1 != nil{
		err=err1
		fmt.Println("httpget fail",err)
		return
	}
	buf:=make([]byte,4*1024)
	for{
		n,_:=resp.Body.Read(buf)
		if n==0{
			break
		}
		result +=string(buf[:n])//累加读取的内容
	}
	return
}
//开始爬取每一个新闻title，content
func spidernews(url string)(title,content string,err error){
	result,err1:=httpget1(url)
	if err1 !=nil{
		err=err1
		return
	}
	//fmt.Println(result)
	//os.Exit(1)

	//编译正则
	re1:=regexp.MustCompile(`<h1>(?s:(.*?))</h1>`)
	if re1==nil{
		err=fmt.Errorf("%s","regexp.mustcompile err")
		return
	}
	//取内容
	temtitle :=re1.FindAllStringSubmatch(result,1)
	for _,data := range temtitle{
		title=data[1]
		//fmt.Println(title)

		title=strings.Replace(title,"\t","",-1)//清洗数据
		break
	}

	re2:=regexp.MustCompile(`<div id="Content">(?s:(.*?))</div>`)

	if re2 ==nil{
		err=fmt.Errorf("%s","regexp error")
		return
	}
	re3:=regexp.MustCompile(`</strong>(?s:(.*?))`)

	if re3 ==nil{
		err=fmt.Errorf("%s","regexp error")
		return
	}

	//取内容
	currentcontent:=re2.FindAllStringSubmatch(result,-1)

	for _,data :=range currentcontent{
		content=data[1]
		if content==""{
        content="null"
		}
		//fmt.Println(content)
        //content=strings.TrimLeft(content,"</strong>")

		cleandata:=strings.Index(content,"</figure>")
		//fmt.Println(cleandata)
		if  cleandata==-1{
			content=strings.Replace(content,"\t","",-1)
			content=strings.Replace(content,"<p>","",-1)
			content=strings.Replace(content,"</p>","",-1)
			content=strings.Replace(content,"<strong>","",-1)
			content=strings.Replace(content,"</strong>","",-1)
			break
		}else{
			//content=content[cleandata:]
			content=content[cleandata:]
			content=strings.Replace(content,"</figure>","",-1)
			content=strings.Replace(content,"\t","",-1)
			content=strings.Replace(content,"<p>","",-1)
			content=strings.Replace(content,"</p>","",-1)
			content=strings.Replace(content,"<strong>","",-1)
			content=strings.Replace(content,"</strong>","",-1)
			break
		}


	}
	return

}
//把内容写到文件
func storedata(i int,filetitle,filecontent []string){
	//新建文件
	f,err:=os.Create(strconv.Itoa(i)+".txt")
	if err !=nil{
		fmt.Println("os.ctrate err",err)
		return
	}
	defer f.Close()
	//写内容
	n:=len(filetitle)
	for i:=0;i<n;i++{
		//写标题
		f.WriteString(filetitle[i]+"\n")
		f.WriteString("\n-----------------------------------------------------------------\n")
		//写内容

		f.WriteString(filecontent[i]+"\n")
		//fmt.Println(filecontent)
		f.WriteString("\n=================================================================\n")

	}

}
func spider1(i int, page chan int){
	//明确爬取的url
	//https://www.pengfu.com/xiaohua_1.html
	url:="http://www.chinadaily.com.cn/china/governmentandpolicy/page_" + strconv.Itoa(i)+".html"
	fmt.Printf("正在爬取第%d个网页\n",i,url)

	//开始爬取页面内容
	result,err:=httpget1(url)
	if err !=nil{
		fmt.Println("httpget err")
		return
	}

	re:=regexp.MustCompile(`<h4><a target="_blank" shape="rect" href="//(?s:(.*?))">`)


	if re == nil{
		fmt.Println("regexp err")
		return
	}
	joyurl:=re.FindAllStringSubmatch(result,-1)
	//fmt.Println(joyurl)
	sub:="https://"
	filetitle:=make([]string,0)
	filecontent:=make([]string,0)
	//取网址
	//第一个返回下标，第二个返回内容
	for _,data:=range joyurl{
		data[1]=sub+data[1]
		title,content,err:=spidernews(data[1])//取出来title和content
		if err !=nil{
			fmt.Printf("spiderinjoy error",err)
			continue
		}
		filetitle=append(filetitle,title)//追加内容
		filecontent=append(filecontent,content)

	}
	storedata(i,filetitle,filecontent)
	page<-i  //写内容，写num


}
func work(start,end int){
	fmt.Printf("准备爬取第%d到%d的网址\n",start,end)
	page:=make(chan int)
	for i:=start;i<=end;i++{
		go spider1(i,page)
	}
	for i:=start ;i<=end;i++{
		fmt.Printf("第%d页完成",<-page)
	}
}
func main(){
	var start,end int
	fmt.Printf("输入起始页：")
	fmt.Scan(&start)
	fmt.Printf("输入终止页：")
	fmt.Scan(&end)
	work(start,end)
}






