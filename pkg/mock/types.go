package mock

import (
	"errors"
	"fmt"

	"github.com/dhuan/mock/internal/utils"
)

var condition_type_code_encoding_map = map[ConditionType]string{
	ConditionType_None:                  "none",
	ConditionType_HeaderMatch:           "header_match",
	ConditionType_MethodMatch:           "method_match",
	ConditionType_JsonBodyMatch:         "json_body_match",
	ConditionType_FormMatch:             "form_match",
	ConditionType_QuerystringMatch:      "querystring_match",
	ConditionType_QuerystringMatchRegex: "querystring_match_regex",
	ConditionType_QuerystringExactMatch: "querystring_exact_match",
	ConditionType_Nth:                   "nth",
	ConditionType_RouteParamMatch:       "route_param_match",
}

type ConditionType int

const (
	ConditionType_None ConditionType = iota
	ConditionType_HeaderMatch
	ConditionType_MethodMatch
	ConditionType_JsonBodyMatch
	ConditionType_FormMatch
	ConditionType_QuerystringMatch
	ConditionType_QuerystringMatchRegex
	ConditionType_QuerystringExactMatch
	ConditionType_Nth
	ConditionType_RouteParamMatch
)

func (this *ConditionType) UnmarshalJSON(data []byte) (err error) {
	conditionTypeText := utils.Unquote(string(data))

	if conditionTypeText == "header_match" {
		*this = ConditionType_HeaderMatch

		return nil
	}

	if conditionTypeText == "method_match" {
		*this = ConditionType_MethodMatch

		return nil
	}

	if conditionTypeText == "json_body_match" {
		*this = ConditionType_JsonBodyMatch

		return nil
	}

	if conditionTypeText == "form_match" {
		*this = ConditionType_FormMatch

		return nil
	}

	if conditionTypeText == "querystring_match" {
		*this = ConditionType_QuerystringMatch

		return nil
	}

	if conditionTypeText == "querystring_match_regex" {
		*this = ConditionType_QuerystringMatchRegex

		return nil
	}

	if conditionTypeText == "querystring_exact_match" {
		*this = ConditionType_QuerystringExactMatch

		return nil
	}

	if conditionTypeText == "nth" {
		*this = ConditionType_Nth

		return nil
	}

	if conditionTypeText == "route_param_match" {
		*this = ConditionType_RouteParamMatch

		return nil
	}

	return errors.New(fmt.Sprintf("Failed to parse Condition Type: %s", conditionTypeText))
}

func (this *ConditionType) MarshalJSON() ([]byte, error) {
	encodingMapPrepared := utils.MapMapValueOnly(
		condition_type_code_encoding_map,
		utils.WrapIn(`"`),
	)

	return utils.MarshalJsonHelper(
		encodingMapPrepared,
		"Failed to parse Condition Type Code: %d",
		this,
	)
}

type ConditionValue string

func (this *ConditionValue) UnmarshalJSON(data []byte) (err error) {
	conditionTypeText := utils.Unquote(string(data))
	var valueStringified ConditionValue = ConditionValue(conditionTypeText)

	*this = valueStringified

	return nil
}

type Condition struct {
	Type      ConditionType          `json:"type"`
	Key       string                 `json:"key"`
	Value     ConditionValue         `json:"value"`
	KeyValues map[string]interface{} `json:"key_values"`
	And       *Condition             `json:"and"`
	Or        *Condition             `json:"or"`
}
