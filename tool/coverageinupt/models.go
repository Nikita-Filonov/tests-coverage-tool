package coverageinupt

type InputCoverageResult struct {
	Method             string   `json:"method"`
	RequestName        string   `json:"requestName"`
	ResponseName       string   `json:"responseName"`
	LogicalService     string   `json:"logicalService"`
	RequestParameters  []string `json:"requestParameters"`
	ResponseParameters []string `json:"responseParameters"`
}
