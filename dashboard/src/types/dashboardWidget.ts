import type { Topic } from "@/types/topic";

export interface DashboardWidget {
  id: string;
  widget: string;
  topic: Topic;
  name: string;
  description: string;
  settings: Record<string, number>;
}

export type NewDashboardWidget = Omit<DashboardWidget, "id" | "settings">;