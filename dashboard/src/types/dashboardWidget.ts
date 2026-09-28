import type { Topic } from "@/types/topic";

export interface DashboardWidget {
  id: string;
  widget: string;
  topic: Topic;
  name: string;
  description: string;
}

export type NewDashboardWidget = Omit<DashboardWidget, "id">;