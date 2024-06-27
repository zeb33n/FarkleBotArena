#include <stdio.h>
#include <stdlib.h>
#include <time.h> 

int main() {
  srand(time(NULL)); 
  int r = rand() % 2;
  printf("%i\n", r);
  return 0; 
}
