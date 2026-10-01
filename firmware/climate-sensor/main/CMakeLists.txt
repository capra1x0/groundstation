idf_component_register(
  SRCS "main.c" "wifi.c" "mqtt.c"
  INCLUDE_DIRS "."
  REQUIRES esp_driver_gpio esp_timer esp_wifi esp_netif esp_event nvs_flash mqtt
)
