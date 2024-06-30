#include <stdio.h>

void foo() {
  printf("wowee\n");
}

int main(int argc, char *argv[])
{
  int i; 
  void (*funptr[])() = {&foo, &foo}; 
  for (i = 0; i<2 ;i++ ) {
    funptr[i]();
  }
  return 1; 
}
