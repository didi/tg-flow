/**
Description : condition test
Author		: dayunzhangyunfeng@didiglobal.com
Date		: 2021-06-14
*/

package test

import (
	"fmt"
)

type ConditionExternal struct {
}

func (c *ConditionExternal) In(key string, collection map[string]interface{}) (string, error){
	if 	collection == nil {
		return "N", fmt.Errorf("collection must no be nil")
	}

	if _, ok := collection[key];ok {
		return "Y", nil
	}

	return "N", nil
}