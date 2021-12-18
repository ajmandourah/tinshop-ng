package repository

// GameID interface
type GameID interface {
	FullID() string
	ShortID() string
	Extension() string
}

// Sources describe all sources type handled
type Sources struct {
	Directories []string `mapstructure:"directories"`
	Nfs         []string `mapstructure:"nfs"`
}

// Config interface
type Config interface {
	RootShop() string
	SetRootShop(string)
	Host() string
	Protocol() string
	Port() int

	DebugNfs() bool
	DebugNoSecurity() bool
	Sources() Sources
	Directories() []string
	NfsShares() []string
	ShopTitle() string
	ShopTemplateData() ShopTemplate
	SetShopTemplateData(ShopTemplate)
}

// ShopTemplate contains all variables used for shop template
type ShopTemplate struct {
	ShopTitle string
}

// HostType new typed string
type HostType string

const (
	// LocalFile Describe local directory file
	LocalFile HostType = "localFile"
	// NFSShare Describe nfs directory file
	NFSShare HostType = "NFS"
)

// FileDesc structure
type FileDesc struct {
	GameID   string
	Size     int64
	GameInfo string
	Path     string
	HostType HostType
}

// GameType structure
type GameType struct {
	Success string                            `json:"success"`
	Titledb map[string]map[string]interface{} `json:"titledb"`
	Files   []interface{}                     `json:"files"`
}
