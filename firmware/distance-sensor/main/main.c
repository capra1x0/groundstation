#include "esp_err.h"
#include "freertos/projdefs.h"
#include "stdbool.h"
#include "esp_log.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "mqtt.h"
#include "ultrasonic.h"
#include "wifi.h"
#include <stdint.h>
#include <stdlib.h>

#define TRIGGER_PIN GPIO_NUM_17
#define ECHO_PIN GPIO_NUM_16
#define MAX_DISTANCE_CM 400
#define SAMPLES 7

static const char *TAG = "distance";

int comp(const void *a, const void *b) {
    int x = *(const int *)a;
    int y = *(const int *)b;

    if (x < y)
        return -1;
    if (x > y)
        return 1;
    return 0;
}

static esp_err_t measure_medain_cm(ultrasonic_sensor_t *sensor, uint32_t *result) {
    uint32_t readings[SAMPLES];
    int count = 0;

    for(int i = 0; i < SAMPLES; i++) {
        uint32_t distance;
        if(ultrasonic_measure_cm(sensor, MAX_DISTANCE_CM, &distance) == ESP_OK) {
            readings[count++] = distance;
        }
        vTaskDelay(pdMS_TO_TICKS(60));
    }

    if(count == 0) {
        return ESP_FAIL;
    }

    qsort(readings, count, sizeof(uint32_t), comp);
    *result = readings[count / 2];
    return ESP_OK;
}


void app_main(void) {
    ultrasonic_sensor_t sensor = {
        .trigger_pin = TRIGGER_PIN,
        .echo_pin = ECHO_PIN,
    };
    ESP_ERROR_CHECK(ultrasonic_init(&sensor));

    wifi_connect();
    mqtt_start();

    while(true) {
        uint32_t distance;

        esp_err_t result = measure_medain_cm(&sensor, &distance);
        if(result == ESP_OK) {
            ESP_LOGI(TAG, "%" PRIu32 " cm", distance);
            mqtt_publish_number("distance", (float)distance);
        } else {
            ESP_LOGW(TAG, "failed to get distance: 0x%x", result);
        }

        vTaskDelay(pdMS_TO_TICKS(500));
    }
}
