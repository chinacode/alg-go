package demo

type MysqlServer struct {
	host     string
	port     int
	user     string
	password string
	database string
}
type BloomServer struct {
	host string
	port int
}
type Config struct {
	mysql MysqlServer
	bloom BloomServer
}

var (
	config = Config{
		//bloom: BloomServer{host: "127.0.0.1", port: 9002},                                                                  // release online
		//mysql: MysqlServer{host: "47.57.68.216", port: 3308, user: "edm", password: "EDM@20201110", database: "edm_crawl"}, //release online

		bloom: BloomServer{host: "8.210.223.207", port: 9002},                                                                 //local test
		mysql: MysqlServer{host: "192.168.1.200", port: 3306, user: "root", password: "Paramida@2019", database: "edm_crawl"}, //local test
	}
)
