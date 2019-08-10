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

func ItemGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		WriteCPUToDB(cpu)
		body.Close()
	}

	return "ok", nil
}

func CardItemGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var card = common.Card{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		card.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			card.Price = 0
		} else {
			card.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		card.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		card.Tb_link = tb_link
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
			case "显卡类型":
				card.Card_type = paraValue
			case "显卡芯片":
				card.Car_core = paraValue
			case "核心频率":
				card.Core_frequency = paraValue
			case "显存频率":
				card.Gra_mem_frequency = paraValue
			case "显存容量":
				card.Gra_mem_capacity = paraValue
			case "显存位宽":
				card.Gra_mem_bit = paraValue
			case "电源接口":
				card.Power_interface = paraValue
			case "供电模式":
				card.Power_mode = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteCardToDB(card)
		body.Close()
	}

	return "ok", nil
}

func MotherboardGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var motherboard = common.Motherboard{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		motherboard.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			motherboard.Price = 0
		} else {
			motherboard.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		motherboard.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		motherboard.Tb_link = tb_link
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
			case "主芯片组":
				motherboard.Chipset = paraValue
			case "音频芯片":
				motherboard.Audio_chip = paraValue
			case "内存类型":
				motherboard.Ram_type = paraValue
			case "最大内存容量":
				motherboard.Max_ram_size = paraValue
			case "主板板型":
				motherboard.Mother_type = paraValue
			case "外形尺寸":
				motherboard.Shape_size = paraValue
			case "电源插口":
				motherboard.Power_socket = paraValue
			case "供电模式":
				motherboard.Power_mode = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteMotherboardToDB(motherboard)
		body.Close()
	}

	return "ok", nil
}

func MemoryGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var memory = common.Memory{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		memory.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			memory.Price = 0
		} else {
			memory.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		memory.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		memory.Tb_link = tb_link
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
				memory.Pc_type = paraValue
			case "内存容量":
				memory.Capacity = paraValue
			case "内存类型":
				memory.Mem_type = paraValue
			case "内存主频":
				memory.Mem_frequency = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteMemoryToDB(memory)
		body.Close()
	}

	return "ok", nil
}

func HarddriveGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var harddrive = common.Harddrive{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		harddrive.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			harddrive.Price = 0
		} else {
			harddrive.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		harddrive.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		harddrive.Tb_link = tb_link
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
				harddrive.Pc_type = paraValue
			case "硬盘尺寸":
				harddrive.Size = paraValue
			case "硬盘容量":
				harddrive.Capacity = paraValue
			case "单碟容量":
				harddrive.Per_capacity = paraValue
			case "缓存":
				harddrive.Cache = paraValue
			case "转速":
				harddrive.Speed = paraValue
			case "接口类型":
				harddrive.Inter_type = paraValue
			case "接口速率":
				harddrive.Inter_speed = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteToDB(harddrive)
		body.Close()
	}

	return "ok", nil
}

func ChassisGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var chassis = common.Chassis{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		chassis.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			chassis.Price = 0
		} else {
			chassis.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		chassis.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		chassis.Tb_link = tb_link
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
			case "机箱类型":
				chassis.Chassis_type = paraValue
			case "机箱结构":
				chassis.Structure = paraValue
			case "适用主板":
				chassis.Motherboard = paraValue
			case "电源设计":
				chassis.Power_design = paraValue
			case "扩展插槽":
				chassis.Extend_socket = paraValue
			case "前置接口":
				chassis.Preinterface = paraValue
			case "机箱材质":
				chassis.Material = paraValue
			case "板材厚度":
				chassis.Thickness = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteToDB(chassis)
		body.Close()
	}

	return "ok", nil
}

func PowerGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var power = common.Power{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		power.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			power.Price = 0
		} else {
			power.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		power.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		power.Tb_link = tb_link
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
			case "电源类型":
				power.Power_type = paraValue
			case "出线类型":
				power.Out_type = paraValue
			case "额定功率":
				power.Rating_power = paraValue
			case "最大功率":
				power.Max_power = paraValue
			case "主板接口":
				power.Mother_interface = paraValue
			case "硬盘接口":
				power.Hard_interface = paraValue
			case "PFC类型":
				power.Pfc_type = paraValue
			case "转换效率":
				power.Swicth = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteToDB(power)
		body.Close()
	}

	return "ok", nil
}

func CoolingGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var cooling = common.Cooling{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		cooling.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			cooling.Price = 0
		} else {
			cooling.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		cooling.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		cooling.Tb_link = tb_link
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
			case "散热器类型":
				cooling.Cooling_type = paraValue
			case "散热方式":
				cooling.Method = paraValue
			case "适用范围":
				cooling.Use_range = paraValue
			case "输入功率":
				cooling.Input_power = paraValue
			case "风扇尺寸":
				cooling.Size = paraValue
			case "轴承类型":
				cooling.Bear_type = paraValue
			case "转数描述":
				cooling.Revolution = paraValue
			case "噪音":
				cooling.Noise = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteToDB(cooling)
		body.Close()
	}

	return "ok", nil
}

func SSDGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var ssd = common.SSD{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		ssd.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			ssd.Price = 0
		} else {
			ssd.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		ssd.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		ssd.Tb_link = tb_link
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
			case "存储容量":
				ssd.Capacity = paraValue
			case "硬盘尺寸":
				ssd.Hard_size = paraValue
			case "接口类型":
				ssd.Inter_type = paraValue
			case "缓存":
				ssd.Cache = paraValue
			case "读取速度":
				ssd.Read_speed = paraValue
			case "写入速度":
				ssd.Write_speed = paraValue
			case "平均寻道时间":
				ssd.Avg_search_time = paraValue
			case "平均无故障时间":
				ssd.Avg_normal_time = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteToDB(ssd)
		body.Close()
	}

	return "ok", nil
}

func CDDriveGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var cddrive = common.Cddrive{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		cddrive.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			cddrive.Price = 0
		} else {
			cddrive.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		cddrive.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		cddrive.Tb_link = tb_link
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
			case "光驱类型":
				cddrive.Drive_type = paraValue
			case "安装方式":
				cddrive.Install = paraValue
			case "接口类型":
				cddrive.Inter_type = paraValue
			case "缓存容量":
				cddrive.Cache_capacity = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteToDB(cddrive)
		body.Close()
	}

	return "ok", nil
}

func SoundcardGet(in chan string) (string, error) {
	for {
		// time.Sleep(1 * time.Second)
		var str string
		str = <-in
		if str == "NULL" {
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
		var soundcard = common.Soundcard{}

		//查找产品名称
		nameObj := doc.Find("div[class=breadcrumb]>span").Eq(0)
		name := nameObj.Text()
		name = common.GbkToUtf8(name)
		soundcard.Name = name
		output = output + name + " "

		//查找产品价格
		priceObj := doc.Find("b[class=price-type]").Eq(0)
		price := priceObj.Text()
		price = common.GbkToUtf8(price)
		priceint, err := strconv.Atoi(price)
		if err != nil {
			soundcard.Price = 0
		} else {
			soundcard.Price = priceint
		}
		output = output + price + " "

		//查找产品图片img_src
		img_srcObj := doc.Find("div[class=big-pic]>a>img").Eq(0)
		img_src, _ := img_srcObj.Attr("src")
		soundcard.Img_src = img_src
		output = output + img_src + " "

		//查找产品购买链接
		tb_linkObj := doc.Find(".select-mol").Filter(".b2c-jd").Find("a.select-hd").Eq(0)
		tb_link, _ := tb_linkObj.Attr("href")
		soundcard.Tb_link = tb_link
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
			case "声卡类别":
				soundcard.Sound_type = paraValue
			case "适用类型":
				soundcard.Usage_type = paraValue
			case "声道系统":
				soundcard.Sound_system = paraValue
			case "安装方式":
				soundcard.Install = paraValue
			}
			output = output + paraValue + " "
		})

		output = output + "\n"
		fmt.Print(output)
		WriteToDB(soundcard)
		body.Close()
	}

	return "ok", nil
}

func WriteCPUToDB(cpu common.CPU) {
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

func WriteCardToDB(card common.Card) {
	if card.Name == "" && card.Price == 0 {
		return
	}

	sqlStr := fmt.Sprintf(`insert into card (id, name, price, img_src, tb_link, card_type, car_core, 
		core_frequency, gra_mem_frequency, power_interface, power_mode, gra_mem_capacity, gra_mem_bit) 
		values ('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s' )`,
		strconv.FormatInt(time.Now().UnixNano(), 10), card.Name, card.Price, card.Img_src, card.Tb_link,
		card.Card_type, card.Car_core, card.Core_frequency, card.Gra_mem_frequency, card.Power_interface,
		card.Power_mode, card.Gra_mem_capacity, card.Gra_mem_bit,
	)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/miniproject")
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	db.Exec(sqlStr)
}

func WriteMotherboardToDB(motherboard common.Motherboard) {
	if motherboard.Name == "" && motherboard.Price == 0 {
		return
	}

	sqlStr := fmt.Sprintf(`insert into motherboard (id, name, price, img_src, tb_link, chipset, audio_chip, 
		ram_type, max_ram_size, mother_type, shape_size, power_socket, power_mode) 
		values ('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s' )`,
		strconv.FormatInt(time.Now().UnixNano(), 10), motherboard.Name, motherboard.Price, motherboard.Img_src,
		motherboard.Tb_link, motherboard.Chipset, motherboard.Audio_chip, motherboard.Ram_type, motherboard.Max_ram_size,
		motherboard.Mother_type, motherboard.Shape_size, motherboard.Power_socket, motherboard.Power_mode,
	)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/miniproject")
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	db.Exec(sqlStr)
}

func WriteMemoryToDB(memory common.Memory) {
	if memory.Name == "" && memory.Price == 0 {
		return
	}

	sqlStr := fmt.Sprintf(`insert into memory (id, name, price, img_src, tb_link, pc_type, capacity, mem_type, 
		mem_frequency) values ('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%s', '%s') `,
		strconv.FormatInt(time.Now().UnixNano(), 10), memory.Name, memory.Price, memory.Img_src, memory.Tb_link,
		memory.Pc_type, memory.Capacity, memory.Mem_type, memory.Mem_frequency,
	)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/miniproject")
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	db.Exec(sqlStr)
}

func WriteToDB(item interface{}) {
	var sqlStr string
	switch v := item.(type) {
	case common.Harddrive:
		harddrive := v
		sqlStr = fmt.Sprintf(`insert into harddrive (id, name, price, img_src, tb_link, pc_type, size, capacity, 
			per_capacity, cache, speed, inter_type, inter_speed) values ('%s', '%s', '%d', '%s', '%s', '%s', '%s', 
			'%s', '%s', '%s', '%s', '%s', '%s')`,
			strconv.FormatInt(time.Now().UnixNano(), 10), harddrive.Name, harddrive.Price, harddrive.Img_src,
			harddrive.Tb_link, harddrive.Pc_type, harddrive.Size, harddrive.Capacity, harddrive.Per_capacity,
			harddrive.Cache, harddrive.Speed, harddrive.Inter_type, harddrive.Inter_speed,
		)
	case common.Chassis:
		chassis := v
		sqlStr = fmt.Sprintf(`insert into chassis (id, name, price, img_src, tb_link, chassis_type, structure, motherboard, 
			power_design, extend_socket, preinterface, material, thickness) values ('%s', '%s', '%d', '%s', '%s', 
			'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') `,
			strconv.FormatInt(time.Now().UnixNano(), 10), chassis.Name, chassis.Price, chassis.Img_src, chassis.Tb_link,
			chassis.Chassis_type, chassis.Structure, chassis.Motherboard, chassis.Power_design, chassis.Extend_socket,
			chassis.Preinterface, chassis.Material, chassis.Thickness,
		)
	case common.Power:
		power := v
		sqlStr = fmt.Sprintf(`insert into power (id, name, price, img_src, tb_link, power_type, out_type, rating_power, 
			max_power, mother_interface, hard_interface, pfc_type, swicth) values ('%s', '%s', '%d', '%s', '%s', '%s', '%s', 
			'%s', '%s', '%s', '%s', '%s', '%s')`,
			strconv.FormatInt(time.Now().UnixNano(), 10), power.Name, power.Price, power.Img_src, power.Tb_link,
			power.Power_type, power.Out_type, power.Rating_power, power.Max_power, power.Mother_interface, power.Hard_interface,
			power.Pfc_type, power.Swicth,
		)
	case common.Cooling:
		cooling := v
		sqlStr = fmt.Sprintf(`insert into cooling (id, name, price, img_src, tb_link, cooling_type, method, use_range, 
			input_power, size, bear_type, revolution, noise) values ('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%s', '%s', 
			'%s', '%s', '%s', '%s') `,
			strconv.FormatInt(time.Now().UnixNano(), 10), cooling.Name, cooling.Price, cooling.Img_src, cooling.Tb_link,
			cooling.Cooling_type, cooling.Method, cooling.Use_range, cooling.Input_power, cooling.Size, cooling.Bear_type,
			cooling.Revolution, cooling.Noise,
		)
	case common.SSD:
		ssd := v
		sqlStr = fmt.Sprintf(`insert into ssd (id, name, price, img_src, tb_link, capacity, hard_size, inter_type, 
			cache, read_speed, write_speed, avg_normal_time, avg_search_time) values ('%s', '%s', '%d', '%s', '%s', '%s', 
			'%s', '%s', '%s', '%s',	'%s', '%s', '%s') `,
			strconv.FormatInt(time.Now().UnixNano(), 10), ssd.Name, ssd.Price, ssd.Img_src, ssd.Tb_link, ssd.Capacity,
			ssd.Hard_size, ssd.Inter_type, ssd.Cache, ssd.Read_speed, ssd.Write_speed, ssd.Avg_normal_time, ssd.Avg_search_time,
		)
	case common.Cddrive:
		cddrive := v
		sqlStr = fmt.Sprintf(`insert into cddrive (id, name, price, img_src, tb_link, drive_type, install, inter_type, 
			cache_capacity) values('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%s', '%s')`,
			strconv.FormatInt(time.Now().UnixNano(), 10), cddrive.Name, cddrive.Price, cddrive.Img_src, cddrive.Tb_link,
			cddrive.Drive_type, cddrive.Install, cddrive.Inter_type, cddrive.Cache_capacity,
		)
	case common.Soundcard:
		soundcard := v
		sqlStr = fmt.Sprintf(`insert into soundcard (id, name, price, img_src, tb_link, sound_type, usage_type, sound_system, 
			install) values('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%s', '%s')`,
			strconv.FormatInt(time.Now().UnixNano(), 10), soundcard.Name, soundcard.Price, soundcard.Img_src, soundcard.Tb_link,
			soundcard.Sound_type, soundcard.Usage_type, soundcard.Sound_system, soundcard.Install,
		)
		// fmt.Println(sqlStr)
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/miniproject")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	_, err = db.Exec(sqlStr)
	if err != nil {
		fmt.Println(err)
	}

}
