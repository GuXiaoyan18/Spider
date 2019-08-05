package spider

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"spiderTool/common"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

func fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong statusCode, %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func ListGet(url string, out chan string) (result string, err error) {
	// resp, err := http.Get(url)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	body, err := fetch(url)
	if err != nil {
		return "", err
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	doc.Find("div[class=pic-mode-box]>ul[class=clearfix]>li>a").Each(func(i int, s *goquery.Selection) {
		str, exists := s.Attr("href")
		if exists && str != "" {
			str = common.GbkToUtf8(str)
			out <- str
		}
	})

	return "ok", nil
}

func ItemGet(in chan string, out chan common.CPU) (result string, err error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
			cpu := common.CPU{}
			out <- cpu
			break
		}

		str = "http://detail.zol.com.cn" + str

		body, err := fetch(str)
		if err != nil {
			return "", err
		}

		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return "", err
		}

		var output string = ""
		var cpu = common.CPU{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		cpu.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			cpu.Price = 0
		} else {
			cpu.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		cpu.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		cpu.Tb_link = tb_link
		output = output + tb_link + " "

		//查找产品参数
		doc.Find("div[class=param-icon]~p").Each(func(i int, s *goquery.Selection) {
			para := s.Text()
			para = common.GbkToUtf8(para)
			//字符串截断，注意是中文需要先转换为rune切片
			rs := []rune(para)
			var paraType, paraValue string
			for index := 0; index < len(rs); index++ {
				if string(rs[index]) == "：" {
					paraType = string(rs[:index])
					paraValue = string(rs[index+1:])
					break
				}
			}
			switch paraType {
			case "适用类型":
				cpu.Pc_type = paraValue
			case "CPU系列":
				cpu.Cpu_series = paraValue
			case "CPU主频":
				cpu.Cpu_frequency = paraValue
			case "最大睿频":
				cpu.Max_frequency = paraValue
			case "插槽类型":
				cpu.Lga_type = paraValue
			case "二级缓存":
				cpu.Second_cache = paraValue
			case "核心数量":
				cpu.Core_num = paraValue
			case "线程数":
				cpu.Thread_num = paraValue
			case "封装大小":
				cpu.Pack_size = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteToDB(cpu)
		body.Close()
	}

	return "ok", nil
}

func WriteToDB(cpu common.CPU) {
	if cpu.Name == "" && cpu.Price == 0 {
		return
	}

	sqlStr := fmt.Sprintf(`insert into cpu (id, name, price, img_src, tb_link, pc_type, cpu_series, cpu_frequency,
		max_frequency, lga_type, second_cache, core_num, thread_num, pack_size) 
		values('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')`,
		strconv.FormatInt(time.Now().UnixNano(), 10), cpu.Name, cpu.Price, cpu.Img_src, cpu.Tb_link, cpu.Pc_type, cpu.Cpu_series,
		cpu.Cpu_frequency, cpu.Max_frequency, cpu.Lga_type, cpu.Second_cache, cpu.Core_num, cpu.Thread_num,
		cpu.Pack_size,
	)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/miniproject")
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	db.Exec(sqlStr)
}
