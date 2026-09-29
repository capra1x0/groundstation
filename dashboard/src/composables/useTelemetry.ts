import { computed, ref, toValue, type MaybeRefOrGetter } from "vue";

import { API_URL } from "@/api/client";
import type { Reading } from "@/types/reading";

const HISTORY_DURATION_MS = 5 * 60 * 1000;

const readings = ref<Record<string, Reading>>({});
const connected = ref(false);
let eventSource: EventSource | null = null;

const history = new Map<string, Reading[]>();

function storeReading(reading: Reading): void {
  readings.value[reading.topic] = reading;

  let list = history.get(reading.topic);
  if(!list) {
    list = [];
    history.set(reading.topic, list);
  }
  list.push(reading);

  const cutoff = Date.now() - HISTORY_DURATION_MS;
  const firstRecent = list.findIndex((entry) => entry.ts >= cutoff);
  const tooOld = firstRecent === -1 ? list.length : firstRecent;
 
  if (tooOld > 0) {
    list.splice(0, tooOld);
  }
}

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
    storeReading(reading)
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

export function getHistory(topic: string): Reading[] {
  connect();

  return history.get(topic) ?? [];
}