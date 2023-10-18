/**
	Description : loader of scene and module config info
	Author		: dayunzhangyunfeng@didiglobal.com
	Date		: 2021-05-14
*/

package wfengine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/redis"
	"io/ioutil"
	"os"
)

const (
	//新版redis key
	RedisKeySceneModule	= "scene_module_app_"
)
func LoadSceneModuleMap(appId int64) (map[int64]*SceneModule, error) {
	//加载redis中的数据到内存
	sceneModuleMap, err := loadSceneModules(appId)
	if err != nil {
		return sceneModuleMap, err
	}

	//根据app_id,获取对应的场景id
	for sceneId, sceneModule := range sceneModuleMap {
		if sceneModule.AppId != appId {
			delete(sceneModuleMap, sceneId)
		}
	}

	return sceneModuleMap, nil
}

//获取redis中的数据
func loadSceneModules(appId int64) (map[int64]*SceneModule, error) {
	sceneModuleMapString, err := redis.Handler.Get(context.TODO(), fmt.Sprintf("%v%v", RedisKeySceneModule, appId))
	if err != nil && err.Error() != redis.ErrNil {
		return nil, err
	}

	var sceneModuleMap map[int64]*SceneModule
	err = json.Unmarshal([]byte(sceneModuleMapString), &sceneModuleMap)
	if err != nil {
		return nil, fmt.Errorf("err:%v, sceneModule:%v", err, sceneModuleMapString)
	}

	return sceneModuleMap, nil
}

func LoadSceneModuleMapFromFile(path string) (map[int64]*SceneModule, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.open error, file:%v, err:%v", path, err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll error, byteValue:%v,err=%v", string(byteValue), err)
	}

	var sceneModuleMap map[int64]*SceneModule
	err = json.Unmarshal(byteValue, &sceneModuleMap)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error, byteValue:%v,err=%v", string(byteValue), err)
	}

	return sceneModuleMap, nil
}