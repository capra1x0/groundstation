import type { Reading } from "@/types/reading";
import { get } from "./client";

export function fetchHistory(topic: string, from?: number, to?: number): Promise<Reading[]> {
  const params = new URLSearchParams({ topic });

  if(from != undefined) {
    params.set("from", from.toString())
  }
  if(to != undefined) {
    params.set("to", to.toString())
  }

  return get<Reading[]>(`/history?${params}`);
}