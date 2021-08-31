package global

// 自动上传标志
var AutoUploadFlag bool

// 自动下载标志
var AutoDowndFlag bool

// 数据处理类型
type ActionType int

const (
	UPLOAD ActionType = iota
	DOWNLOAD
	DELETE
)

type ObjectData struct {
	InstanceKey int64
	Key         string
	Type        ActionType
	Path        string
	Count       int // 执行次数
}

var (
	ObjectDataChan chan ObjectData
)
