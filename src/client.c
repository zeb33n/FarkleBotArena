#include <netinet/in.h>
#include <stdio.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <unistd.h>

char *read_user() {
  static char out[1] = {0};
  read(0, out, 1);
  return out;
}

char *read_server(int socketfd) {
  static char recv_buf[256] = {0};
  if (recv(socketfd, recv_buf, 255, 0) == 0) {
    // add error handling
  }
  return recv_buf;
}

int main() {
  int socketfd = socket(AF_INET, SOCK_STREAM,
                        0); // what does AF_INET // what does 0 mean here
  struct sockaddr address = {AF_INET, htons(8998), 0};

  int con_err = connect(socketfd, &address, sizeof(address));
  if (con_err == -1) {
    perror("connection error");
  }

  struct pollfd fds[2] = {{0, POLLIN, 0},
                          {socketfd, POLLIN, 0}}; // what does POLLIN mean

  for (;;) {
    poll(fds, 2, 50000);
    if (fds[0].revents & POLLIN) {
      char *msg = read_user();
      printf("%s\n", msg);
      send(socketfd, msg, 1, 0);
      if (*msg == 'q') {
        return 0;
      }

    } else if (fds[1].revents & POLLIN) {
      char *game_state = read_server(socketfd);
      printf("%s\n", game_state);
    }
  }

  return 0;
}

// char *read_pipe() {
//   int fd = open("../pipe", O_RDONLY);
//   static char out[256] = {0};
//   int _ = read(fd, out, 255);
//   return out;
// }
