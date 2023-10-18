import { h ,withModifiers} from "vue";
import { NTooltip, NIcon, NIconWrapper } from "naive-ui";
import {
  SyncSharp,
  ArrowDownCircleOutline,
  AddCircleOutline,
  RemoveCircleOutline,
  EllipsisHorizontalCircleOutline,
  PauseCircleOutline,
  TimeOutline,
  ReloadCircleOutline,
  NavigateCircleOutline,
} from "@vicons/ionicons5";

export const renderTooltip = (
  text: string,
  click: () => void,
  size?: number
) => {
  return h(NTooltip, null, {
    default: () => text,
    trigger: () => {
      return h(
        "p",
        {
          onClick: withModifiers(click, ['stop']),

          style: {
            color: colorAndType[text].color,
            margin: "0",
            cursor:"pointer"
          },
        },
        { default: () => text }
        // NIconWrapper,
        // {
        //   style: {
        //     margin: "0 5px",
        //   },
        //   color: colorAndType[text].color,
        // },
        // [
        //   h(
        //     NIcon,
        //     {
        //       size: size ? size : 24,
        //       onClick: click,
        //     },
        //     {
        //       default: () => h(colorAndType[text].icon),
        //     }
        //   ),
        // ]
      );
    },
  });
};

const colorAndType = {
  ["修改"]: {
    icon: SyncSharp,
    color: "#2080f0",
  },
  ["提交"]: {
    icon: SyncSharp,
    color: "#2080f0",
  },
  ["导出"]: {
    icon: ArrowDownCircleOutline,
    color: "#2080f0",
  },
  ["导入"]: {
    icon: ArrowDownCircleOutline,
    color:"#2080f0",
  },
  ["复制"]: {
    icon: AddCircleOutline,
    color: "#2080f0",
  },
  ["删除"]: {
    icon: RemoveCircleOutline,
    color: "#d03050",
  },
  ["详情"]: {
    icon: EllipsisHorizontalCircleOutline,
    color: "#2080f0",
  },
  ["更新"]: {
    icon: ReloadCircleOutline,
    color: "#2080f0",
  },
  ["部署"]: {
    icon: PauseCircleOutline,
    color: "#2080f0",
  },
  ["下线"]: {
    icon: ArrowDownCircleOutline,
    color: "#2080f0",
  },
  ["运行历史"]: {
    icon: TimeOutline,
    color: "#2080f0",
  },
  ["例行"]: {
    icon: NavigateCircleOutline,
    color: "#2080f0",
  },
} as any;
// 操作数据配置
