import type { Component } from "vue";

import BigNumberWidget from "./BigNumberWidget.vue";
import GaugeWidget from "./GaugeWidget.vue";
import LineChartWidget from "./LineChartWidget.vue";

export const widgetComponents: Record<string, Component> = {
  "big-number": BigNumberWidget,
  "gauge": GaugeWidget,
  "line-chart": LineChartWidget,
};