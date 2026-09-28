import type { Component } from "vue";

import BigNumberWidget from "./BigNumberWidget.vue";

export const widgetComponents: Record<string, Component> = {
  "big-number": BigNumberWidget,
};