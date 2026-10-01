#include "mqtt.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include "esp_event_base.h"
#include "esp_log.h"
#include "mqtt_client.h"
#include "sdkconfig.h"

static const char *TAG = "mqtt";
static esp_mqtt_client_handle_t client;
static volatile bool connected = false;

static void on_mqtt_event(void *arg, esp_event_base_t base, int32_t id, void *data) {
    if(id == MQTT_EVENT_CONNECTED) {
        connected = true;
        ESP_LOGI(TAG, "connected to broker");
    } else if(id == MQTT_EVENT_DISCONNECTED) {
        connected = false;
        ESP_LOGW(TAG, "connection to broker lost"); 
    }
}

void mqtt_start(void) {
    esp_mqtt_client_config_t config = {
        .broker.address.uri = CONFIG_MQTT_BROKER_URI,
    };

    client = esp_mqtt_client_init(&config);
    esp_mqtt_client_register_event(client, ESP_EVENT_ANY_ID, on_mqtt_event, NULL);
    esp_mqtt_client_start(client);
}

void mqtt_publish_number(const char *sensor, float value) {
    if(!connected) {
        return;
    }

    char topic[64];
    char payload[32];

    snprintf(topic, sizeof(topic), "telemetry/esp/%s", sensor);
    snprintf(payload, sizeof(payload), "{\"value\":%.1f}", value);

    esp_mqtt_client_publish(client, topic, payload, 0, 0, 0);
}
