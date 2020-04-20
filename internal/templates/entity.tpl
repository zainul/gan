package entity

// Generate by Thor

{{ .CustomTemplateFromInput }}

{{ .ReqTemplate }}

// Validate {{.Key}} entity...
func (e *{{.Name}}) Validate() error {
	return xsvalidator.Validate(e)
}

func (e *{{.Name}}) TableName() string {
	return "{{.TableName}}"
}