package common

type CPU struct {
	Name          string
	Price         int
	Img_src       string
	Tb_link       string
	Pc_type       string
	Cpu_series    string
	Cpu_frequency string
	Max_frequency string
	Lga_type      string
	Second_cache  string
	Core_num      string
	Thread_num    string
	Pack_size     string
}

type Card struct {
	Name              string
	Price             int
	Img_src           string
	Tb_link           string
	Card_type         string
	Car_core          string
	Core_frequency    string
	Gra_mem_frequency string
	Gra_mem_capacity  string
	Gra_mem_bit       string
	Power_interface   string
	Power_mode        string
}

type Motherboard struct {
	Name         string
	Price        int
	Img_src      string
	Tb_link      string
	Chipset      string
	Audio_chip   string
	Ram_type     string
	Max_ram_size string
	Mother_type  string
	Shape_size   string
	Power_socket string
	Power_mode   string
}

type Memory struct {
	Name          string
	Price         int
	Img_src       string
	Tb_link       string
	Pc_type       string
	Capacity      string
	Mem_type      string
	Mem_frequency string
}

type Harddrive struct {
	Name         string
	Price        int
	Img_src      string
	Tb_link      string
	Pc_type      string
	Size         string
	Capacity     string
	Per_capacity string
	Cache        string
	Speed        string
	Inter_type   string
	Inter_speed  string
}

type Chassis struct {
	Name          string
	Price         int
	Img_src       string
	Tb_link       string
	Chassis_type  string
	Structure     string
	Motherboard   string
	Power_design  string
	Extend_socket string
	Preinterface  string
	Material      string
	Thickness     string
}

type Power struct {
	Name             string
	Price            int
	Img_src          string
	Tb_link          string
	Power_type       string
	Out_type         string
	Rating_power     string
	Max_power        string
	Mother_interface string
	Hard_interface   string
	Pfc_type         string
	Swicth           string
}

type Cooling struct {
	Name         string
	Price        int
	Img_src      string
	Tb_link      string
	Cooling_type string
	Method       string
	Use_range    string
	Input_power  string
	Size         string
	Bear_type    string
	Revolution   string
	Noise        string
}

type SSD struct {
	Name            string
	Price           int
	Img_src         string
	Tb_link         string
	Capacity        string
	Hard_size       string
	Inter_type      string
	Cache           string
	Read_speed      string
	Write_speed     string
	Avg_normal_time string
	Avg_search_time string
}

type Cddrive struct {
	Name           string
	Price          int
	Img_src        string
	Tb_link        string
	Drive_type     string
	Install        string
	Inter_type     string
	Cache_capacity string
}

type Soundcard struct {
	Name         string
	Price        int
	Img_src      string
	Tb_link      string
	Sound_type   string
	Usage_type   string
	Sound_system string
	Install      string
}
