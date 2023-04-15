package types

type FileType int

const (
	UnknowFileType FileType = iota
	DirectoryFileType
	AudioFileType
	TextFileType
	ImageFileType
	VideoFileType
	ArchiveFileType
)

type File struct {
	Id           int      `json:"id"`
	FileType     FileType `json:"fileType"`
	FileName     string   `json:"fileName"`
	Size         string   `json:"size"`
	PathFile     string   `json:"pathFile"`
	StreamPath   string   `json:"streamPath"`
	PathDownload string   `json:"pathDownload"`
}
