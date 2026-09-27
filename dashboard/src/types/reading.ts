export interface Reading {
  topic: string;
  ts: number;
  value: number | boolean | string;
  unit: string;
}