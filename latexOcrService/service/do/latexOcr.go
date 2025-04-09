package do

type AdminUploadRecord struct {
	System      string
	Scene       string
	FileKye     string
	AdminUserID int64
	IsPublic    int32
	FileName    string
	FileSize    int64
	FileType    string
	ClientIp    string
}
