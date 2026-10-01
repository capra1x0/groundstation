idf_component_register(
  SRCS "wifi.c" "mqtt.c"
  INCLUDE_DIRS "include"
  PRIV_REQUIRES esp_wifi esp_netif esp_event nvs_flash mqtt
)
