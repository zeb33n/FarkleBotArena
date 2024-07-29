#include <fcntl.h>
#include <stdio.h>
#include <sys/poll.h>
#include <sys/types.h>
#include <unistd.h>

char read_user() {
  char out = '\0';
  read(0, &out, 1);
  return out;
}

int main() {
  int p = 0;
  for (;;) {
    char in = read_user();
    if (in == '1') {
      p = fork();
      in = 0;
    }
    if (p != 0) {
      break;
    }
  }
  pid_t id = getpid();
  printf("p: %d, id: %d\n", p, id);
  for (;;) {
    printf("");
  }
  return 0;
}
