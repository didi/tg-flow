import { sortByKey } from '../utils/arrayOperation';
// 只处理到二级菜单
export function formatMenus(menusUpm: Array<any>) {
	menusUpm = sortByKey(menusUpm, 'sortVal');
	const rootMenu: any[] = [];
	menusUpm.map((item: any) => {
		if (item.pid === 0) {
			rootMenu.push(makeMenu(item));
		}
	})
	for (let i = 0; i < rootMenu.length; i++) {
		const rootMenuItem = rootMenu[i];
		const children: any[] = [];
		for (let i = 0; i < menusUpm.length; i++) {
			const subMenu = menusUpm[i];
			const isMenu = subMenu.isMenu;
			if (!isMenu || subMenu.pid !== rootMenuItem.id) {
				continue;
			}
			children.push(makeMenu(subMenu));
		}
		rootMenuItem.children = sortByKey(children, 'sortVal'); 
	}
	return rootMenu;
}

function makeMenu(menuUpm: any) {
	const childMenu = {
		id: menuUpm.id,
		path: menuUpm.url,
		name: menuUpm.name,
		component: menuUpm.featureKey,
		meta: {
			title: menuUpm.name,
			isLink: '',
			isHide: false,
			isKeepAlive: true,
			isAffix: true,
			isIframe: false,
			roles: ['common'],
			icon: menuUpm.icon,
		},
	}
	return childMenu;
}