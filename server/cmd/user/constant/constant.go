package constant

const (
	JWTIssuer  = "FreeCar"
	ThirtyDays = 60 * 60 * 24 * 30
	AccountID  = "accountID"
	ID         = "id"

	ApiConfigPath      = "./config.yaml"
	VideoConfigPath    = "./config.yaml"
	CommentConfigPath  = "./config.yaml"
	FavoriteConfigPath = "./config.yaml"
	MessageConfigPath  = "./config.yaml"
	RelationConfigPath = "./config.yaml"
	UserConfigPath     = "./config.yaml"

	ApiGroup      = "API_GROUP"
	VideoGroup    = "VIDEO_GROUP"
	FavoriteGroup = "FAVORITE_GROUP"
	CommentGroup  = "COMMENT_GROUP"
	RelationGroup = "RELATION_GROUP"
	UserGroup     = "USER_GROUP"
	MessageGroup  = "MESSAGE_GROUP"

	NacosLogDir   = "tmp/nacos/log"
	NacosCacheDir = "tmp/nacos/cache"
	NacosLogLevel = "debug"

	HlogFilePath = "./tmp/hlog/logs/"
	KlogFilePath = "./tmp/klog/logs/"

	MySqlDSN    = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	MongoURI    = "mongodb://%s:%s@%s:%d"
	RabbitMqURI = "amqp://%s:%s@%s:%d/"

	IPFlagName  = "ip"
	IPFlagValue = "0.0.0.0"
	IPFlagUsage = "address"

	PortFlagName  = "port"
	PortFlagUsage = "port"

	TCP = "tcp"

	FreePortAddress  = "localhost:0"
	SERVICEName      = "user"
	DefaultLicNumber = "100000000001"
	DefaultName      = "douyin"
	DefaultGender    = 1
	DefaultBirth     = 631152000000
)
