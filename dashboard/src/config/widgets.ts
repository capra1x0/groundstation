import type { ValueType } from "@/types/topic";

export interface WidgetOption {
  id: string;
  name: string;
  description: string;
  preview: string;
  valueTypes: ValueType[];
}

export const widgetOptions: WidgetOption[] = [
  {
    id: "big-number",
    name: "Big Number",
    description: "Latest Value as a large number",
    preview: "/previews/big-number.svg",
    valueTypes: ["number"],
  },
  {
    id: "gauge",
    name: "Gauge",
    description: "Current value within a range",
    preview: "/previews/gauge.svg",
    valueTypes: ["number"],
  },
  {
    id: "line-chart",
    name: "Line chart",
    description: "Value over time",
    preview: "/previews/line-chart.svg",
    valueTypes: ["number"],
  },
]