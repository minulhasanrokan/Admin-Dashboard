package requests

type GetDropdownRequest struct {
	TableName     string `param:"table_name" json:"table_name" form:"table_name" query:"table_name" validate:"required"`
	FieldName     string `param:"field_name" json:"field_name" form:"field_name" query:"field_name" validate:"required"`
	Value         string `param:"value" json:"value" form:"value" query:"value"`
	OtherField    string `param:"other_field" json:"other_field" form:"other_field" query:"other_field"`
	OtherValue    string `param:"other_value" json:"other_value" form:"other_value" query:"other_value"`
	NotEqualField string `param:"not_equal_field" json:"not_equal_field" form:"not_equal_field" query:"not_equal_field"`
	NotEqualValue string `param:"not_equal_value" json:"not_equal_value" form:"not_equal_value" query:"not_equal_value"`
}

type UUIDRequest struct {
	UUID string `param:"uuid" json:"uuid" form:"uuid" query:"uuid" validate:"required,uuid4"`
}
