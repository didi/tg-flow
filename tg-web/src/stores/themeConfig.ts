import { defineStore } from 'pinia';
import { ThemeConfigStates, ThemeConfigState } from './interface';

/**
 * 布局配置
 * 修复：https://gitee.com/lyt-top/vue-next-admin/issues/I567R1，感谢@lanbao123
 * 2020.05.28 by lyt 优化。开发时配置不生效问题
 * 修改配置时：
 * 1、需要每次都清理 `window.localStorage` 浏览器永久缓存
 * 2、或者点击布局配置最底部 `一键恢复默认` 按钮即可看到效果
 */
export const useThemeConfig = defineStore('themeConfig', {
	state: (): ThemeConfigStates => ({
		themeConfig: {
			// 是否开启布局配置抽屉
			isDrawer: false,

			/**
			 * 全局主题
			 */
			// 默认 primary 主题颜色
			primary: '#409eff',

			/**
			 * 菜单 / 顶栏
			 * 注意：v1.0.17 版本去除设置布局切换，重置主题样式（initSetLayoutChange），
			 * 切换布局需手动设置样式，设置的样式自动同步各布局，
			 * 代码位置：/@/layout/navBars/breadcrumb/setings.vue
			 */
			// 默认顶栏导航背景颜色
			topBar: '#1d2128',
			// 默认顶栏导航字体颜色
			topBarColor: '#ffffff',
			// 是否开启顶栏背景颜色渐变
			isTopBarColorGradual: false,
			// 默认菜单导航背景颜色
			menuBar: '#ffffff',
			// 默认菜单导航字体颜色
			menuBarColor: '#000000',
			// 是否开启菜单背景颜色渐变
			isMenuBarColorGradual: false,
			// 默认分栏菜单背景颜色
			columnsMenuBar: '#545c64',
			// 默认分栏菜单字体颜色
			columnsMenuBarColor: '#e6e6e6',
			// 是否开启分栏菜单背景颜色渐变
			isColumnsMenuBarColorGradual: false,

			/**
			 * 界面设置
			 */
			// 是否开启菜单水平折叠效果
			isCollapse: false,
			// 是否开启菜单手风琴效果
			isUniqueOpened: false,
			// 是否开启固定 Header
			isFixedHeader: false,
			// 初始化变量，用于更新菜单 el-scrollbar 的高度，请勿删除
			isFixedHeaderChange: false,
			// 是否开启经典布局分割菜单（仅经典布局生效）
			isClassicSplitMenu: false,
			// 是否开启自动锁屏
			isLockScreen: false,
			// 开启自动锁屏倒计时(s/秒)
			lockScreenTime: 30,

			/**
			 * 界面显示
			 */
			// 是否开启侧边栏 Logo
			isShowLogo: true,
			// 初始化变量，用于 el-scrollbar 的高度更新，请勿删除
			isShowLogoChange: false,
			// 是否开启 Breadcrumb，强制经典、横向布局不显示
			isBreadcrumb: true,
			// 是否开启 Tagsview
			isTagsview: false,
			// 是否开启 Breadcrumb 图标
			isBreadcrumbIcon: false,
			// 是否开启 Tagsview 图标
			isTagsviewIcon: false,
			// 是否开启 TagsView 缓存
			isCacheTagsView: false,
			// 是否开启 TagsView 拖拽
			isSortableTagsView: true,
			// 是否开启 TagsView 共用
			isShareTagsView: false,
			// 是否开启 Footer 底部版权信息
			isFooter: false,
			// 是否开启灰色模式
			isGrayscale: false,
			// 是否开启色弱模式
			isInvert: false,
			// 是否开启深色模式
			isIsDark: false,
			// 是否开启水印
			isWartermark: false,
			// 水印文案
			wartermarkText: 'small@小柒',

			/**
			 * 其它设置
			 */
			// Tagsview 风格：可选值"<tags-style-one|tags-style-four|tags-style-five>"，默认 tags-style-five
			// 定义的值与 `/src/layout/navBars/tagsView/tagsView.vue` 中的 class 同名
			tagsStyle: 'tags-style-five',
			// 主页面切换动画：可选值"<slide-right|slide-left|opacitys>"，默认 slide-right
			animation: 'slide-right',
			// 分栏高亮风格：可选值"<columns-round|columns-card>"，默认 columns-round
			columnsAsideStyle: 'columns-round',
			// 分栏布局风格：可选值"<columns-horizontal|columns-vertical>"，默认 columns-horizontal
			columnsAsideLayout: 'columns-vertical',

			/**
			 * 布局切换
			 * 注意：为了演示，切换布局时，颜色会被还原成默认，代码位置：/@/layout/navBars/breadcrumb/setings.vue
			 * 中的 `initSetLayoutChange(设置布局切换，重置主题样式)` 方法
			 */
			// 布局切换：可选值"<defaults|classic|transverse|columns>"，默认 defaults
			layout: 'classic',

			/**
			 * 后端控制路由
			 */
			// 是否开启后端控制路由
			isRequestRoutes: false,

			/**
			 * 全局网站标题 / 副标题
			 */
			// 网站主标题（菜单导航、浏览器当前网页标题）
			globalTitle: import.meta.env.VITE_TITLE,
			// 网站副标题（登录页顶部文字）
			globalViceTitle: import.meta.env.VITE_TITLE,
			// 默认初始语言，可选值"<zh-cn|en|zh-tw>"，默认 zh-cn
			globalI18n: 'zh-cn',
			// 默认全局组件大小，可选值"<large|'default'|small>"，默认 'large'
			globalComponentSize: 'large',
		},
	}),
	actions: {
		setThemeConfig(data: ThemeConfigState) {
			this.themeConfig = data;
		},
	},
});
