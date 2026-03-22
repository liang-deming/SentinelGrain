package quota

// Rule 配额规则（与 api types.QuotaRule 对齐）
type Rule struct {
	AppId      string
	Resource   string
	Threshold  int64
	Period     int64
	UpdateTime int64
}

// ListQuery 列表查询参数
type ListQuery struct {
	AppId    string
	Resource string
	Page     int64
	Size     int64
}
