#include <WiFiNINA.h>
#include <ArduinoHttpClient.h>
#include <ArduinoJson.h>


// WiFi
const char ssid[] = "FamAK";
const char pass[] = "contrasenadewifi";
const char URL[] = "192.168.1.184";
const int port = 8080;

WiFiClient wifi;
HttpClient client = HttpClient(wifi, URL, port);


// IO
int PIN_A = 8;
int PIN_B = 7;
int PIN_C = 6;
int PIN_D = 2;


// State
bool makingRequest = false;


void setup() {

  Serial.begin(9600);

  // WiFi
  while (WiFi.begin(ssid, pass) != WL_CONNECTED) {
    Serial.println("Connecting to WiFi...");
    delay(1000);
  }

  Serial.println("Connected to WiFi");

  pinMode(PIN_A, INPUT_PULLUP);
  pinMode(PIN_B, INPUT_PULLUP);
  pinMode(PIN_C, INPUT_PULLUP);
  pinMode(PIN_D, INPUT_PULLUP);
}

void loop() {
  delay(10);
  int pinAValue = digitalRead(PIN_A);
  int pinBValue = digitalRead(PIN_B);
  int pinCValue = digitalRead(PIN_C);
  int pinDValue = digitalRead(PIN_D);


  StaticJsonDocument<200> doc;
  doc["button"] = -1;
  if (!pinDValue) {
    doc["button"] = 0;
  } else if (!pinAValue) {
    doc["button"] = 1;
  } else if (!pinBValue) {
    doc["button"] = 2;
  } else if (!pinCValue) {
    doc["button"] = 3;
  }

  if (doc["button"] != -1) {
    Serial.print("Making request for: ");
    Serial.println((int)doc["button"]);
    String jsonString;
    serializeJson(doc, jsonString);

    client.put("/update", "application/json", jsonString);
    int statusCode = client.responseStatusCode();
    String response = client.responseBody();

    Serial.print("Status code: ");
    Serial.println(statusCode);
    Serial.print("Response: ");
    Serial.println(response);
  }
}
