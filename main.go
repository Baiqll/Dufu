package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
	"strings"
)


func main() {

	var banner = `
         __      ____     
    ____/ /_  __/ __/_  __
   / __  / / / / /_/ / / /
  / /_/ / /_/ / __/ /_/ / 
  \__,_/\__,_/_/  \__,_/  v1.0
						  
    `
	
	now:=time.Now().Format("2006-01-02 15:04:05")

	var url string
	var wordlists string
	var silent bool


	flag.StringVar(&url, "u", "", "URL")
	flag.StringVar(&wordlists, "w", "", "FUZZ字典")
	flag.BoolVar(&silent, "silent", false, "静默状态")

	// 解析命令行参数写入注册的flag里
	flag.Parse()


	// 静默状态下 不打印banner 信息
	if !silent{
		fmt.Println(string(banner))
		fmt.Println("[*] Starting combination @ ",now)
		fmt.Println("[*] FUZZ URL :\n")
	}

	if wordlists == "" {
		return 
	}


	// 如果管道有参数传递
	if has_stdin() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			if validation(s.Text()){
				combination(s.Text(), wordlists)
			}
			return
		}
	}

	if validation(url){
		combination(url, wordlists)
	}

}

func validation(url string) bool {
	// 验证 url 是否包含有 FUZZ
	return strings.Contains(url, "FUZZ") 
}


func combination(url string, wordlists string){

	// 使用字典
	for _, path := range get_word_dict_list(wordlists) {
		new_url := strings.Replace(url, "FUZZ",strings.Replace(path, "\n", "", -1), -1)
		fmt.Println(new_url)
	}
	
}


// 返回字典列表
func get_word_dict_list(file_path string) []string {

	f, err := os.Open(file_path)
	if err != nil {
		fmt.Println(err.Error())
	}
	//建立缓冲区，把文件内容放到缓冲区中
	buf := bufio.NewReader(f)
	var dict_list []string
	for {
		//遇到\n结束读取
		b, errR := buf.ReadBytes('\n')

		if errR == io.EOF {
			break
		}
		dict_list = append(dict_list, string(b))
	}

	return dict_list

}

// 获取 linux 管道传递的参数
func has_stdin() bool {

	// resList := make([]string, 0, 0)

	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
		return false
	}

	return true
	
}
