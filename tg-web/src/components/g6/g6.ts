import G6 from '@antv/g6';
import insertCss from 'insert-css';

insertCss(`
  .g6-component-toolbar li {
    list-style-type: none !important;
  }
  .g6-minimap-container {
    border: 1px solid #e2e2e2;
  }
  .g6-minimap-viewport {
    border: 2px solid rgb(25, 128, 255);
  }
`);

const ICON_MAP = {
  a: 'https://gw.alipayobjects.com/mdn/rms_8fd2eb/afts/img/A*0HC-SawWYUoAAAAAAAAAAABkARQnAQ',
  b: 'https://gw.alipayobjects.com/mdn/rms_8fd2eb/afts/img/A*sxK0RJ1UhNkAAAAAAAAAAABkARQnAQ',
};
const color = '#5B8FF9';

function drawRect(cfg: any, group: any) {
  // console.log(cfg)
  const r = 2;
  const width = 150
  const shape = group.addShape('rect', {
    attrs: {
      x: 0,
      y: 0,
      width,
      height: 60,
      stroke: color,
      radius: r,
      cursor: 'move',
      anchorPoints: [
        [0, 0.5],
        [1, 0.5],
      ],
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'main-box',
    draggable: true,
  });

  group.addShape('rect', {
    attrs: {
      x: 0,
      y: 0,
      width,
      height: 30,
      fill: color,
      radius: [r, r, 0, 0],
      cursor: 'move',
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'title-box',
    draggable: true,
  });

  // left icon
  group.addShape('image', {
    attrs: {
      x: 4,
      y: 7,
      height: 16,
      width: 16,
      cursor: 'pointer',
      img: ICON_MAP['a'],
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'node-icon',
  });

  // title text
  group.addShape('text', {
    attrs: {
      textBaseline: 'top',
      y: 9,
      x: 24,
      lineHeight: 20,
      fontSize: 12,
      fontWeight: 500,
      text: 'ID：' + cfg.id,
      fill: '#fff',
      cursor: 'pointer',
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'title',
  });

  // The content list
  // name text
  group.addShape('text', {
    attrs: {
      textBaseline: 'top',
      y: 40,
      x: 10,
      lineHeight: 20,
      fontSize: 12,
      text: fittingString(cfg.label, 130, 12),
      fill: 'rgba(0,0,0, 0.4)',
      cursor: 'pointer',
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: `desc`,
  });
  return shape
}

function drawTimeout(cfg: any, group: any) {
  const r = 2;
  const width = 150
  const shape = group.addShape('rect', {
    attrs: {
      x: 0,
      y: 0,
      width,
      height: 60,
      stroke: color,
      radius: r,
      anchorPoints: [
        [0, 0.5],
        [1, 0.5],
      ],
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'main-box',
    draggable: true,
  });

  group.addShape('rect', {
    attrs: {
      x: 0,
      y: 0,
      width,
      height: 30,
      fill: color,
      radius: [r, r, 0, 0],
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'title-box',
    draggable: true,
  });

  // left icon
  group.addShape('image', {
    attrs: {
      x: 4,
      y: 7,
      height: 16,
      width: 16,
      cursor: 'pointer',
      img: ICON_MAP['a'],
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'node-icon',
  });

  // title text
  group.addShape('text', {
    attrs: {
      textBaseline: 'top',
      y: 9,
      x: 24,
      lineHeight: 20,
      fontSize: 12,
      fontWeight: 500,
      text: 'ID：' + cfg.id,
      fill: '#fff',
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'title',
  });

  // The content list
  // name text
  group.addShape('text', {
    attrs: {
      textBaseline: 'top',
      y: 40,
      x: 10,
      lineHeight: 20,
      fontSize: 12,
      text: fittingString(cfg.label, 130, 12),
      fill: '#666',
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: `desc`,
  });
  return shape
}

function drawCondition(cfg: any, group: any) {
  const width = 80
  const shape = group.addShape('polygon', {
    attrs: {
      points: [
        [40, 0],
        [0, 30],
        [40, 60],
        [80, 30],
      ],
      stroke: color,
      fill: color,
      cursor: 'move',
      anchorPoints: [
        [0, 0.5],
        [1, 0.5],
      ],
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'main-box',
    draggable: true,
  });
  group.addShape('text', {
    attrs: {
      textBaseline: 'top',
      y: 25,
      x: 32,
      lineHeight: 20,
      fontSize: 12,
      fontWeight: 500,
      text: cfg.label,
      fill: '#fff',
      cursor: 'pointer',
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'title',
  });
  return shape
}

function drawFlow(cfg: any, group: any) {
  const shape = group.addShape('circle', {
    attrs: {
      x: 0,
      y: 0,
      r: 50,
      stroke: color,
      cursor: 'move',
      fill: color,
      anchorPoints: [
        [0, 0.5],
        [1, 0.5],
      ],
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'main-box',
    draggable: true,
  });
  group.addShape('text', {
    attrs: {
      textBaseline: 'middle',
      lineHeight: 20,
      fontSize: 12,
      fontWeight: 500,
      textAlign: 'center',
      text: cfg.label,
      fill: '#fff',
      cursor: 'pointer',
    },
    // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
    name: 'title',
  });
  return shape
}

G6.registerNode(
  'timeout',
  {
    draw: function drawShape(cfg: any, group: any) {
      return drawTimeout(cfg, group)
    },
    afterDraw(cfg, group) {
      const bbox = group.getBBox();
      const anchorPoints = this.getAnchorPoints(cfg)
      anchorPoints.forEach((anchorPos, i) => {
        group.addShape('circle', {
          attrs: {
            r: 5,
            x: bbox.x + bbox.width * anchorPos[0],
            y: bbox.y + bbox.height * anchorPos[1],
            fill: '#fff',
            stroke: '#5F95FF'
          },
          // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
          name: `anchor-point`, // the name, for searching by group.find(ele => ele.get('name') === 'anchor-point')
          anchorPointIdx: i, // flag the idx of the anchor-point circle
          links: 0, // cache the number of edges connected to this shape
          visible: false, // invisible by default, shows up when links > 1 or the node is in showAnchors state
          draggable: true // allow to catch the drag events on this shape
        })
      })
    },
    getAnchorPoints(cfg) {
      return cfg.anchorPoints || [[0, 0.5], [0.33, 0], [0.66, 0], [1, 0.5], [0.33, 1], [0.66, 1]];
    },
    // update: function drawShape(cfg: any, item: any) {
    //   console.log('---',cfg, item)
    //   item.label = cfg.label
    // }
    setState(name, value, item) {
      const group = item!.getContainer();
      if (name === 'selected') {
        const box = group.find((element) => element.get('name') === 'main-box');
        const lineWidth = value ? 2 : 1
        box.attr('lineWidth', lineWidth)
      }
      if (name === 'showAnchors') {
        const anchorPoints = item.getContainer().findAll(ele => ele.get('name') === 'anchor-point');
        anchorPoints.forEach(point => {
          if (value || point.get('links') > 0) point.show()
          else point.hide()
        })
      }
    }
  },
);

G6.registerNode(
  'task',
  {
    draw: function drawShape(cfg: any, group: any) {
      return drawRect(cfg, group)
    },
    afterDraw(cfg, group) {
      const bbox = group.getBBox();
      const anchorPoints = this.getAnchorPoints(cfg)
      anchorPoints.forEach((anchorPos, i) => {
        group.addShape('circle', {
          attrs: {
            r: 5,
            x: bbox.x + bbox.width * anchorPos[0],
            y: bbox.y + bbox.height * anchorPos[1],
            fill: '#fff',
            stroke: '#5F95FF'
          },
          // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
          name: `anchor-point`, // the name, for searching by group.find(ele => ele.get('name') === 'anchor-point')
          anchorPointIdx: i, // flag the idx of the anchor-point circle
          links: 0, // cache the number of edges connected to this shape
          visible: false, // invisible by default, shows up when links > 1 or the node is in showAnchors state
          draggable: true // allow to catch the drag events on this shape
        })
      })
    },
    getAnchorPoints(cfg) {
      return cfg.anchorPoints || [[0, 0.5], [0.33, 0], [0.66, 0], [1, 0.5], [0.33, 1], [0.66, 1]];
    },
    // update: function drawShape(cfg: any, item: any) {
    //   console.log('---',cfg, item)
    //   item.label = cfg.label
    // }
    setState(name, value, item) {
      const group = item!.getContainer();
      if (name === 'selected') {
        const box = group.find((element) => element.get('name') === 'main-box');
        const lineWidth = value ? 2 : 1
        box.attr('lineWidth', lineWidth)
      }
      if (name === 'showAnchors') {
        const anchorPoints = item.getContainer().findAll(ele => ele.get('name') === 'anchor-point');
        anchorPoints.forEach(point => {
          if (value || point.get('links') > 0) point.show()
          else point.hide()
        })
      }
    }
  },
);

G6.registerNode(
  'condition',
  {
    draw: function drawShape(cfg: any, group: any) {
      return drawCondition(cfg, group)
    },
    afterDraw(cfg, group) {
      const bbox = group.getBBox();
      const anchorPoints = this.getAnchorPoints(cfg)
      anchorPoints.forEach((anchorPos, i) => {
        group.addShape('circle', {
          attrs: {
            r: 5,
            x: bbox.x + bbox.width * anchorPos[0],
            y: bbox.y + bbox.height * anchorPos[1],
            fill: '#fff',
            stroke: '#5F95FF'
          },
          // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
          name: `anchor-point`, // the name, for searching by group.find(ele => ele.get('name') === 'anchor-point')
          anchorPointIdx: i, // flag the idx of the anchor-point circle
          links: 0, // cache the number of edges connected to this shape
          visible: false, // invisible by default, shows up when links > 1 or the node is in showAnchors state
          draggable: true // allow to catch the drag events on this shape
        })
      })
    },
    getAnchorPoints(cfg) {
      return cfg.anchorPoints || [[0, 0.5], [0.33, 0], [0.66, 0], [1, 0.5], [0.33, 1], [0.66, 1]];
    },
    // update: function drawShape(cfg: any, item: any) {
    //   console.log('---',cfg, item)
    //   item.label = cfg.label
    // }
    setState(name, value, item) {
      const group = item!.getContainer();
      if (name === 'selected') {
        const box = group.find((element) => element.get('name') === 'main-box');
        const lineWidth = value ? 2 : 1
        box.attr('lineWidth', lineWidth)
      }
      if (name === 'showAnchors') {
        const anchorPoints = item.getContainer().findAll(ele => ele.get('name') === 'anchor-point');
        anchorPoints.forEach(point => {
          if (value || point.get('links') > 0) point.show()
          else point.hide()
        })
      }
    }
  },
);

G6.registerNode(
  'flow',
  {
    draw: function drawShape(cfg: any, group: any) {
      return drawFlow(cfg, group)
    },
    afterDraw(cfg, group) {
      const bbox = group.getBBox();
      const anchorPoints = this.getAnchorPoints(cfg)
      anchorPoints.forEach((anchorPos, i) => {
        group.addShape('circle', {
          attrs: {
            r: 5,
            x: bbox.x + bbox.width * anchorPos[0],
            y: bbox.y + bbox.height * anchorPos[1],
            fill: '#fff',
            stroke: '#5F95FF'
          },
          // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
          name: `anchor-point`, // the name, for searching by group.find(ele => ele.get('name') === 'anchor-point')
          anchorPointIdx: i, // flag the idx of the anchor-point circle
          links: 0, // cache the number of edges connected to this shape
          visible: false, // invisible by default, shows up when links > 1 or the node is in showAnchors state
          draggable: true // allow to catch the drag events on this shape
        })
      })
    },
    getAnchorPoints(cfg) {
      return cfg.anchorPoints || [[0, 0.5], [0.33, 0], [0.66, 0], [1, 0.5], [0.33, 1], [0.66, 1]];
    },
    // update: function drawShape(cfg: any, item: any) {
    //   console.log('---',cfg, item)
    //   item.label = cfg.label
    // }
    setState(name, value, item) {
      const group = item!.getContainer();
      if (name === 'selected') {
        const box = group.find((element) => element.get('name') === 'main-box');
        const lineWidth = value ? 2 : 1
        box.attr('lineWidth', lineWidth)
      }
      if (name === 'showAnchors') {
        const anchorPoints = item.getContainer().findAll(ele => ele.get('name') === 'anchor-point');
        anchorPoints.forEach(point => {
          if (value || point.get('links') > 0) point.show()
          else point.hide()
        })
      }
    }
  },
);

G6.registerNode(
  'card-node',
  {
    draw: function drawShape(cfg: any, group: any) {
      return drawRect(cfg, group)
    },
    afterDraw(cfg, group) {
      const bbox = group.getBBox();
      const anchorPoints = this.getAnchorPoints(cfg)
      anchorPoints.forEach((anchorPos, i) => {
        group.addShape('circle', {
          attrs: {
            r: 5,
            x: bbox.x + bbox.width * anchorPos[0],
            y: bbox.y + bbox.height * anchorPos[1],
            fill: '#fff',
            stroke: '#5F95FF'
          },
          // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
          name: `anchor-point`, // the name, for searching by group.find(ele => ele.get('name') === 'anchor-point')
          anchorPointIdx: i, // flag the idx of the anchor-point circle
          links: 0, // cache the number of edges connected to this shape
          visible: false, // invisible by default, shows up when links > 1 or the node is in showAnchors state
          draggable: true // allow to catch the drag events on this shape
        })
      })
    },
    getAnchorPoints(cfg) {
      return cfg.anchorPoints || [[0, 0.5], [0.33, 0], [0.66, 0], [1, 0.5], [0.33, 1], [0.66, 1]];
    },
    // update: function drawShape(cfg: any, item: any) {
    //   console.log('---',cfg, item)
    //   item.label = cfg.label
    // }
    setState(name, value, item) {
      const group = item!.getContainer();
      if (name === 'selected') {
        const box = group.find((element) => element.get('name') === 'main-box');
        const lineWidth = value ? 2 : 1
        box.attr('lineWidth', lineWidth)
      }
      if (name === 'showAnchors') {
        const anchorPoints = item.getContainer().findAll(ele => ele.get('name') === 'anchor-point');
        anchorPoints.forEach(point => {
          if (value || point.get('links') > 0) point.show()
          else point.hide()
        })
      }
    }
  },
);

const collapseIcon = (x, y, r) => {
  return [
    ['M', x - r, y],
    ['a', r, r, 0, 1, 0, r * 2, 0],
    ['a', r, r, 0, 1, 0, -r * 2, 0],
    ['M', x - r + 4, y],
    ['L', x - r + 2 * r - 4, y],
  ];
};

const expandIcon = (x, y, r) => {
  return [
    ['M', x - r, y],
    ['a', r, r, 0, 1, 0, r * 2, 0],
    ['a', r, r, 0, 1, 0, -r * 2, 0],
    ['M', x - r + 4, y],
    ['L', x - r + 2 * r - 4, y],
    ['M', x - r + r, y - r + 4],
    ['L', x, y + r - 4],
  ];
};

G6.registerCombo(
  'cCircle',
  {
    drawShape: function draw(cfg, group) {
      const self = this;
      // Get the shape style, where the style.r corresponds to the R in the Illustration of Built-in Rect Combo
      const style = self.getShapeStyle(cfg);
      const circle = group.addShape('circle', {
        attrs: {
          ...style,
          x: 0,
          y: 0,
          // r: 100,
          r: style.r,
          fill: '#fff',
          // fill: '#000',
        },
        draggable: true,
        // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
        name: 'combo-keyShape',
      });
      // Add the marker on the bottom
      const marker = group.addShape('marker', {
        attrs: {
          ...style,
          fill: '#fff',
          // fill: '#000',
          opacity: 1,
          x: 0,
          // y: 100,
          y: style.r,
          r: 10,
          symbol: collapseIcon,
        },
        draggable: true,
        // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
        name: 'combo-marker-shape',
      });

      return circle;
    },
    // Define the updating logic for the marker
    afterUpdate: function afterUpdate(cfg, combo) {
      const self = this;
      // Get the shape style, where the style.r corresponds to the R in the Illustration of Built-in Rect Combo
      const style = self.getShapeStyle(cfg);
      const group = combo.get('group');
      // Find the marker shape in the graphics group of the Combo
      const marker = group.find((ele) => ele.get('name') === 'combo-marker-shape');
      // Update the marker shape
      marker.attr({
        x: 0,
        y: style.r,
        // The property 'collapsed' in the combo data represents the collapsing state of the Combo
        // Update the symbol according to 'collapsed'
        symbol: cfg.collapsed ? expandIcon : collapseIcon,
      });
    },
  },
  'circle',
);

G6.registerCombo(
  'cRect',
  {
    drawShape: function draw(cfg, group) {
      const self = this;
      cfg.padding = cfg.padding || [3, 3, 3, 3];
      // Get the shape style, where the style.r corresponds to the R in the Illustration of Built-in Rect Combo
      const style = self.getShapeStyle(cfg);
      const rect = group.addShape('rect', {
        attrs: {
          ...style,
          x:  - (cfg.padding[3] - cfg.padding[1]) / 2,
          y:  - (cfg.padding[0] - cfg.padding[2]) / 2,
          width: style.width,
          height: style.height,
          fill: '#fff',
          // fill: '#000',
        },
        draggable: true,
        // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
        name: 'combo-keyShape',
      });
      // Add the marker on the bottom
      const marker = group.addShape('marker', {
        attrs: {
          ...style,
          fill: '#fff',
          // fill: '#000',
          opacity: 1,
          x: cfg.style.width / 2 + cfg.padding[1],
          y: (cfg.padding[2] - cfg.padding[0]) / 2,
          r: 10,
          symbol: collapseIcon,
        },
        draggable: true,
        // must be assigned in G6 3.3 and later versions. it can be any string you want, but should be unique in a custom item type
        name: 'combo-marker-shape',
      });

      return rect;
    },
    // Define the updating logic for the marker
    afterUpdate: function afterUpdate(cfg, combo) {
      const group = combo.get('group');
      // Find the circle shape in the graphics group of the Combo by name
      const marker = group.find((ele) => ele.get('name') === 'combo-marker-shape');
      // Update the position of the right circle
      marker.attr({
        // cfg.style.width and cfg.style.heigth correspond to the innerWidth and innerHeight in the figure of Illustration of Built-in Rect Combo
        x: cfg.style.width / 2 + cfg.padding[1],
        y: (cfg.padding[2] - cfg.padding[0]) / 2,
        // The property 'collapsed' in the combo data represents the collapsing state of the Combo
        // Update the symbol according to 'collapsed'
        symbol: cfg.collapsed ? expandIcon : collapseIcon,
      });
    },
  },
  'rect',
);

const fittingString = (str: string, maxWidth: number, fontSize: number) => {
  if (!str) {
    return ''
  }
  const ellipsis = "...";
  const ellipsisLength = G6.Util.getTextSize(ellipsis, fontSize)[0];
  let currentWidth = 0;
  let res = str;
  const pattern = new RegExp("[\u4E00-\u9FA5]+"); // distinguish the Chinese charactors and letters
  str.split("").forEach((letter, i) => {
    if (currentWidth > maxWidth - ellipsisLength) return;
    if (pattern.test(letter)) {
      // Chinese charactors
      currentWidth += fontSize;
    } else {
      // get the width of single letter according to the fontSize
      currentWidth += G6.Util.getLetterWidth(letter, fontSize);
    }
    if (currentWidth > maxWidth - ellipsisLength) {
      res = `${str.substr(0, i)}${ellipsis}`;
    }
  });
  return res;
};

const processParallelEdgesOnAnchorPoint = (
  edges,
  offsetDiff = 15,
  multiEdgeType = 'cubic-horizontal',
  singleEdgeType = undefined,
  loopEdgeType = undefined
) => {
  const len = edges.length;
  const cod = offsetDiff * 2;
  const loopPosition = [
    'top',
    'top-right',
    'right',
    'bottom-right',
    'bottom',
    'bottom-left',
    'left',
    'top-left',
  ];
  const edgeMap = {};
  const tags = [];
  const reverses = {};
  for (let i = 0; i < len; i++) {
    const edge = edges[i];
    const { source, target, sourceAnchor, targetAnchor } = edge;
    const sourceTarget = `${source}|${sourceAnchor}-${target}|${targetAnchor}`;

    if (tags[i]) continue;
    if (!edgeMap[sourceTarget]) {
      edgeMap[sourceTarget] = [];
    }
    tags[i] = true;
    edgeMap[sourceTarget].push(edge);
    for (let j = 0; j < len; j++) {
      if (i === j) continue;
      const sedge = edges[j];
      const { source: src, target: dst, sourceAnchor: srcAnchor, targetAnchor: dstAnchor } = sedge;

      // 两个节点之间共同的边
      // 第一条的source = 第二条的target
      // 第一条的target = 第二条的source
      if (!tags[j]) {
        if (source === dst && sourceAnchor === dstAnchor
            && target === src && targetAnchor === srcAnchor) {
          edgeMap[sourceTarget].push(sedge);
          tags[j] = true;
          reverses[`${src}|${srcAnchor}|${dst}|${dstAnchor}|${edgeMap[sourceTarget].length - 1}`] = true;
        } else if (source === src && sourceAnchor === srcAnchor
           && target === dst  && targetAnchor === dstAnchor) {
          edgeMap[sourceTarget].push(sedge);
          tags[j] = true;
        }
      }
    }
  }

  for (const key in edgeMap) {
    const arcEdges = edgeMap[key];
    const { length } = arcEdges;
    for (let k = 0; k < length; k++) {
      const current = arcEdges[k];
      if (current.source === current.target) {
        if (loopEdgeType) current.type = loopEdgeType;
        // 超过8条自环边，则需要重新处理
        current.loopCfg = {
          position: loopPosition[k % 8],
          dist: Math.floor(k / 8) * 20 + 50,
        };
        continue;
      }
      if (length === 1 && singleEdgeType && (current.source !== current.target || current.sourceAnchor !== current.targetAnchor)) {
        current.type = singleEdgeType;
        continue;
      }
      current.type = multiEdgeType;
      const sign =
        (k % 2 === 0 ? 1 : -1) * (reverses[`${current.source}|${current.sourceAnchor}|${current.target}|${current.targetAnchor}|${k}`] ? -1 : 1);
      if (length % 2 === 1) {
        current.curveOffset = sign * Math.ceil(k / 2) * cod;
      } else {
        current.curveOffset = sign * (Math.floor(k / 2) * cod + offsetDiff);
      }
    }
  }
  return edges;
};

export function renderFlow(containerId: string, nodes: any, edges: any, combos: any, events: any, editabel: boolean) {
  const container = document.getElementById(containerId);
  const width = container!.scrollWidth;
  const height = container!.scrollHeight;
  const toolbar = new G6.ToolBar({
    container: 'g6-toolbar',
    position: { x: 10, y: 10 },
  });
  const minimap = new G6.Minimap({
    container: 'g6-minimap',
    size: [150, 100]
  });
  let sourceAnchorIdx, targetAnchorIdx
  const graph = new G6.Graph({
    container: containerId,
    width: width - 15,
    height,
    animate: true,
    fitView: false,
    fitCenter: false,
    groupByTypes: false,
    plugins: [toolbar, minimap],
    // 设置为true，启用 redo & undo 栈功能
    enabledStack: true,
    modes: {
      default: [
        'drag-canvas', 'zoom-canvas', 'drag-combo',
        {
          type: 'drag-node',
          shouldBegin: e => {
            if (e.target.get('name') === 'anchor-point') return false;
            if(e.item?._cfg?.model?.hasOwnProperty('comboId') && e.item?._cfg?.model?.comboId !== undefined) return false;
            return true;
          }
        },
        {
          type: 'create-edge',
          trigger: 'drag', // set the trigger to be drag to make the create-edge triggered by drag
          shouldBegin: e => {
            // avoid beginning at other shapes on the node
            if (e.target && e.target.get('name') !== 'anchor-point') return false;
            sourceAnchorIdx = e.target.get('anchorPointIdx');
            e.target.set('links', e.target.get('links') + 1); // cache the number of edge connected to this anchor-point circle
            return true;
          },
          shouldEnd: e => {
            // avoid ending at other shapes on the node
            if (e.target && e.target.get('name') !== 'anchor-point') return false;
            if (e.target) {
              targetAnchorIdx = e.target.get('anchorPointIdx');
              e.target.set('links', e.target.get('links') + 1);  // cache the number of edge connected to this anchor-point circle
              return true;
            }
            targetAnchorIdx = undefined;
            return true;
          },
        },
      ],
    },
    // layout: {
    //   // type: 'dagre',
    //   // rankdir: 'LR',
    //   // align: 'UL',
    //   // controlPoints: false,
    //   // nodesepFunc: () => 20,
    //   // ranksepFunc: () => 50,
    // },
    defaultNode: {
      type: 'card-node',
      anchorPoints: [[0, 0.5], [1, 0.5]]
    },
    defaultEdge: {
      type: 'cubic-horizontal',
      size: 1,
      style: {
        stroke: '#AAB7C4',
        endArrow: {
          // path: G6.Arrow.circle(3, 2),
          path: G6.Arrow.vee(6, 10, 2),
          d: 2,
          fill: '#AAB7C4',
        },
        radius: 20,
      },
    },
    defaultCombo: {
      type: 'cRect',
      // type: 'cCircle',
      style:{
        lineWidth: 3,
        stroke: color,
      },
      fixCollapseSize: 5,
      labelCfg: {
        /* label's offset to the keyShape */
        refY: -15,
        refX: -5,
        // position: 'top',
        style: {
          fontSize: 12,
        },
      },
    },
    nodeStateStyles: {
      hover: {
        lineWidth: 5,
        stroke: '#000',
        fill: '#000'
      },
    },
  });
  const newNodes: any = []
  nodes.map((node: any) => {
    const n: any = {}
    n.id = node.id
    n.label = node.label
    n.data = node
    n.type = node.type
    n.x = node.x
    n.y = node.y
    n.comboId = node.comboId
    newNodes.push(n)
  })
  const newEdeges: any = []
  edges.map((edge: any) => {
    newEdeges.push({
      source: edge.source,
      target: edge.target,
      data: edge,
      comboId: edge.comboId,
    })
  })
  const newCombos: any = []
  combos.map((combo: any) => {
    newCombos.push({
      id: combo.id,
      label: combo.label,
      collapsed: combo.collapsed,
    })
  })
  // console.log(JSON.stringify(newNodes))
  graph.data({ nodes: newNodes, edges: newEdeges, combos: newCombos});
  graph.node((node) => {
    const data: any = node.data
    // if (data.type === 'condition') {
    if (data.type === 'task') {
      node.type = 'task'
    }
    if (data.type === 'timeout') {
      node.type = 'timeout'
    }
    if (data.type === 'condition') {
      node.type = 'condition'
    }
    if (data.type === 'flow') {
      node.type = 'flow'
    }
    return node
  });
  graph.render();

  // const temp_nodes = graph.getNodes();
  // for (let temp_k in temp_nodes) {
  //   const temp_node = temp_nodes[temp_k]
  //   console.log(temp_node)
  //   // afterDrawCircle(temp_node._cfg, temp_node._cfg?.group)
  // }


  graph.moveTo(0, graph.getHeight() / 2 - 75)

  graph.on('click', () => {
  });
  graph.on('node:click', (evt: any) => {
    const { item } = evt;
    console.log(item)
    graph.getNodes().forEach((node) => {
      graph.clearItemStates(node);
    });
    setTimeout(() => graph.setItemState(item, 'selected', true));
    if (editabel){
      events['click'](item)
    }

  });
  graph.on('edge:mouseenter', (evt: any) => {
    const { item } = evt;
    // console.log(item)
    graph.getEdges().forEach((edge) => {
      graph.clearItemStates(edge);
    });
    setTimeout(() => graph.setItemState(item, 'selected', true));
    // events['clickEdge'](item)
  });

  graph.on('edge:mouseleave', (evt: any) => {
    const { item } = evt;
    // console.log(item)
    graph.getEdges().forEach((edge) => {
      graph.clearItemStates(edge);
    });
    setTimeout(() => graph.setItemState(item, 'selected', false));
    // events['clickEdge'](item)
  });

  // after drag from the first node, the edge is created, update the sourceAnchor
  graph.on('afteradditem', e => {
    if (e.item && e.item.getType() === 'edge') {
      graph.updateItem(e.item, {
        sourceAnchor: sourceAnchorIdx
      });
    }
  })

  // if create-edge is canceled before ending, update the 'links' on the anchor-point circles
  graph.on('afterremoveitem', e => {
    if (e.item && e.item.source && e.item.target) {
      const sourceNode = graph.findById(e.item.source);
      const targetNode = graph.findById(e.item.target);
      const { sourceAnchor, targetAnchor } = e.item;
      if (sourceNode && !isNaN(sourceAnchor)) {
        const sourceAnchorShape = sourceNode.getContainer().find(ele => (ele.get('name') === 'anchor-point' && ele.get('anchorPointIdx') === sourceAnchor));
        sourceAnchorShape.set('links', sourceAnchorShape.get('links') - 1);
      }
      if (targetNode && !isNaN(targetAnchor)) {
        const targetAnchorShape = targetNode.getContainer().find(ele => (ele.get('name') === 'anchor-point' && ele.get('anchorPointIdx') === targetAnchor));
        targetAnchorShape.set('links', targetAnchorShape.get('links') - 1);
      }
    }
  })

  graph.on('aftercreateedge', (e) => {
    graph.updateItem(e.edge, {
      sourceAnchor: sourceAnchorIdx,
      targetAnchor: targetAnchorIdx
    })

    // update the curveOffset for parallel edges
    const edges = graph.save().edges;
    processParallelEdgesOnAnchorPoint(edges);
    graph.getEdges().forEach((edge, i) => {
      graph.updateItem(edge, {
        curveOffset: edges[i].curveOffset,
        curvePosition: edges[i].curvePosition,
      });
    });
  });

  const dataChange = () => {
    // const nodes: any = [], edges: any = []
    // graph.getNodes().map(node => {
    //   nodes.push(node.getModel())
    // })
    // graph.getEdges().map(edge => {
    //   edges.push(edge.getModel())
    // })
    // const data = { nodes, edges }
    const data: any = graph.save()
    const nodes: any = []
    const edges: any = []
    data.nodes.forEach((node: any) => {
      nodes.push(Object.assign({}, {...node}))
    })
    data.edges.forEach((edge: any) => {
      edges.push(Object.assign({}, {...edge}))
    })
    
    const newData = {nodes, edges}
    if (editabel){
      events['datachange'](newData)
    }
    // console.log('kkkkkkk',newData)
  }

  graph.on('afterupdateitem', (e) => {
    dataChange()
  });

  graph.on('stackchange', (e) => {
    dataChange()
  });

  graph.on('afterrender', (e) => {
    dataChange()
  });

  // some listeners to control the state of nodes to show and hide anchor-point circles
  graph.on('node:mouseenter', e => {
    graph.setItemState(e.item, 'showAnchors', true);
    // setShowState('showAnchors', true, e.item);
  });
  graph.on('node:mouseleave', e => {
    graph.setItemState(e.item, 'showAnchors', false);
    // setShowState('showAnchors', false, e.item);
  })
  graph.on('node:dragenter', e => {
    graph.setItemState(e.item, 'showAnchors', true);
    // setShowState('showAnchors', true, e.item);
  });
  graph.on('node:dragleave', e => {
    graph.setItemState(e.item, 'showAnchors', false);
    // setShowState('showAnchors', false, e.item);
  });
  graph.on('node:dragstart', e => {
    graph.setItemState(e.item, 'showAnchors', true);
    // setShowState('showAnchors', true, e.item);
  });
  graph.on('node:dragout', e => {
    graph.setItemState(e.item, 'showAnchors', false);
    // setShowState('showAnchors', false, e.item);
  });

  graph.on('combo:mouseenter', (evt) => {
    const { item } = evt;
    graph.setItemState(item, 'active', true);
  });
  
  graph.on('combo:mouseleave', (evt) => {
    const { item } = evt;
    graph.setItemState(item, 'active', false);
  });
  // graph.on('combo:click', (evt) => {
  //   const { item } = evt;
  //   graph.setItemState(item, 'selected', true);
  // });
  // graph.on('canvas:click', (evt) => {
  //   graph.getCombos().forEach((combo) => {
  //     graph.clearItemStates(combo);
  //   });
  // });

  graph.on('combo:click', (e) => {
    if (e.target.get('name') === 'combo-marker-shape') {
      // graph.collapseExpandCombo(e.item.getModel().id);
      graph.collapseExpandCombo(e.item);
      if (graph.get('layout')) graph.layout();
      else graph.refreshPositions();
    }
  });


  graph.on('canvas:click', (evt) => {
    graph.getEdges().forEach((edge) => {
      graph.clearItemStates(edge);
    });
    graph.getNodes().forEach((node) => {
      graph.clearItemStates(node);
    });
  });

  graph.on('keyup', (evt) => {
    const srcEleTag = evt.srcElement!.tagName
    console.log(srcEleTag)
    if (evt.key === 'Backspace' && editabel && srcEleTag === 'BODY') {
      console.log(evt)
      graph.getEdges().forEach((edge) => {
        edge.getStates().includes('selected') && graph.removeItem(edge);
      });
      graph.getNodes().forEach((node) => {
        node.getStates().includes('selected') && graph.removeItem(node);
      });
      dataChange()
    }
  })

  if (typeof window !== 'undefined')
    window.onresize = () => {
      if (!graph || graph.get('destroyed')) return;
      if (!container || !container.scrollWidth || !container.scrollHeight) return;
      graph.changeSize(container.scrollWidth, container.scrollHeight);
    };
  return graph
}


