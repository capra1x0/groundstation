import { get } from "./client";

import type { Topic } from "@/types/topic";

export function getTopics(): Promise<Topic[]> {
  return get<Topic[]>("/topics");
}