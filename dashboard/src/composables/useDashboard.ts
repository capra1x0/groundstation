import { ref } from "vue";
import { v4 as uuidv4 } from "uuid";

import type { DashboardWidget, NewDashboardWidget } from "@/types/dashboardWidget";

const widgets = ref<DashboardWidget[]>([]);

function addWidget(widget: NewDashboardWidget): void {
  widgets.value.push({ id: uuidv4(), settings: {}, ...widget });
}

function removeWidget(id: string): void {
  widgets.value = widgets.value.filter((widget) => widget.id !== id);
}

function updateSettings(id: string, settings: Record<string, number>): void {
  const widget = widgets.value.find((widget) => widget.id === id);
  if (widget) {
    widget.settings = { ...widget.settings, ...settings };
  }
}

export function useDashboard() {
  return { widgets, addWidget, removeWidget, updateSettings };
}