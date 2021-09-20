package dto


type ActiveList struct {
	Host     string                      `json:"host"`
	Service  string                      `json:"service"`
	FunName  string                      `json:"fun_name"`
	UseTLS   bool                        `json:"use_tls"`
	Restart  bool                        `json:"restart"`
	Body     map[string]interface{}      `json:"body"`
}

type InvRes struct {
	Time      string `json:"timer"`
	Result    string `json:"result"`
	MapResult map[string]interface{} `json:"map_result"`
}

type Desc struct {
	Schema   string `json:"schema"`
	Template string `json:"template"`
	MapSchema   map[string]interface{}  `json:"map_schema"`
	MapTemplate map[string]interface{}  `json:"map_template"`
}