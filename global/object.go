package global

type ActionType int

const (
	UPLOAD ActionType = iota
	DOWNLOAD
	DELETE
)

type ObjectData struct {
	InstanceKey  int64
	Key          string
	Type         ActionType
	SyncStrategy string
	Path         string
}

var (
	ObjectDataChan chan ObjectData
)
