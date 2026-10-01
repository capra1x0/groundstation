#include <stdio.h>
#include "esp_adc/adc_oneshot.h"
#include "esp_err.h"
#include "esp_log.h"
#include "freertos/FreeRTOS.h"
#include "freertos/projdefs.h"
#include "freertos/task.h"
#include "hal/adc_types.h"
#include "mqtt.h"
#include "wifi.h"

#define CHANNEL ADC_CHANNEL_0
#define SAMPLES 16
#define ADC_MAX 4095.0f

static const char *TAG = "potentiometer";

static adc_oneshot_unit_handle_t adc_init(void) {
    adc_oneshot_unit_handle_t adc;

    adc_oneshot_unit_init_cfg_t unit_config = {
        .unit_id = ADC_UNIT_1,
    };
    ESP_ERROR_CHECK(adc_oneshot_new_unit(&unit_config, &adc));

    adc_oneshot_chan_cfg_t channel_config = {
        .atten = ADC_ATTEN_DB_12,
        .bitwidth = ADC_BITWIDTH_12,
    };
    ESP_ERROR_CHECK(adc_oneshot_config_channel(adc, CHANNEL, &channel_config));

    return adc;
}

static float read_percent(adc_oneshot_unit_handle_t adc) {
    int sum = 0;

    for(int i = 0; i < SAMPLES; i++) {
        int raw = 0;
        ESP_ERROR_CHECK(adc_oneshot_read(adc, CHANNEL, &raw));
        sum += raw;
    }

    float average = (float)sum / SAMPLES;
    return average / ADC_MAX * 100.0f;
}

void app_main(void) {
    adc_oneshot_unit_handle_t adc = adc_init();

    wifi_connect();
    mqtt_start();

    while(true) {
        float percent = read_percent(adc);
        ESP_LOGI(TAG, "%.1f %%", percent);
        mqtt_publish_number("potentiometer", percent);

        vTaskDelay(pdMS_TO_TICKS(100));
    }
}
