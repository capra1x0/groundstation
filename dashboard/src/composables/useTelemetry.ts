import { computed, ref, toValue, type MaybeRefOrGetter } from "vue";

import { API_URL } from "@/api/client";
import type { Reading } from "@/types/reading";

const readings = ref<Record<string, Reading>>({});
const connected = ref(false);
let eventSource: EventSource | null = null;

function connect(): void {
  if (eventSource) {
    return;
  }

  eventSource = new EventSource(`${API_URL}/events`);

  eventSource.onopen = () => {
    connected.value = true;
  };

  eventSource.onerror = () => {
    connected.value = false;
  };

  eventSource.onmessage = (event: MessageEvent<string>) => {
    const reading = JSON.parse(event.data) as Reading;

    readings.value[`telemetry/${reading.source}/${reading.sensor}`] = reading;
  };
}

export function useTelemetry() {
  connect();

  const topics = computed(() => Object.keys(readings.value).sort());

  return { readings, connected, topics };
}

export function useReading(topic: MaybeRefOrGetter<string>) {
  const { readings } = useTelemetry();

  return computed<Reading | undefined>(() => readings.value[toValue(topic)]);
}