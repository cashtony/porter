package define

type Task struct {
	Bduss    string
	Nickname string
	Videos   []*TaskVideo
}
type TaskVideo struct {
	AwemeID     string
	Desc        string
	DownloadURL string
}

type TaskFinished struct {
	AwemeID string
}
