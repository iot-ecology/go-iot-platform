import type { App } from "vue";
import { Checkbox, Column, Grid, Table, Tooltip, VXETable } from "vxe-table";
import "xe-utils";
import "vxe-table/lib/style.css";

const components = [Checkbox, Column, Grid, Table, Tooltip];

// 单元格默认渲染器
VXETable.renderer.add("cellRender", {
  renderDefault(_renderOpts, params) {
    const { row, column } = params;
    const value = row[column.field];
    if (value === null || value === "" || (Array.isArray(value) && value.length === 0)) {
      return "-";
    } else if (Array.isArray(value)) {
      return value.join(",");
    } else {
      return value;
    }
  },
});

export default {
  install(app: App) {
    for (const comp of components) {
      app.use(comp);
    }
  },
};
