package google

type GmailMessage struct {
	Id           string           `json:"id"`
	ThreadId     string           `json:"threadId"`
	LabelIds     []string         `json:"labelIds"`
	Snippet      string           `json:"snippet"`
	HistoryId    string           `json:"historyId"`
	InternalDate string           `json:"internalDate"`
	Payload      GmailMessagePart `json:"payload"`
	SizeEstimate int              `json:"sizeEstimate"`
	Raw          string           `json:"raw"`
}

type GmailMessagePart struct {
	PartId   string               `json:"partId"`
	MimeType string               `json:"mimeType"`
	FileName string               `json:"filename"`
	Headers  []GmailMessageHeader `json:"headers"`
	Body     GmailMessagePartBody `json:"body"`
	Parts    []GmailMessagePart   `json:"parts"`
}

type GmailMessagePartBody struct {
	AttachmentId string `json:"attachmentId"`
	Size         int    `json:"size"`
	Data         string `json:"data"`
}

type GmailMessageHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GmailMessagesResponse struct {
	Messages           []GmailMessage `json:"messages"`
	ResultSizeEstimate int            `json:"resultSizeEstimate"`
	NextPageToken      string         `json:"nextPageToken"`
}

type GmailSendMessageResponse struct {
	GmailMessage
}
