package constant

var Env string

var SsoAppId string
var SsoAppKey string
var LoginUrl string
var LogoutUrl string
var CheckCodeUrl string
var CheckTicketUrl string
var GetUserNameUrl string

var OldSsoAppId string
var OldSsoAppKey string
var OldLoginUrl string
var OldLogoutUrl string

var RdId string
var AdminId string
var UserInfoPermissUrl string
var UpmPermissUrl string

var HttpTimeOut int64 = 500

var Msg string = "请申请相应角色,获取权限："

const NameSpaceOnlineDebug = "140"
const NameSpacePointSys = "point_sys"
const ToggleIdOnlineDebug = "155765"
const WhiteKeyOnlineDebug = "phone"

const ToggleNotFound = "toggle not found"

const ApolloSettingNameSpace = "point_arch"
const ApolloTypeAB = "ab"
const ApolloTypeGray = "gray"

// WorkflowConfigName Apollo workflow 配置名称
const WorkflowConfigName = "workflow"

// ConfigBlacklist 不允许在自定义配置中使用的 Key
var ConfigBlacklist = map[string]struct{}{
	WorkflowConfigName: {},
}

const ApolloServicePrefix = "map_base_"
