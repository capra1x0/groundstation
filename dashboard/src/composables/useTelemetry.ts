import { computed, ref, toValue, type MaybeRefOrGetter } from "vue";

import type { Reading } from "@/types/telemetry";

const EVENTS_URL = "http://localhost:1880/events";

const readings = ref<Record<string, Reading>>({});
const connected = ref(false);
let eventSource: EventSource | null = null;

function connect(): void {
  if (eventSource) {
    return;
  }

  eventSource = new EventSource(EVENTS_URL);

  eventSource.onopen = () => {
    connected.value = true;
  };

  eventSource.onerror = () => {
    connected.value = false;
  };

  eventSource.onmessage = (event: MessageEvent<string>) => {
    const reading = JSON.parse(event.data) as Reading;
    readings.value[`${reading.source}/${reading.sensor}`] = reading;
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