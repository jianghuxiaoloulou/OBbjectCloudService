package setting

import "time"

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type GeneralSettingS struct {
	LogSavePath string
	LogFileName string
	LogFileExt  string
	LogMaxSize  int
	LogMaxAge   int
	MaxThreads  int
	MaxTasks    int
}

type DatabaseSettingS struct {
	DBConn       string
	DBType       string
	MaxIdleConns int
	MaxOpenConns int
}

type ObjectSettingS struct {
	OBJECT_BucketId       string
	OBJECT_MDID           string
	OBJECT_Sync           string
	OBJECT_GET_Version    string
	OBJECT_POST_Upload    string
	OBJECT_GET_Download   string
	OBJECT_DEL_Delete     string
	OBJECT_PATH           string
	TOKEN_Username        string
	TOKEN_Password        string
	TOKEN_URL             string
	DOWN_Dest_Code        int
	DOWN_Dest_Root        string
	UPLOAD_ROOT           string
	OBJECT_Upload_SUCCESS int
	OBJECT_Upload_ERROR   int
	OBJECT_Down_SUCCESS   int
	OBJECT_Down_ERROR     int
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
