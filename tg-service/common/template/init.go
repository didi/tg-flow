package template

import (
	"github.com/didi/tg-flow/tg-core/common/path"
	"html/template"
	"net/http"
)

var Handler *template.Template

//模型状态编号和中文映射关系
var ModelNameMap map[string]int = make(map[string]int)
var ModelStatusMap map[int]string = make(map[int]string)

//场景编号和中文映射关系
var AppNameMap map[string]int = make(map[string]int)
var AppIdMap map[int]string = make(map[int]string)

//索引状态编号和中文映射关系
var IndexNameMap map[string]int = make(map[string]int)
var IndexIdMap map[int]string = make(map[int]string)

//场景编号和场景名称映射关系
var SceneNameAndIdMap map[string]int = make(map[string]int)
var SceneIdAndNameMap map[int]string = make(map[int]string)

//场景组编号和场景组名称映射关系
var SceneNameAndIdsMap map[string]int = make(map[string]int)
var SceneIdsAndNameMap map[int]string = make(map[int]string)

//流量主从和主从编号映射关系
var DefultNameAndIdMap map[string]int = make(map[string]int)
var DefultIdAndNameMap map[int]string = make(map[int]string)

func init() {
	Handler = template.Must(template.ParseGlob(path.Root + "/template/*.html"))
	ModelNameMap["未知"] = -1
	ModelNameMap["待上线"] = 1
	ModelNameMap["上线失败"] = 2
	ModelNameMap["已上线"] = 3

	ModelStatusMap[-1] = "未知"
	ModelStatusMap[1] = "待上线"
	ModelStatusMap[2] = "上线失败"
	ModelStatusMap[3] = "已上线"

	IndexNameMap["离线"] = 0
	IndexNameMap["在线"] = 1

	IndexIdMap[0] = "离线"
	IndexIdMap[1] = "在线"

	DefultNameAndIdMap["从"] = 0
	DefultNameAndIdMap["主"] = 1

	DefultIdAndNameMap[0] = "从"
	DefultIdAndNameMap[1] = "主"
}

func Render(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return Handler.ExecuteTemplate(w, name, data)
}
