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
	CronSpec    string
	WorkMode    int
}

type DatabaseSettingS struct {
	DBConn       string
	DBType       string
	MaxIdleConns int
	MaxOpenConns int
}

type ObjectSettingS struct {
	OBJECT_BucketId                 string
	OBJECT_ResId                    string
	OBJECT_MDID                     string
	OBJECT_AK                       string
	OBJECT_Sync                     string
	OBJECT_GET_Version              string
	OBJECT_POST_Upload              string
	OBJECT_GET_Download             string
	OBJECT_DEL_Delete               string
	OBJECT_PATH                     string
	TOKEN_Username                  string
	TOKEN_Password                  string
	TOKEN_URL                       string
	DOWN_Dest_Code                  int
	DOWN_Dest_Root                  string
	UPLOAD_ROOT                     string
	OBJECT_Upload_Success_Code      int
	OBJECT_Upload_SUCCESS           int
	OBJECT_Upload_ERROR             int
	OBJECT_Down_SUCCESS             int
	OBJECT_Down_ERROR               int
	OBJECT_Upload_Flag              int
	OBJECT_Down_Flag                int
	OBJECT_TASK                     int
	OBJECT_Count                    int
	File_Fragment_Size              int64
	Each_Section_Size               int64
	File_Split_Temp                 string
	OBJECT_Multipart_Init_URL       string
	OBJECT_Multipart_Upload_URL     string
	OBJECT_Multipart_Completion_URL string
	OBJECT_Multipart_Abortion_URL   string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
