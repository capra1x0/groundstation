#include "esp_err.h"
#include "freertos/FreeRTOS.h"
#include "freertos/projdefs.h"
#include "freertos/task.h"
#include "mqtt.h"
#include "nvs.h"
#include "nvs_flash.h"
#include "wifi.h"
#include "dht.h"
#include "esp_log.h"

#define DHT_PIN GPIO_NUM_4

static const char *TAG = "climate";

void app_main(void) {
    esp_err_t err = nvs_flash_init();
    if(err == ESP_ERR_NVS_NO_FREE_PAGES || err == ESP_ERR_NVS_NEW_VERSION_FOUND) {
        ESP_ERROR_CHECK(nvs_flash_erase());
        err = nvs_flash_init();
    }
    ESP_ERROR_CHECK(err);

    wifi_connect();
    mqtt_start();

    while(true) {
        float humidity;
        float temp;

        esp_err_t result = dht_read_float_data(DHT_TYPE_DHT11, DHT_PIN, &humidity, &temp);
        if(result == ESP_OK) {
            ESP_LOGI(TAG, "%.1f C, %.1f %%", temp, humidity);
            mqtt_publish_number("temp", temp);
            mqtt_publish_number("humidity", humidity);
        } else {
            ESP_LOGW(TAG, "reading failed: %s", esp_err_to_name(result));
        }

        vTaskDelay(pdMS_TO_TICKS(10000));
    }
}
