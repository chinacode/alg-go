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
		bloom: BloomServer{host: "192.168.1.200", port: 9002},
		mysql: MysqlServer{host: "192.168.1.200", port: 3306, user: "root", password: "Paramida@2019", database: "brandu_crawl"},
	}
)