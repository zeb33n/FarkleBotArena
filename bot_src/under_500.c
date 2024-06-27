#include <cjson/cJSON.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
  char *json = argv[1];
  cJSON *root = cJSON_Parse(json); 
  int score = cJSON_GetObjectItem(root, "round_score")->valueint; 
  if (score < 500) {
    printf("1\n");
  } else {
    printf("0\n");
  }
  return 0;
}
