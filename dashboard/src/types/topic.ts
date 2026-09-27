export type ValueType = "number" | "boolean" | "string"

export interface Topic {
  topic: string;
  name: string;
  description: string;
  sourceLabel: string;
  unit: string;
  valueType: ValueType,
}