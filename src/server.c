#include <fcntl.h>
#include <netinet/in.h>
#include <stdio.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <unistd.h>

char *read_pipe(int fd) {
  static char out[256] = {0};
  int bytes_read = read(fd, out, 255);
  if (bytes_read < 0) {
    perror("read");
    return NULL;
  }
  printf("Bytes read: %d\n", bytes_read);
  return out;
}

char *read_line(int fd) {
  char c = 0;
  int i = 0;
  static char out[256];
  while (read(fd, &c, 1) != 0) {
    if (c == '\n') {
      break;
    }
    out[i] = c;
    i++;
  }

  printf("size of out: %lu\n", sizeof(out));
  return out;
}

int main() {
  int pipefd = open("../pipe", O_RDONLY);
  int socketfd = socket(AF_INET, SOCK_STREAM, 0);
  struct sockaddr address = {AF_INET, htons(9999), 0};

  int true = 1;

  // int opt_err = setsockopt(socketfd, SOL_SOCKET, SO_REUSEADDR, &true,
  // sizeof(int)); if (opt_err == -1) {
  //   perror("opt_err");
  //   return 1;
  // }

  int bind_err = bind(socketfd, &address, sizeof(address));
  if (bind_err == -1) {
    perror("bind error");
    return 1;
  }

  int listen_err = listen(socketfd, 10);
  if (listen_err == -1) {
    perror("listen error");
    return 1;
  }

  // struct pollfd pipefds[1] = {{pipefd, POLLIN, 0}};
  // for(;;){
  //   poll(pipefds, 1, 50000);
  //   if (pipefds[0].revents & POLLIN) {
  //     char *game_state = read_pipe(pipefd);
  //     printf("%s\n", game_state);
  //   }
  // }

  int clientfd = accept(socketfd, 0, 0);
  if (clientfd == -1) {
    perror("accept error");
    return 1;
  }
  printf("connection recieved\n");

  struct pollfd fds[2] = {{pipefd, POLLIN, 0}, {clientfd, POLLIN, 0}};
  for (;;) {
    char buffer[256] = {0};
    poll(fds, 2, 50000);
    if (fds[0].revents & POLLIN) {
      static char out[256] = {0};
      int bytes_read = read(pipefd, out, 255);
      printf("%lu\n", sizeof(out)); 
      long send_err = send(clientfd, out, 255, 0);
      if (send_err == -1) {
        perror("Send Error");
      }

    } else if (fds[1].revents & POLLIN) {
      if (recv(clientfd, buffer, 1, 0) == 0) {
        printf("exiting\n");
        return 0;
      }
      printf("%s\n", buffer);
    }
  }

  close(clientfd);
  close(socketfd);

  return 0;
}
